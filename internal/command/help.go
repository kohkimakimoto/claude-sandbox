package command

// RootHelpTemplate is the help template for the root command.
const RootHelpTemplate = `Usage: claude-sandbox [<command>]|[claude [<args of claude command...>]]

{{template "descriptionTemplate" .}}

Builtin commands:{{template "visibleCommandCategoryTemplate" .}}

Configuration:
   claude-sandbox looks for config files in the following order:

   1. $HOME/.claude/sandbox.toml (user-level)
   2. .claude/sandbox.toml (project-level)
   3. .claude/sandbox.local.toml (local overrides, gitignore-friendly)

   See: https://github.com/kohkimakimoto/claude-sandbox#configuration-file

Example Usage:
   # Create project-specific config file
   $ claude-sandbox init

   # Create local override config file (not for version control)
   $ claude-sandbox init-local

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
Commit: {{ index (ExtraInfo) "CommitHash" }}
{{template "copyrightTemplate" .}}
`

// HelpTemplate is the help template for subcommands.
const HelpTemplate = `Usage: {{template "usageTemplate" .}}

{{template "helpNameTemplate" .}}

Options:{{template "visibleFlagTemplate" .}}
`
