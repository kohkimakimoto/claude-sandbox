package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GetWorkdir returns the working directory for sandbox execution.
// It uses CLAUDE_SANDBOX_WORKDIR if set, otherwise falls back to the current directory.
func GetWorkdir() string {
	if v := os.Getenv("CLAUDE_SANDBOX_WORKDIR"); v != "" {
		return v
	}
	wd, _ := os.Getwd()
	return wd
}

// GetClaudeBin returns the path to the claude binary.
// It checks CLAUDE_SANDBOX_CLAUDE_BIN, then searches PATH, then falls back to ~/.claude/local/claude.
func GetClaudeBin() string {
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

// SocketPath returns the path for the daemon's Unix Domain Socket.
func SocketPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("claude-sandbox-unboxexec-%d.sock", os.Getpid()))
}
