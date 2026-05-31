// Package workspace provides agent tools that execute inside a Docker container workspace.
// These tools enable the agent to read, write, list, and delete files in a sandboxed
// environment, as well as get preview URLs and publish content to the store.
package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Abraxas-365/vendex/internal/containerx"
)

// ContainerAccessor provides access to the running workspace container for a session.
// The agent session manager implements this interface.
type ContainerAccessor interface {
	// ContainerID returns the Docker container ID for the current workspace.
	ContainerID() containerx.ID
	// PreviewBaseURL returns the base URL where workspace files are publicly accessible.
	PreviewBaseURL() string
}

// WriteFileTool writes content to a file inside the workspace container.
type WriteFileTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

func NewWriteFileTool(mgr containerx.Manager, accessor ContainerAccessor) *WriteFileTool {
	return &WriteFileTool{mgr: mgr, accessor: accessor}
}

func (t *WriteFileTool) Name() string        { return "write_file" }
func (t *WriteFileTool) Description() string { return "Write content to a file in the workspace" }
func (t *WriteFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    map[string]any{"type": "string", "description": "Relative path in workspace (e.g. pages/landing.html)"},
			"content": map[string]any{"type": "string", "description": "File content to write"},
		},
		"required": []string{"path", "content"},
	}
}

func (t *WriteFileTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	// Create parent directory and write file via exec
	dir := filepath.Dir("/workspace/" + safePath)
	cmd := []string{"sh", "-c", fmt.Sprintf("mkdir -p %q && cat > %q", dir, "/workspace/"+safePath)}

	// Write content via a temporary approach using exec with stdin piping.
	// Since Docker exec doesn't natively support stdin content easily,
	// we use a printf-based approach for the content.
	escaped := strings.ReplaceAll(req.Content, "'", "'\\''")
	writeCmd := []string{"sh", "-c", fmt.Sprintf("mkdir -p %q && printf '%%s' '%s' > %q", dir, escaped, "/workspace/"+safePath)}

	_, err = t.mgr.Exec(ctx, t.accessor.ContainerID(), writeCmd)
	if err != nil {
		return "failed to write file: " + err.Error(), nil
	}

	_ = cmd // suppresses unused warning for earlier approach
	return fmt.Sprintf("Written %d bytes to %s", len(req.Content), safePath), nil
}

// ReadFileTool reads a file from the workspace container.
type ReadFileTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

func NewReadFileTool(mgr containerx.Manager, accessor ContainerAccessor) *ReadFileTool {
	return &ReadFileTool{mgr: mgr, accessor: accessor}
}

func (t *ReadFileTool) Name() string        { return "read_file" }
func (t *ReadFileTool) Description() string { return "Read a file from the workspace" }
func (t *ReadFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "Relative path in workspace"},
		},
		"required": []string{"path"},
	}
}

func (t *ReadFileTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	output, err := t.mgr.Exec(ctx, t.accessor.ContainerID(), []string{"cat", "/workspace/" + safePath})
	if err != nil {
		return "file not found or cannot be read: " + safePath, nil
	}

	return string(output), nil
}

// ListFilesTool lists files in a workspace directory.
type ListFilesTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

func NewListFilesTool(mgr containerx.Manager, accessor ContainerAccessor) *ListFilesTool {
	return &ListFilesTool{mgr: mgr, accessor: accessor}
}

func (t *ListFilesTool) Name() string        { return "list_files" }
func (t *ListFilesTool) Description() string { return "List files in the workspace directory" }
func (t *ListFilesTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "Directory path relative to workspace root", "default": "."},
		},
	}
}

func (t *ListFilesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	if req.Path == "" {
		req.Path = "."
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	output, err := t.mgr.Exec(ctx, t.accessor.ContainerID(), []string{"find", "/workspace/" + safePath, "-maxdepth", "2", "-type", "f"})
	if err != nil {
		return "directory not found: " + safePath, nil
	}

	// Clean up paths to be relative to /workspace/
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		rel := strings.TrimPrefix(line, "/workspace/")
		result = append(result, rel)
	}

	if len(result) == 0 {
		return "No files found in " + safePath, nil
	}

	return strings.Join(result, "\n"), nil
}

// DeleteFileTool deletes a file from the workspace.
type DeleteFileTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

func NewDeleteFileTool(mgr containerx.Manager, accessor ContainerAccessor) *DeleteFileTool {
	return &DeleteFileTool{mgr: mgr, accessor: accessor}
}

func (t *DeleteFileTool) Name() string        { return "delete_file" }
func (t *DeleteFileTool) Description() string { return "Delete a file from the workspace" }
func (t *DeleteFileTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "Relative path in workspace"},
		},
		"required": []string{"path"},
	}
}

func (t *DeleteFileTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	_, err = t.mgr.Exec(ctx, t.accessor.ContainerID(), []string{"rm", "-f", "/workspace/" + safePath})
	if err != nil {
		return "failed to delete: " + err.Error(), nil
	}

	return "Deleted " + safePath, nil
}

// PreviewURLTool returns the live preview URL for a workspace file.
type PreviewURLTool struct {
	accessor ContainerAccessor
}

func NewPreviewURLTool(accessor ContainerAccessor) *PreviewURLTool {
	return &PreviewURLTool{accessor: accessor}
}

func (t *PreviewURLTool) Name() string        { return "preview_url" }
func (t *PreviewURLTool) Description() string { return "Get the live preview URL for a page in the workspace" }
func (t *PreviewURLTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "Relative path to the HTML file (e.g. pages/landing.html)"},
		},
		"required": []string{"path"},
	}
}

func (t *PreviewURLTool) Execute(_ context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	url := t.accessor.PreviewBaseURL() + "/" + safePath
	return fmt.Sprintf("Preview URL: %s\n\nThe merchant can open this URL to see the live preview of the page.", url), nil
}

// safeguardPath validates and cleans a relative path to prevent directory traversal.
func safeguardPath(path string) (string, error) {
	// Clean the path
	cleaned := filepath.Clean(path)
	// Reject absolute paths
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be relative, got: %s", path)
	}
	// Reject directory traversal
	if strings.HasPrefix(cleaned, "..") || strings.Contains(cleaned, "/../") {
		return "", fmt.Errorf("path cannot escape workspace: %s", path)
	}
	return cleaned, nil
}
