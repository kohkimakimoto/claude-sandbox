package command

// RootHelpTemplate is the help template for the root command.
const RootHelpTemplate = `Usage: claude-sandbox [<command>]|[claude [<args of claude command...>]]

A wrapper around the claude command to run it in a sandboxed environment.

Builtin commands:{{template "visibleCommandCategoryTemplate" .}}

Configuration:
   All settings are managed through a single TOML configuration file.
   claude-sandbox looks for config files in the following order:

   1. .claude/sandbox.toml (project-specific config)
   2. $HOME/.claude/sandbox.toml (global config)

   The project-specific config takes precedence over the global config.
   If neither exists, built-in defaults are used.

   The [sandbox] section can contain:
   - profile: sandbox-exec profile content (if not set, built-in default is used)
   - workdir: override working directory
   - claude_bin: override path to claude binary

   The [unboxexec] section can contain:
   - allowed_commands: regex patterns for commands allowed via unboxexec

Example Usage:
   # Create project-specific config file
   $ claude-sandbox init

   # Create user config file
   $ claude-sandbox init-user

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
