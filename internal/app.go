package internal

import (
	"context"

	"github.com/urfave/cli/v3"
)

func Run(args []string) error {
	return newApp().Run(context.Background(), args)
}

func newApp() *cli.Command {
	app := &cli.Command{
		Name:                          "claude-sandbox",
		HideVersion:                   true,
		Version:                       Version,
		Copyright:                     "Copyright (c) Kohki Makimoto",
		SkipFlagParsing:               true,
		CustomRootCommandHelpTemplate: rootHelpTemplate,
	}

	app.Commands = []*cli.Command{
		InitCommand,
		InitGlobalCommand,
		ProfileCommand,
		ClaudeCommand,
	}

	app.Action = func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Present() {
			first := cmd.Args().First()
			if first == "help" || first == "--help" || first == "-h" {
				return cli.ShowAppHelp(cmd)
			}
			// If args are present and not a builtin command, run claude with all args
			return runClaudeAction(ctx, cmd, cmd.Args().Slice())
		}

		// No args: run claude without arguments
		return runClaudeAction(ctx, cmd, nil)
	}

	return app
}
