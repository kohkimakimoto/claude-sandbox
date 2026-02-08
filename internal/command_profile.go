package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var ProfileCommand = &cli.Command{
	Name:               "profile",
	Usage:              "Print evaluated profile and exit",
	CustomHelpTemplate: helpTemplate,
	Action:             profileAction,
}

func profileAction(ctx context.Context, cmd *cli.Command) error {
	profilePath, cleanup, err := buildProfile()
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
