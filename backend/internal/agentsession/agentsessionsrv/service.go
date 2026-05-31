// Package agentsessionsrv implements business logic for agent session lifecycle management.
package agentsessionsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/agentsession"
	"github.com/Abraxas-365/vendex/internal/containerx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace/marketplacesrv"
)

// Service manages agent session lifecycle.
type Service struct {
	sessionRepo agentsession.SessionRepository
	chatRepo    agentsession.ChatRepository
	containers  containerx.Manager
	presetSvc   *marketplacesrv.PresetService
}

// NewService creates a new agentsession Service.
func NewService(
	sessionRepo agentsession.SessionRepository,
	chatRepo agentsession.ChatRepository,
	containers containerx.Manager,
	presetSvc *marketplacesrv.PresetService,
) *Service {
	return &Service{
		sessionRepo: sessionRepo,
		chatRepo:    chatRepo,
		containers:  containers,
		presetSvc:   presetSvc,
	}
}

// CreateSession creates a new agent session from a preset.
// It pulls the Docker image, creates a volume/network, starts the container,
// and persists the session record.
func (s *Service) CreateSession(ctx context.Context, tenantID kernel.TenantID, req agentsession.CreateSessionRequest) (agentsession.Session, error) {
	// Verify the preset exists
	preset, err := s.presetSvc.Get(ctx, req.PresetID)
	if err != nil {
		return agentsession.Session{}, err
	}

	now := time.Now()
	sessionID := kernel.AgentSessionID(uuid.New().String())

	// Create initial session record in "creating" state
	sess := agentsession.Session{
		ID:        sessionID,
		TenantID:  tenantID,
		PresetID:  req.PresetID,
		Status:    agentsession.SessionStatusCreating,
		Metadata:  req.Config,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if len(sess.Metadata) == 0 {
		sess.Metadata = json.RawMessage("{}")
	}

	sess, err = s.sessionRepo.Create(ctx, sess)
	if err != nil {
		return agentsession.Session{}, err
	}

	// Provision Docker resources
	volumeName := fmt.Sprintf("agentsession-%s", string(sessionID))
	networkName := fmt.Sprintf("agentsession-net-%s", string(sessionID))

	if err := s.containers.CreateVolume(ctx, volumeName); err != nil {
		_ = s.markFailed(ctx, sess)
		return agentsession.Session{}, err
	}

	networkID, err := s.containers.CreateNetwork(ctx, networkName)
	if err != nil {
		_ = s.containers.RemoveVolume(ctx, volumeName)
		_ = s.markFailed(ctx, sess)
		return agentsession.Session{}, err
	}

	spec := containerx.Spec{
		Image: preset.Image,
		Name:  fmt.Sprintf("agentsession-%s", string(sessionID)),
		Env: map[string]string{
			"SESSION_ID": string(sessionID),
			"TENANT_ID":  string(tenantID),
			"PRESET_ID":  string(req.PresetID),
		},
		Ports: []containerx.Port{
			{ContainerPort: preset.FrontendPort, HostPort: 0, Protocol: "tcp"},
		},
		Volumes: []containerx.Volume{
			{Name: volumeName, ContainerPath: "/workspace"},
		},
		Network: networkID,
		Labels: map[string]string{
			"managed-by": "hada-agentsession",
			"session-id": string(sessionID),
			"tenant-id":  string(tenantID),
		},
	}

	containerID, err := s.containers.Create(ctx, spec)
	if err != nil {
		_ = s.containers.RemoveNetwork(ctx, networkID)
		_ = s.containers.RemoveVolume(ctx, volumeName)
		_ = s.markFailed(ctx, sess)
		return agentsession.Session{}, err
	}

	if err := s.containers.Start(ctx, containerID); err != nil {
		_ = s.containers.Remove(ctx, containerID)
		_ = s.containers.RemoveNetwork(ctx, networkID)
		_ = s.containers.RemoveVolume(ctx, volumeName)
		_ = s.markFailed(ctx, sess)
		return agentsession.Session{}, err
	}

	// Get runtime status to discover assigned host ports
	status, err := s.containers.Status(ctx, containerID)
	if err != nil {
		_ = s.containers.Stop(ctx, containerID, 10*time.Second)
		_ = s.containers.Remove(ctx, containerID)
		_ = s.containers.RemoveNetwork(ctx, networkID)
		_ = s.containers.RemoveVolume(ctx, volumeName)
		_ = s.markFailed(ctx, sess)
		return agentsession.Session{}, err
	}

	// Determine the frontend URL from assigned host port
	frontendURL := ""
	for _, p := range status.Ports {
		if p.ContainerPort == preset.FrontendPort && p.HostPort > 0 {
			frontendURL = fmt.Sprintf("http://localhost:%d", p.HostPort)
			break
		}
	}

	// Update session to running
	sess.ContainerID = string(containerID)
	sess.Status = agentsession.SessionStatusRunning
	sess.FrontendURL = frontendURL
	sess.UpdatedAt = time.Now()

	return s.sessionRepo.Update(ctx, sess)
}

// GetSession retrieves a session by ID, scoped to a tenant.
func (s *Service) GetSession(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentSessionID) (agentsession.Session, error) {
	return s.sessionRepo.GetByID(ctx, tenantID, id)
}

// ListSessions returns paginated sessions for a tenant.
func (s *Service) ListSessions(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agentsession.Session], error) {
	return s.sessionRepo.ListByTenant(ctx, tenantID, p)
}

