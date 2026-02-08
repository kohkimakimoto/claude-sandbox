package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/urfave/cli/v3"
)

var ClaudeCommand = &cli.Command{
	Name:               "claude",
	Usage:              "Run the claude command in a sandboxed environment",
	SkipFlagParsing:    true,
	CustomHelpTemplate: helpTemplate,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return runClaudeAction(ctx, cmd, cmd.Args().Slice())
	},
}

// runClaudeAction executes claude inside a macOS sandbox using sandbox-exec.
func runClaudeAction(ctx context.Context, cmd *cli.Command, args []string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	profilePath, cleanup, err := buildProfile()
	if err != nil {
		return err
	}
	defer cleanup()

	workdir := getWorkdir()
	home, _ := os.UserHomeDir()
	claudeBin := getClaudeBin()

	// Build sandbox-exec command arguments:
	// sandbox-exec -D WORKDIR=<workdir> -D HOME=<home> -f <profile> <claude-bin> [args...]
	sandboxExecArgs := []string{
		"sandbox-exec",
		"-D", "WORKDIR=" + workdir,
		"-D", "HOME=" + home,
		"-f", profilePath,
		claudeBin,
	}
	sandboxExecArgs = append(sandboxExecArgs, args...)

	// Use syscall.Exec to replace the current process, matching the shell script's `exec` behavior.
	sandboxExecPath, err := exec.LookPath("sandbox-exec")
	if err != nil {
		return fmt.Errorf("sandbox-exec not found: %w", err)
	}

	return syscall.Exec(sandboxExecPath, sandboxExecArgs, os.Environ())
}
