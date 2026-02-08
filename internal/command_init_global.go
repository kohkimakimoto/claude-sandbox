package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var InitGlobalCommand = &cli.Command{
	Name:               "init-global",
	Usage:              "Create $HOME/.claude/sandbox.sb file if it doesn't exist",
	CustomHelpTemplate: helpTemplate,
	Action:             initGlobalAction,
}

func initGlobalAction(ctx context.Context, cmd *cli.Command) error {
	home, _ := os.UserHomeDir()
	sandboxFile := filepath.Join(home, ".claude", "sandbox.sb")

	if _, err := os.Stat(sandboxFile); err == nil {
		return fmt.Errorf("global sandbox profile file already exists: %s", sandboxFile)
	}

	// Create .claude directory if it doesn't exist
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(sandboxFile, []byte(globalProfileTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	fmt.Fprintf(cmd.Writer, "Created global sandbox profile file: %s\n", sandboxFile)
	return nil
}
