package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kohkimakimoto/claude-sandbox/internal/sandbox"
	"github.com/urfave/cli/v3"
)

var InitCommand = &cli.Command{
	Name:               "init",
	Usage:              "Create .claude/sandbox.sb file if it doesn't exist",
	CustomHelpTemplate: HelpTemplate,
	Action:             initAction,
}

func initAction(ctx context.Context, cmd *cli.Command) error {
	workdir := sandbox.GetWorkdir()
	sandboxFile := filepath.Join(workdir, ".claude", "sandbox.sb")

	if _, err := os.Stat(sandboxFile); err == nil {
		return fmt.Errorf("sandbox profile file already exists: %s", sandboxFile)
	}

	// Create .claude directory if it doesn't exist
	if err := os.MkdirAll(filepath.Join(workdir, ".claude"), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(sandboxFile, []byte(sandbox.ProjectProfileTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	fmt.Fprintf(cmd.Writer, "Created sandbox profile file: %s\n", sandboxFile)
	return nil
}
