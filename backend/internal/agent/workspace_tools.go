package agent

import (
	"github.com/Abraxas-365/vendex/internal/agent/workspace"
	"github.com/Abraxas-365/vendex/internal/containerx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront/storefrontsrv"
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
		workspace.NewExecCommandTool(mgr, accessor),
		workspace.NewScreenshotPageTool(mgr, accessor),
	}
}

// WorkspaceToolsWithPublish returns workspace tools including the publish tool.
// This requires the storefront service for publishing content to the live store.
// Returns nil if mgr or accessor is nil.
func WorkspaceToolsWithPublish(mgr containerx.Manager, accessor workspace.ContainerAccessor, sf *storefrontsrv.Service, tenantID kernel.TenantID) []Tool {
	if mgr == nil || accessor == nil {
		return nil
	}
	base := []Tool{
		workspace.NewWriteFileTool(mgr, accessor),
		workspace.NewReadFileTool(mgr, accessor),
		workspace.NewListFilesTool(mgr, accessor),
		workspace.NewPreviewURLTool(accessor),
		workspace.NewDeleteFileTool(mgr, accessor),
		workspace.NewExecCommandTool(mgr, accessor),
		workspace.NewScreenshotPageTool(mgr, accessor),
	}
	if sf != nil {
		base = append(base, workspace.NewPublishWorkspacePageTool(mgr, accessor, sf, tenantID))
	}
	return base
}
