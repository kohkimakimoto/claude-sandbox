package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/kohkimakimoto/claude-sandbox/internal/config"
	"github.com/kohkimakimoto/claude-sandbox/internal/sandbox"
	"github.com/kohkimakimoto/claude-sandbox/internal/unboxexec"
	"github.com/urfave/cli/v3"
)

var ClaudeCommand = &cli.Command{
	Name:               "claude",
	Usage:              "Run the claude command in a sandboxed environment",
	SkipFlagParsing:    true,
	CustomHelpTemplate: HelpTemplate,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		// When invoked as a subcommand, load config
		cfg, err := config.Load(config.ResolveConfigPath())
		if err != nil {
			return err
		}
		return RunClaudeAction(ctx, cmd, cmd.Args().Slice(), cfg)
	},
}

// RunClaudeAction executes claude inside a macOS sandbox using sandbox-exec.
// It starts an internal daemon for sandbox-external command execution,
// then runs sandbox-exec as a child process.
func RunClaudeAction(ctx context.Context, cmd *cli.Command, args []string, cfg *config.Config) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	profilePath, cleanup, err := sandbox.BuildProfile()
	if err != nil {
		return err
	}
	defer cleanup()

	workdir := sandbox.GetWorkdir()
	home, _ := os.UserHomeDir()
	claudeBin := sandbox.GetClaudeBin()

	// Compile allowed command patterns
	allowedCommands, err := config.CompileAllowedCommands(cfg.Unboxexec.AllowedCommands)
	if err != nil {
		return fmt.Errorf("failed to compile allowed_commands: %w", err)
	}

	// Start the daemon for sandbox-external command execution
	sockPath := sandbox.SocketPath()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := unboxexec.StartDaemon(ctx, sockPath, allowedCommands); err != nil {
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
	eCmd.Env = append(os.Environ(), "CLAUDE_SANDBOX=1", "CLAUDE_SANDBOX_UNBOXEXEC_SOCK="+sockPath)
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
