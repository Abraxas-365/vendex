package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Abraxas-365/vendex/internal/containerx"
)

// ScreenshotPageTool captures a screenshot of a workspace page using headless Chromium.
// Requires the workspace container to have chromium installed (available in webdev preset).
type ScreenshotPageTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

// NewScreenshotPageTool creates a new ScreenshotPageTool.
func NewScreenshotPageTool(mgr containerx.Manager, accessor ContainerAccessor) *ScreenshotPageTool {
	return &ScreenshotPageTool{mgr: mgr, accessor: accessor}
}

func (t *ScreenshotPageTool) Name() string { return "screenshot_page" }

func (t *ScreenshotPageTool) Description() string {
	return "Take a screenshot of an HTML page in the workspace using headless Chromium. Saves the screenshot as a PNG in the workspace .screenshots directory. Useful for visually inspecting your work."
}

func (t *ScreenshotPageTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":      map[string]any{"type": "string", "description": "Path to the HTML file to screenshot (e.g. pages/landing.html)"},
			"width":     map[string]any{"type": "integer", "description": "Viewport width in pixels (default: 1280)", "default": 1280},
			"height":    map[string]any{"type": "integer", "description": "Viewport height in pixels (default: 800)", "default": 800},
			"full_page": map[string]any{"type": "boolean", "description": "Capture full scrollable page (default: false)", "default": false},
		},
		"required": []string{"path"},
	}
}

func (t *ScreenshotPageTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Path     string `json:"path"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		FullPage bool   `json:"full_page"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	safePath, err := safeguardPath(req.Path)
	if err != nil {
		return err.Error(), nil
	}

	if req.Width <= 0 {
		req.Width = 1280
	}
	if req.Height <= 0 {
		req.Height = 800
	}

	// Generate output path for the screenshot — replace slashes with underscores for a flat filename.
	screenshotFilename := strings.ReplaceAll(safePath, "/", "_") + ".png"
	outputPath := "/workspace/.screenshots/" + screenshotFilename

	// Ensure the screenshot directory exists.
	mkdirCmd := []string{"sh", "-c", "mkdir -p /workspace/.screenshots"}
	_, _ = t.mgr.Exec(ctx, t.accessor.ContainerID(), mkdirCmd)

	// The workspace HTTP server serves files on localhost:9091.
	pageURL := fmt.Sprintf("http://localhost:9091/%s", safePath)

	// Build the chromium screenshot command.
	// --full-page is a chromium CLI flag for full-page capture.
	fullPageFlag := ""
	if req.FullPage {
		fullPageFlag = "--full-page"
	}

	chromiumCmd := fmt.Sprintf(
		"chromium-browser --headless --no-sandbox --disable-gpu --screenshot=%s --window-size=%d,%d %s %s 2>/dev/null || "+
			"chromium --headless --no-sandbox --disable-gpu --screenshot=%s --window-size=%d,%d %s %s 2>/dev/null",
		outputPath, req.Width, req.Height, fullPageFlag, pageURL,
		outputPath, req.Width, req.Height, fullPageFlag, pageURL,
	)
	cmd := []string{"sh", "-c", chromiumCmd}

	_, err = t.mgr.Exec(ctx, t.accessor.ContainerID(), cmd)
	if err != nil {
		return fmt.Sprintf(
			"Screenshot failed. Ensure chromium or chromium-browser is available in the workspace container.\nError: %s\nTip: The page is accessible at %s",
			err.Error(), pageURL,
		), nil
	}

	// Verify the screenshot was actually created.
	checkCmd := []string{"sh", "-c", fmt.Sprintf("test -f %s && echo OK", outputPath)}
	checkOutput, _ := t.mgr.Exec(ctx, t.accessor.ContainerID(), checkCmd)
	if !strings.Contains(string(checkOutput), "OK") {
		return "Screenshot command ran but output file was not created. The page may have issues loading.", nil
	}

	// Read the file size for reporting.
	sizeCmd := []string{"sh", "-c", fmt.Sprintf("wc -c < %s", outputPath)}
	sizeOutput, _ := t.mgr.Exec(ctx, t.accessor.ContainerID(), sizeCmd)
	size := strings.TrimSpace(string(sizeOutput))
	sizeBytes := parseInt(size)

	relOutput := ".screenshots/" + screenshotFilename
	previewURL := t.accessor.PreviewBaseURL() + "/" + relOutput

	// For small screenshots, read and return base64 for inline vision capability.
	if sizeBytes > 0 && sizeBytes < 500000 {
		b64Cmd := []string{"sh", "-c", fmt.Sprintf("base64 %s", outputPath)}
		b64Output, b64Err := t.mgr.Exec(ctx, t.accessor.ContainerID(), b64Cmd)
		if b64Err == nil {
			encoded := strings.TrimSpace(string(b64Output))
			return fmt.Sprintf(
				"Screenshot saved to %s (%s bytes)\nPreview URL: %s\n\n[IMAGE:data:image/png;base64,%s]",
				relOutput, size, previewURL, encoded,
			), nil
		}
	}

	return fmt.Sprintf("Screenshot saved to %s (%s bytes)\nPreview URL: %s", relOutput, size, previewURL), nil
}

// parseInt parses a trimmed string as an integer, returning 0 on failure.
func parseInt(s string) int {
	var n int
	fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n
}
