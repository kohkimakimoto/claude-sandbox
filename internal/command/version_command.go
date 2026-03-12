package command

import (
	"context"
	"fmt"

	"github.com/kohkimakimoto/claude-sandbox/v2/internal/version"
	"github.com/urfave/cli/v3"
)

func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:               "version",
		Usage:              "Print version and exit",
		CustomHelpTemplate: HelpTemplate,
		Action:             versionAction,
	}
}

func versionAction(ctx context.Context, cmd *cli.Command) error {
	_, err := fmt.Fprintf(cmd.Root().Writer, "%s version %s (commit: %s)\n", cmd.Root().Name, version.Version, version.CommitHash)
	return err
}
