package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// getWorkdir returns the working directory for sandbox execution.
// It uses CLAUDE_SANDBOX_WORKDIR if set, otherwise falls back to the current directory.
func getWorkdir() string {
	if v := os.Getenv("CLAUDE_SANDBOX_WORKDIR"); v != "" {
		return v
	}
	wd, _ := os.Getwd()
	return wd
}

// getClaudeBin returns the path to the claude binary.
// It checks CLAUDE_SANDBOX_CLAUDE_BIN, then searches PATH, then falls back to ~/.claude/local/claude.
func getClaudeBin() string {
	if v := os.Getenv("CLAUDE_SANDBOX_CLAUDE_BIN"); v != "" {
		return v
	}

	if p, err := exec.LookPath("claude"); err == nil {
		return p
	}

	home, _ := os.UserHomeDir()
	localClaude := filepath.Join(home, ".claude", "local", "claude")
	if _, err := os.Stat(localClaude); err == nil {
		return localClaude
	}

	return "claude"
}

// socketPath returns the path for the daemon's Unix Domain Socket.
func socketPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("claude-sandbox-unbox-exec-%d.sock", os.Getpid()))
}
