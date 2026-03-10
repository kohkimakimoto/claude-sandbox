package command

import (
	"context"
	"fmt"
	"os"

	"github.com/kohkimakimoto/claude-sandbox/v2/internal/config"
	"github.com/kohkimakimoto/claude-sandbox/v2/internal/sandbox"
	"github.com/urfave/cli/v3"
)

func NewProfileCommand() *cli.Command {
	return &cli.Command{
		Name:               "profile",
		Usage:              "Print evaluated profile and exit",
		CustomHelpTemplate: HelpTemplate,
		Action:             profileAction,
	}
}

func profileAction(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	profilePath, cleanup, err := sandbox.BuildProfile(cfg.Sandbox.Profile)
	if err != nil {
		return err
	}
	defer cleanup()

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	_, err = cmd.Writer.Write(content)
	return err
}
