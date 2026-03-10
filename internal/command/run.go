package command

import (
	"context"

	"github.com/kohkimakimoto/claude-sandbox/v2/internal/config"
	"github.com/kohkimakimoto/claude-sandbox/v2/internal/version"
	"github.com/urfave/cli/v3"
)

func Run(args []string) error {
	cfg, err := config.LoadMerged()
	if err != nil {
		return err
	}

	return newApp(cfg).Run(context.Background(), args)
}

func newApp(cfg *config.Config) *cli.Command {
	app := &cli.Command{
		Name:                          "claude-sandbox",
		HideVersion:                   true,
		Version:                       version.Version,
		ExtraInfo:                     func() map[string]string { return map[string]string{"CommitHash": version.CommitHash} },
		Copyright:                     "Copyright (c) Kohki Makimoto",
		SkipFlagParsing:               true,
		CustomRootCommandHelpTemplate: RootHelpTemplate,
	}

	app.Commands = []*cli.Command{
		NewInitCommand(),
		NewInitLocalCommand(),
		NewInitUserCommand(),
		NewInitGlobalCommand(),
		NewProfileCommand(),
		NewClaudeCommand(),
		NewUnboxexecCommand(),
	}

	app.Action = func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Present() {
			first := cmd.Args().First()
			if first == "help" || first == "--help" || first == "-h" {
				return cli.ShowAppHelp(cmd)
			}
			// If args are present and not a builtin command, run claude with all args
			return RunClaudeAction(ctx, cmd, cmd.Args().Slice(), cfg)
		}

		// No args: run claude without arguments
		return RunClaudeAction(ctx, cmd, nil, cfg)
	}

	return app
}
