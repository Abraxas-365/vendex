// Package containerxdocker implements containerx.Manager using the Docker Engine SDK.
package containerxdocker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Abraxas-365/vendex/internal/containerx"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Compile-time check.
var _ containerx.Manager = (*Client)(nil)

// Client wraps the Docker Engine SDK client.
type Client struct {
	docker *client.Client
}

// New creates a new Docker containerx manager.
// It connects to the Docker daemon using the default environment settings
// (DOCKER_HOST, DOCKER_TLS_VERIFY, etc.).
func New() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: create docker client: %w", err)
	}
	return &Client{docker: cli}, nil
}

// NewFromClient wraps an existing Docker client.
func NewFromClient(cli *client.Client) *Client {
	return &Client{docker: cli}
}

func (c *Client) Create(ctx context.Context, spec containerx.Spec) (containerx.ID, error) {
	// Build environment slice.
	env := make([]string, 0, len(spec.Env))
	for k, v := range spec.Env {
		env = append(env, k+"="+v)
	}

	// Build port bindings and exposed ports.
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}
	for _, p := range spec.Ports {
		proto := p.Protocol
		if proto == "" {
			proto = "tcp"
		}
		containerPort := nat.Port(fmt.Sprintf("%d/%s", p.ContainerPort, proto))
		exposedPorts[containerPort] = struct{}{}
		portBindings[containerPort] = []nat.PortBinding{
			{HostPort: strconv.Itoa(p.HostPort)},
		}
	}

	// Build mounts.
	binds := make([]string, 0, len(spec.Mounts)+len(spec.Volumes))
	for _, m := range spec.Mounts {
		bind := m.HostPath + ":" + m.ContainerPath
		if m.ReadOnly {
			bind += ":ro"
		}
		binds = append(binds, bind)
	}
	for _, v := range spec.Volumes {
		bind := v.Name + ":" + v.ContainerPath
		if v.ReadOnly {
			bind += ":ro"
		}
		binds = append(binds, bind)
	}

	// Build labels.
	labels := spec.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	// Restart policy.
	restartPolicy := containertypes.RestartPolicy{Name: containertypes.RestartPolicyDisabled}
	if spec.RestartPolicy != "" {
		parts := strings.SplitN(spec.RestartPolicy, ":", 2)
		switch parts[0] {
		case "always":
			restartPolicy = containertypes.RestartPolicy{Name: containertypes.RestartPolicyAlways}
		case "unless-stopped":
			restartPolicy = containertypes.RestartPolicy{Name: containertypes.RestartPolicyUnlessStopped}
		case "on-failure":
			restartPolicy = containertypes.RestartPolicy{Name: containertypes.RestartPolicyOnFailure}
			if len(parts) == 2 {
				if n, err := strconv.Atoi(parts[1]); err == nil {
					restartPolicy.MaximumRetryCount = n
				}
			}
		}
	}

	config := &containertypes.Config{
		Image:        spec.Image,
		Env:          env,
		ExposedPorts: exposedPorts,
		Labels:       labels,
		WorkingDir:   spec.WorkDir,
	}
	if len(spec.Cmd) > 0 {
		config.Cmd = spec.Cmd
	}

	hostConfig := &containertypes.HostConfig{
		Binds:         binds,
		PortBindings:  portBindings,
		RestartPolicy: restartPolicy,
		Resources: containertypes.Resources{
			CPUShares: spec.Resources.CPUShares,
			Memory:    spec.Resources.MemoryBytes,
		},
	}

	networkConfig := &network.NetworkingConfig{}
	if spec.Network != "" {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			spec.Network: {},
		}
	}

	resp, err := c.docker.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, spec.Name)
	if err != nil {
		return "", fmt.Errorf("containerxdocker: create: %w", err)
	}

	return containerx.ID(resp.ID), nil
}

func (c *Client) Start(ctx context.Context, id containerx.ID) error {
	if err := c.docker.ContainerStart(ctx, string(id), containertypes.StartOptions{}); err != nil {
		return fmt.Errorf("containerxdocker: start %s: %w", id, err)
	}
	return nil
}

func (c *Client) Stop(ctx context.Context, id containerx.ID, timeout time.Duration) error {
	timeoutSec := int(timeout.Seconds())
	opts := containertypes.StopOptions{Timeout: &timeoutSec}
	if err := c.docker.ContainerStop(ctx, string(id), opts); err != nil {
		return fmt.Errorf("containerxdocker: stop %s: %w", id, err)
	}
	return nil
}

func (c *Client) Remove(ctx context.Context, id containerx.ID) error {
	opts := containertypes.RemoveOptions{Force: true, RemoveVolumes: false}
	if err := c.docker.ContainerRemove(ctx, string(id), opts); err != nil {
		return fmt.Errorf("containerxdocker: remove %s: %w", id, err)
	}
	return nil
}

func (c *Client) Status(ctx context.Context, id containerx.ID) (containerx.Status, error) {
	info, err := c.docker.ContainerInspect(ctx, string(id))
	if err != nil {
		return containerx.Status{}, fmt.Errorf("containerxdocker: inspect %s: %w", id, err)
	}

	state := mapDockerState(info.State)
	ports := extractPorts(info.NetworkSettings)
	errMsg := ""
	if info.State != nil && info.State.Error != "" {
		errMsg = info.State.Error
	}

	return containerx.Status{
		ID:    id,
		State: state,
		Ports: ports,
		Error: errMsg,
	}, nil
}

func (c *Client) Logs(ctx context.Context, id containerx.ID) (io.ReadCloser, error) {
	opts := containertypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}
	reader, err := c.docker.ContainerLogs(ctx, string(id), opts)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: logs %s: %w", id, err)
	}
	return reader, nil
}

