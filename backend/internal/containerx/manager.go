// Package containerx defines the interface for Docker container lifecycle management.
// It is used by the agentsession domain to manage agent workspace containers.
package containerx

import "context"

// ContainerConfig holds the configuration for creating a container.
type ContainerConfig struct {
	Image       string            // Docker image reference
	Name        string            // Container name
	Env         map[string]string // Environment variables
	Ports       []PortBinding     // Port bindings
	VolumeName  string            // Named volume to mount
	NetworkName string            // Network to attach to
	Labels      map[string]string // Labels
}

// PortBinding maps a container port to a host port.
type PortBinding struct {
	ContainerPort int
	HostPort      int    // 0 = dynamic allocation
	Protocol      string // "tcp" or "udp"
}

// ContainerInfo holds runtime information about a running container.
type ContainerInfo struct {
	ID          string
	Name        string
	Status      string
	Ports       []PortBinding
	NetworkID   string
	VolumeName  string
}

// Manager defines the interface for Docker container lifecycle management.
type Manager interface {
	// CreateAndStart creates and starts a container from the given config.
	// Returns the container ID and port bindings assigned.
	CreateAndStart(ctx context.Context, cfg ContainerConfig) (ContainerInfo, error)

	// Stop stops a running container gracefully.
	Stop(ctx context.Context, containerID string) error

	// Remove removes a container (must be stopped first).
	Remove(ctx context.Context, containerID string) error

	// Inspect returns current info for a container.
	Inspect(ctx context.Context, containerID string) (ContainerInfo, error)

	// CreateVolume creates a named Docker volume.
	CreateVolume(ctx context.Context, name string) error

	// RemoveVolume removes a named Docker volume.
	RemoveVolume(ctx context.Context, name string) error

	// CreateNetwork creates a Docker network.
	CreateNetwork(ctx context.Context, name string) (string, error)

	// RemoveNetwork removes a Docker network.
	RemoveNetwork(ctx context.Context, networkID string) error
}
