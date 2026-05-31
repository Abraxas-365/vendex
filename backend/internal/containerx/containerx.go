// Package containerx defines the container lifecycle abstraction for agent workspaces.
// It provides a Manager interface for creating, starting, stopping, and interacting
// with Docker containers that serve as isolated workspaces for agent presets.
package containerx

import (
	"context"
	"io"
	"time"
)

// ID is a unique container identifier (Docker container ID).
type ID string

// State represents the lifecycle state of a container.
type State string

const (
	StateCreating State = "creating"
	StateRunning  State = "running"
	StateStopped  State = "stopped"
	StateError    State = "error"
)

// Port maps a host port to a container port.
type Port struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"` // "tcp" or "udp"
}

// Mount describes a host-to-container filesystem binding.
type Mount struct {
	HostPath      string
	ContainerPath string
	ReadOnly      bool
}

// Volume describes a named Docker volume to mount into a container.
type Volume struct {
	Name          string
	ContainerPath string
	ReadOnly      bool
}

// ResourceLimits constrains CPU and memory usage for a container.
type ResourceLimits struct {
	CPUShares   int64 // relative weight (default 1024)
	MemoryBytes int64 // hard memory limit in bytes
}

// Spec defines the configuration used to create a container.
type Spec struct {
	Name    string // optional container name
	Image   string
	Cmd     []string // optional command override
	Env     map[string]string
	Mounts  []Mount
	Volumes []Volume
	Ports   []Port
	Labels  map[string]string

	Resources     ResourceLimits
	Network       string // Docker network name to connect to
	WorkDir       string // working directory inside the container
	RestartPolicy string // "no", "on-failure:3", "unless-stopped", "always"
}

// Status holds the observed runtime state of a container.
type Status struct {
	ID    ID
	State State
	Ports []Port // actual port mappings (host ports resolved)
	Error string
}

// Stats holds resource usage metrics for a running container.
type Stats struct {
	CPUPercent  float64   `json:"cpu_percent"`
	MemoryUsage uint64   `json:"memory_usage"`
	MemoryLimit uint64   `json:"memory_limit"`
	Timestamp   time.Time `json:"timestamp"`
}

// Manager is the core abstraction for managing container lifecycles.
// Implementations include Docker (production) and mock (testing).
type Manager interface {
	// Create provisions a new container from the given spec and returns its ID.
	Create(ctx context.Context, spec Spec) (ID, error)

	// Start runs a previously created container.
	Start(ctx context.Context, id ID) error

	// Stop signals a running container to stop, waiting up to timeout.
	Stop(ctx context.Context, id ID, timeout time.Duration) error

	// Remove deletes a stopped container and its resources.
	Remove(ctx context.Context, id ID) error

	// Status returns the current runtime status of a container.
	Status(ctx context.Context, id ID) (Status, error)

	// Logs returns a stream of the container's combined stdout/stderr output.
	Logs(ctx context.Context, id ID) (io.ReadCloser, error)

	// Exec runs a command inside a running container and returns combined output.
	Exec(ctx context.Context, id ID, cmd []string) ([]byte, error)

	// ExecStream runs a command inside a running container and streams output.
	ExecStream(ctx context.Context, id ID, cmd []string) (io.ReadCloser, error)

	// CreateNetwork provisions a new Docker bridge network and returns its ID.
	CreateNetwork(ctx context.Context, name string) (string, error)

	// RemoveNetwork removes a Docker network by ID.
	RemoveNetwork(ctx context.Context, id string) error

	// CreateVolume provisions a new named Docker volume.
	CreateVolume(ctx context.Context, name string) error

	// RemoveVolume removes a named Docker volume.
	RemoveVolume(ctx context.Context, name string) error

	// CopyToContainer copies data into the container at dstPath.
	// Content should be a tar archive.
	CopyToContainer(ctx context.Context, id ID, dstPath string, content io.Reader) error

	// CopyFromContainer reads a single file at srcPath from the container.
	CopyFromContainer(ctx context.Context, id ID, srcPath string) ([]byte, error)

	// Stats returns a one-shot resource usage snapshot for a running container.
	Stats(ctx context.Context, id ID) (*Stats, error)

	// PullImage pulls a container image if not present locally.
	PullImage(ctx context.Context, image string) error
}
