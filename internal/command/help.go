package command

// RootHelpTemplate is the help template for the root command.
const RootHelpTemplate = `Usage: claude-sandbox [<command>]|[claude [<args of claude command...>]]

A wrapper around the claude command to run it in a sandboxed environment.

Builtin commands:{{template "visibleCommandCategoryTemplate" .}}

Profile:
   The profile file defines sandbox environment.
   This is a sandbox-exec (Apple Seatbelt) profile.
   claude-sandbox looks for profile files in the following order:

   1. .claude/sandbox.sb (project-specific profile)
   2. $HOME/.claude/sandbox.sb (global profile)

   The project-specific profile takes precedence over the global profile.
   If neither exists, claude-sandbox will run its built-in default profile.

   You can use parameters in the profile file using (param "NAME") syntax.
   The following parameters are provided by claude-sandbox:

   - WORKDIR: The current working directory where claude-sandbox is executed.
   - HOME: The user's home directory.

Example Usage:
   # Create project-specific sandbox profile
   $ claude-sandbox init

   # Create global sandbox profile
   $ claude-sandbox init-global

   # Print the evaluated sandbox profile
   $ claude-sandbox profile

   # Run Claude Code in a sandboxed environment
   $ claude-sandbox claude

   # Run Claude Code with arguments in a sandboxed environment
   $ claude-sandbox claude --dangerously-skip-permissions

   # You can also run Claude Code without the 'claude' command prefix.
   $ claude-sandbox
   $ claude-sandbox --dangerously-skip-permissions

   Commands or options that conflict with claude-sandbox can be used with the claude command prefix.
   For example, the following command shows the claude help, not the claude-sandbox help.
   $ claude-sandbox claude -h

Version: {{ .Version }}
{{template "copyrightTemplate" .}}
`

// HelpTemplate is the help template for subcommands.
const HelpTemplate = `Usage: {{template "usageTemplate" .}}

{{template "helpNameTemplate" .}}

Options:{{template "visibleFlagTemplate" .}}
`
