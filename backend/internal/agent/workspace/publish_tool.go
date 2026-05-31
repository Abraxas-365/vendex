package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Abraxas-365/vendex/internal/containerx"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront/storefrontsrv"
)

// PublishWorkspacePageTool reads HTML+CSS from workspace and creates/updates a storefront page.
type PublishWorkspacePageTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
	sf       *storefrontsrv.Service
	tenantID kernel.TenantID
}

// NewPublishWorkspacePageTool creates a new PublishWorkspacePageTool.
func NewPublishWorkspacePageTool(mgr containerx.Manager, accessor ContainerAccessor, sf *storefrontsrv.Service, tenantID kernel.TenantID) *PublishWorkspacePageTool {
	return &PublishWorkspacePageTool{mgr: mgr, accessor: accessor, sf: sf, tenantID: tenantID}
}

func (t *PublishWorkspacePageTool) Name() string { return "publish_workspace_page" }

func (t *PublishWorkspacePageTool) Description() string {
	return "Publish an HTML page from the workspace to the live storefront. Reads the file content and creates or updates a storefront page with the given slug and title. Agent-created pages land in pending_review and require admin approval before going live."
}

func (t *PublishWorkspacePageTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     map[string]any{"type": "string", "description": "Path to the HTML file in workspace (e.g. pages/landing.html)"},
			"slug":     map[string]any{"type": "string", "description": "URL slug for the published page (e.g. 'summer-sale')"},
			"title":    map[string]any{"type": "string", "description": "Page title shown in browser tab and admin"},
			"css_path": map[string]any{"type": "string", "description": "Optional path to a separate CSS file to include"},
		},
		"required": []string{"path", "slug", "title"},
	}
}

func (t *PublishWorkspacePageTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path    string `json:"path"`
		Slug    string `json:"slug"`
		Title   string `json:"title"`
		CSSPath string `json:"css_path"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	if req.Path == "" || req.Slug == "" || req.Title == "" {
		return "path, slug, and title are required", nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	// Read HTML content from workspace
	htmlBytes, err := t.mgr.Exec(ctx, t.accessor.ContainerID(), []string{"cat", "/workspace/" + safePath})
	if err != nil {
		return "failed to read HTML file: " + safePath, nil
	}
	html := strings.TrimSpace(string(htmlBytes))

	// Optionally read CSS
	css := ""
	if req.CSSPath != "" {
		cssPath, cerr := safeguardPath(req.CSSPath)
		if cerr != nil {
			return cerr.Error(), nil
		}
		cssBytes, cerr := t.mgr.Exec(ctx, t.accessor.ContainerID(), []string{"cat", "/workspace/" + cssPath})
		if cerr != nil {
			return "failed to read CSS file: " + cssPath, nil
		}
		css = strings.TrimSpace(string(cssBytes))
	}

	// Try to find an existing page by slug.
	existingPage, err := t.sf.GetPageBySlug(ctx, t.tenantID, req.Slug)
	if err != nil {
		if !errx.IsNotFound(err) {
			return "failed to look up existing page: " + err.Error(), nil
		}
		// Page does not exist — create it.
		page, cerr := t.sf.CreatePage(ctx, storefrontsrv.CreatePageInput{
			TenantID:  t.tenantID,
			Slug:      req.Slug,
			Title:     req.Title,
			HTML:      html,
			CSS:       css,
			CreatedBy: "agent",
			ByAgent:   true,
		})
		if cerr != nil {
			return "failed to create page: " + cerr.Error(), nil
		}
		return fmt.Sprintf(
			"Published new page '%s' at slug /%s (page ID: %s)\nStatus: %s — an admin must approve before it goes live.",
			page.Title, page.Slug, page.ID, page.Status,
		), nil
	}

	// Page exists — update it.
	page, err := t.sf.UpdatePage(ctx, storefrontsrv.UpdatePageInput{
		TenantID: t.tenantID,
		ID:       existingPage.ID,
		Title:    &req.Title,
		HTML:     &html,
		CSS:      &css,
		EditedBy: "agent",
		Comment:  "updated from workspace",
	})
	if err != nil {
		return "failed to update page: " + err.Error(), nil
	}
	return fmt.Sprintf(
		"Updated page '%s' at slug /%s (page ID: %s)\nStatus: %s — version %d saved.",
		page.Title, page.Slug, page.ID, page.Status, page.Version,
	), nil
}
