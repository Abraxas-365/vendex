package agent

import (
	"github.com/Abraxas-365/hada-commerce/internal/agent/workspace"
	"github.com/Abraxas-365/hada-commerce/internal/containerx"
)

// WorkspaceTools returns all workspace tools as a slice of agent.Tool.
// These can be passed directly to AdaptTools() and appended to domain tools.
// Returns nil if mgr or accessor is nil.
func WorkspaceTools(mgr containerx.Manager, accessor workspace.ContainerAccessor) []Tool {
	if mgr == nil || accessor == nil {
		return nil
	}
	return []Tool{
		workspace.NewWriteFileTool(mgr, accessor),
		workspace.NewReadFileTool(mgr, accessor),
		workspace.NewListFilesTool(mgr, accessor),
		workspace.NewPreviewURLTool(accessor),
		workspace.NewDeleteFileTool(mgr, accessor),
	}
}
