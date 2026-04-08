package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ClaudeProvider uses the Claude CLI subprocess.
type ClaudeProvider struct {
	Model string
}

func (c *ClaudeProvider) Name() string { return "claude" }

func (c *ClaudeProvider) Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return "", fmt.Errorf("claude CLI not found on PATH: %w", err)
	}

	args := []string{"-p", "--output-format", "text"}
	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}
	if c.Model != "" {
		args = append(args, "--model", c.Model)
	}
	args = append(args, userMessage)

	cmd := exec.CommandContext(ctx, "claude", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude CLI error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("claude CLI error: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