func (c *Client) Exec(ctx context.Context, id containerx.ID, cmd []string) ([]byte, error) {
	execConfig := containertypes.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := c.docker.ContainerExecCreate(ctx, string(id), execConfig)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: exec create %s: %w", id, err)
	}

	attachResp, err := c.docker.ContainerExecAttach(ctx, execResp.ID, containertypes.ExecStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: exec attach %s: %w", id, err)
	}
	defer attachResp.Close()

	output, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: exec read %s: %w", id, err)
	}

	return output, nil
}

func (c *Client) ExecStream(ctx context.Context, id containerx.ID, cmd []string) (io.ReadCloser, error) {
	execConfig := containertypes.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := c.docker.ContainerExecCreate(ctx, string(id), execConfig)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: exec create %s: %w", id, err)
	}

	attachResp, err := c.docker.ContainerExecAttach(ctx, execResp.ID, containertypes.ExecStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: exec attach %s: %w", id, err)
	}

	return io.NopCloser(attachResp.Reader), nil
}

func (c *Client) CreateNetwork(ctx context.Context, name string) (string, error) {
	resp, err := c.docker.NetworkCreate(ctx, name, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return "", fmt.Errorf("containerxdocker: create network %q: %w", name, err)
	}
	return resp.ID, nil
}

func (c *Client) RemoveNetwork(ctx context.Context, id string) error {
	if err := c.docker.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("containerxdocker: remove network %s: %w", id, err)
	}
	return nil
}

func (c *Client) CreateVolume(ctx context.Context, name string) error {
	_, err := c.docker.VolumeCreate(ctx, volume.CreateOptions{Name: name})
	if err != nil {
		return fmt.Errorf("containerxdocker: create volume %q: %w", name, err)
	}
	return nil
}

func (c *Client) RemoveVolume(ctx context.Context, name string) error {
	if err := c.docker.VolumeRemove(ctx, name, true); err != nil {
		return fmt.Errorf("containerxdocker: remove volume %q: %w", name, err)
	}
	return nil
}

func (c *Client) CopyToContainer(ctx context.Context, id containerx.ID, dstPath string, content io.Reader) error {
	opts := containertypes.CopyToContainerOptions{AllowOverwriteDirWithFile: true}
	if err := c.docker.CopyToContainer(ctx, string(id), dstPath, content, opts); err != nil {
		return fmt.Errorf("containerxdocker: copy to %s:%s: %w", id, dstPath, err)
	}
	return nil
}

func (c *Client) CopyFromContainer(ctx context.Context, id containerx.ID, srcPath string) ([]byte, error) {
	reader, _, err := c.docker.CopyFromContainer(ctx, string(id), srcPath)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: copy from %s:%s: %w", id, srcPath, err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: read from %s:%s: %w", id, srcPath, err)
	}
	return data, nil
}

func (c *Client) Stats(ctx context.Context, id containerx.ID) (*containerx.Stats, error) {
	resp, err := c.docker.ContainerStatsOneShot(ctx, string(id))
	if err != nil {
		return nil, fmt.Errorf("containerxdocker: stats %s: %w", id, err)
	}
	defer resp.Body.Close()

	var raw struct {
		CPUStats struct {
			CPUUsage struct {
				TotalUsage uint64 `json:"total_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage uint64 `json:"system_cpu_usage"`
		} `json:"cpu_stats"`
		PreCPUStats struct {
			CPUUsage struct {
				TotalUsage uint64 `json:"total_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage uint64 `json:"system_cpu_usage"`
		} `json:"precpu_stats"`
		MemoryStats struct {
			Usage uint64 `json:"usage"`
			Limit uint64 `json:"limit"`
		} `json:"memory_stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("containerxdocker: decode stats %s: %w", id, err)
	}

	cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage - raw.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPUStats.SystemCPUUsage - raw.PreCPUStats.SystemCPUUsage)
	cpuPercent := 0.0
	if sysDelta > 0 {
		cpuPercent = (cpuDelta / sysDelta) * 100.0
	}

	return &containerx.Stats{
		CPUPercent:  cpuPercent,
		MemoryUsage: raw.MemoryStats.Usage,
		MemoryLimit: raw.MemoryStats.Limit,
		Timestamp:   time.Now().UTC(),
	}, nil
}

func (c *Client) PullImage(ctx context.Context, img string) error {
	// Check if image exists locally first.
	_, _, err := c.docker.ImageInspectWithRaw(ctx, img)
	if err == nil {
		return nil // already present
	}

	reader, err := c.docker.ImagePull(ctx, img, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("containerxdocker: pull %q: %w", img, err)
	}
	defer reader.Close()

	// Drain the pull output to complete the operation.
	_, _ = io.Copy(io.Discard, reader)
	return nil
}

// --- Helpers ---

func mapDockerState(state *types.ContainerState) containerx.State {
	if state == nil {
		return containerx.StateError
	}
	switch {
	case state.Running:
		return containerx.StateRunning
	case state.Restarting:
		return containerx.StateCreating
	case state.Dead || state.OOMKilled:
		return containerx.StateError
	default:
		return containerx.StateStopped
	}
}

func extractPorts(ns *types.NetworkSettings) []containerx.Port {
	if ns == nil {
		return nil
	}
	var ports []containerx.Port
	for containerPort, bindings := range ns.Ports {
		for _, b := range bindings {
			hp, _ := strconv.Atoi(b.HostPort)
			ports = append(ports, containerx.Port{
				HostPort:      hp,
				ContainerPort: containerPort.Int(),
				Protocol:      containerPort.Proto(),
			})
		}
	}
	return ports
}

