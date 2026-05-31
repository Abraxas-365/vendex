package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Abraxas-365/vendex/internal/containerx"
)

// ExecCommandTool runs a shell command inside the workspace container.
type ExecCommandTool struct {
	mgr      containerx.Manager
	accessor ContainerAccessor
}

// NewExecCommandTool creates a new ExecCommandTool.
func NewExecCommandTool(mgr containerx.Manager, accessor ContainerAccessor) *ExecCommandTool {
	return &ExecCommandTool{mgr: mgr, accessor: accessor}
}

func (t *ExecCommandTool) Name() string { return "exec_command" }

func (t *ExecCommandTool) Description() string {
	return "Execute a shell command inside the workspace container. Use for running build tools, linters, or file operations. Commands run in /workspace as working directory."
}

func (t *ExecCommandTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command":         map[string]any{"type": "string", "description": "Shell command to execute (e.g. 'ls -la pages/' or 'cat pages/index.html | wc -l')"},
			"timeout_seconds": map[string]any{"type": "integer", "description": "Max execution time in seconds (default: 30, max: 120)", "default": 30},
		},
		"required": []string{"command"},
	}
}

func (t *ExecCommandTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Command string `json:"command"`
		Timeout int    `json:"timeout_seconds"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}

	if req.Command == "" {
		return "command is required", nil
	}

	// Enforce timeout bounds
	if req.Timeout <= 0 {
		req.Timeout = 30
	}
	if req.Timeout > 120 {
		req.Timeout = 120
	}

	// Block dangerous commands
	if isDangerous(req.Command) {
		return "command rejected: potentially dangerous operation", nil
	}

	execCtx, cancel := context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Second)
	defer cancel()

	cmd := []string{"sh", "-c", fmt.Sprintf("cd /workspace && %s", req.Command)}
	output, err := t.mgr.Exec(execCtx, t.accessor.ContainerID(), cmd)
	if err != nil {
		return fmt.Sprintf("command failed: %s\nError: %s", req.Command, err.Error()), nil
	}

	result := strings.TrimSpace(string(output))
	if len(result) > 10000 {
		result = result[:10000] + "\n... (output truncated at 10000 chars)"
	}
	if result == "" {
		result = "(no output)"
	}

	return result, nil
}

// isDangerous checks if a command could be harmful to the container or host.
func isDangerous(cmd string) bool {
	lower := strings.ToLower(cmd)
	dangerous := []string{
		"rm -rf /",
		"mkfs",
		"dd if=",
		"> /dev/",
		"chmod 777 /",
		"curl | sh",
		"wget | sh",
		"curl|sh",
		"wget|sh",
	}
	for _, d := range dangerous {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}