// StopSession stops a running container and marks the session as stopped.
func (s *Service) StopSession(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentSessionID) (agentsession.Session, error) {
	sess, err := s.sessionRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return agentsession.Session{}, err
	}
	if sess.Status != agentsession.SessionStatusRunning {
		return agentsession.Session{}, agentsession.ErrSessionNotRunning
	}

	if sess.ContainerID != "" {
		if err := s.containers.Stop(ctx, containerx.ID(sess.ContainerID), 30*time.Second); err != nil {
			return agentsession.Session{}, err
		}
	}

	now := time.Now()
	sess.Status = agentsession.SessionStatusStopped
	sess.StoppedAt = &now
	sess.UpdatedAt = now

	return s.sessionRepo.Update(ctx, sess)
}

// SendMessage appends a message to the session's chat history.
func (s *Service) SendMessage(ctx context.Context, tenantID kernel.TenantID, sessionID kernel.AgentSessionID, role, content, toolName string) (agentsession.ChatMessage, error) {
	// Verify the session exists and belongs to the tenant
	if _, err := s.sessionRepo.GetByID(ctx, tenantID, sessionID); err != nil {
		return agentsession.ChatMessage{}, err
	}

	msg := agentsession.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		ToolName:  toolName,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.SaveMessage(ctx, msg); err != nil {
		return agentsession.ChatMessage{}, err
	}
	return msg, nil
}

// GetHistory returns paginated chat history for a session.
func (s *Service) GetHistory(ctx context.Context, tenantID kernel.TenantID, sessionID kernel.AgentSessionID, p kernel.PaginationOptions) (kernel.Paginated[agentsession.ChatMessage], error) {
	// Verify the session exists and belongs to the tenant
	if _, err := s.sessionRepo.GetByID(ctx, tenantID, sessionID); err != nil {
		return kernel.Paginated[agentsession.ChatMessage]{}, err
	}
	return s.chatRepo.ListMessages(ctx, sessionID, p)
}

// markFailed updates a session's status to failed.
func (s *Service) markFailed(ctx context.Context, sess agentsession.Session) error {
	sess.Status = agentsession.SessionStatusFailed
	sess.UpdatedAt = time.Now()
	_, err := s.sessionRepo.Update(ctx, sess)
	return err
}
