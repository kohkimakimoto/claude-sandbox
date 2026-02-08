package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
// It starts an internal daemon for sandbox-external command execution,
// then runs sandbox-exec as a child process.
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

	// Start the daemon for sandbox-external command execution
	sockPath := socketPath()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := startDaemon(ctx, sockPath); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Build sandbox-exec command arguments
	sandboxExecArgs := []string{
		"-D", "WORKDIR=" + workdir,
		"-D", "HOME=" + home,
		"-f", profilePath,
		claudeBin,
	}
	sandboxExecArgs = append(sandboxExecArgs, args...)

	// Run sandbox-exec as a child process
	eCmd := exec.CommandContext(ctx, "sandbox-exec", sandboxExecArgs...)
	eCmd.Env = append(os.Environ(), "CLAUDE_SANDBOX_UNBOX_EXEC_SOCK="+sockPath)
	eCmd.Stdin = os.Stdin
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr

	err = eCmd.Run()
	cancel() // shut down daemon

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return cli.Exit("", exitErr.ExitCode())
		}
		return err
	}

	return nil
}
