# claude-sandbox

A wrapper around Claude Code (`claude` command) to run it in a sandboxed environment using macOS's `sandbox-exec`.

> [!NOTE]
> v2 is a full rewrite from a shell script to Go. This enables single binary distribution and integrated process management.
>
> **v2 is currently under active development and may introduce breaking changes without notice.**
>
> Key changes from v1:
>
> - **Reimplemented in Go** — No more shell script. Everything is a single compiled binary.
> - **Sandbox-external command execution (unboxexec)** — A built-in daemon that lets Claude Code selectively bypass the sandbox for tools that require it (e.g. Playwright). Commands are restricted via a TOML configuration file.
> - **Configuration file support** — Project-specific or global `.claude/sandbox.toml` to control allowed unboxexec commands.

## Installation

Build from source:

```bash
git clone https://github.com/kohkimakimoto/claude-sandbox.git
cd claude-sandbox
make build
# Binary is at .dev/build/dev/claude-sandbox
```

## Usage

`claude-sandbox` can be used as a drop-in replacement for the `claude` command, but runs in a sandboxed environment that restricts file system write access.

```bash
# Instead of: claude
claude-sandbox

# Instead of: claude --dangerously-skip-permissions
claude-sandbox --dangerously-skip-permissions
```

You can also use the explicit `claude` subcommand. These commands are equivalent to the above:

```bash
claude-sandbox claude
claude-sandbox claude --dangerously-skip-permissions
```

Commands or options that conflict with claude-sandbox's own can be passed using the `claude` subcommand prefix. For example, the following shows the claude help, not the claude-sandbox help:

```bash
claude-sandbox claude -h
```

### Viewing the Sandbox Profile

By default, `claude-sandbox` uses a built-in sandbox profile that restricts file system write access to the current working directory, Claude Code configuration, and temporary directories.
You can view the actual profile being used:

```bash
claude-sandbox profile
```

Example output:

```scheme
;; This is a default built-in sandbox profile for claude-sandbox.
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code
    (regex (string-append "^" (param "HOME") "/\\.claude"))

    ;; Keychain access for Claude Code credentials
    (subpath (string-append (param "HOME") "/Library/Keychains"))

    ;; Temporary directories and files
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")

    ;; Home directory
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))
    (subpath (string-append (param "HOME") "/Library/Caches"))
    (regex (string-append "^" (param "HOME") "/\\.viminfo"))

    ;; devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
    (literal "/dev/dtracehelper")
    (regex #"^/dev/tty*")
)
```

The sandbox uses macOS's `sandbox-exec` (Apple Seatbelt) technology. Even if Claude Code tried to execute a command like `rm -rf /usr/bin` or modify system configuration files, the sandbox would block these operations.

## Configuring Sandbox Profiles

To customize the sandbox environment, you can create a sandbox profile. There are two types of profiles:

### Project-Specific Profile

Create a project-specific sandbox profile that applies only to the current project:

```bash
claude-sandbox init
```

This creates `.claude/sandbox.sb` in your current directory.
You can then edit this file to customize the sandbox permissions for your project.

### Global Profile

Create a global sandbox profile that applies to all projects:

```bash
claude-sandbox init-global
```

This creates `~/.claude/sandbox.sb`.

**Profile Priority**: Project-specific profiles take precedence over global profiles. If neither exists, a built-in default profile is used.

### Parameters

The profile uses parameters that are passed from claude-sandbox automatically:

- `WORKDIR`: The current working directory where claude-sandbox is executed
- `HOME`: The user's home directory

You can use these parameters in your sandbox profile like this:

```scheme
(allow file-write*
    (subpath (param "WORKDIR"))
    (subpath (string-append (param "HOME") "/.claude"))
)
```

## Sandbox-External Command Execution

Some tools (e.g. Playwright) cannot run inside the macOS sandbox because they use their own sandboxing mechanisms, which conflict with the nested sandbox environment.

`claude-sandbox` includes a built-in mechanism called **unboxexec** that allows commands to be executed outside the sandbox. When `claude-sandbox` starts, it launches an internal daemon that accepts command execution requests from inside the sandbox.

### The `unboxexec` Subcommand

The `claude-sandbox unboxexec` subcommand is used from inside the sandbox to execute commands outside of it.

```bash
claude-sandbox unboxexec [options] -- <command> [args...]
```

#### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-C` | Specify the working directory for the command |
| `--timeout` | `-t` | Timeout in seconds (default: 60) |
| `--env` | `-e` | Environment variable in `KEY=VALUE` format (can be specified multiple times) |

#### Examples

```bash
# Execute a command outside the sandbox
claude-sandbox unboxexec -- echo "hello from outside"

# Execute with a specified working directory
claude-sandbox unboxexec --dir /tmp -- ls -la

# Execute with an extended timeout
claude-sandbox unboxexec --timeout 300 -- long-running-command

# Execute with environment variables
claude-sandbox unboxexec --env API_KEY=secret --env DEBUG=1 -- my-command
```

### Command Restrictions

By default, all commands executed via `unboxexec` are **rejected** unless explicitly allowed by a configuration file. See [Configuration File](#configuration-file) for details.

## Configuration File

`claude-sandbox` supports a TOML configuration file for controlling behavior. The configuration file is resolved in the following order:

1. `.claude/sandbox.toml` in the current working directory (project-specific)
2. `~/.claude/sandbox.toml` (global)

The project-specific configuration takes precedence over the global configuration. If neither file exists, the default configuration is used (all `unboxexec` commands are rejected).

### Example

```toml
# .claude/sandbox.toml or ~/.claude/sandbox.toml

[unboxexec]
# Regex patterns for allowed commands.
# The command and its arguments are joined by spaces, and the resulting string
# is matched against each pattern. If any pattern matches, the command is allowed.
# If empty or not configured, all commands are rejected.
allowed_commands = [
    "^playwright-cli",
]
```

### `[unboxexec]` Section

| Key | Type | Description |
|-----|------|-------------|
| `allowed_commands` | Array of strings | Regex patterns that define which commands are allowed to execute via `unboxexec`. The command and arguments are joined with spaces and matched against each pattern. If any pattern matches, the command is permitted. |

## Environment Variables

| Variable | Description |
|---|---|
| `CLAUDE_SANDBOX_WORKDIR` | Override working directory for sandbox execution |
| `CLAUDE_SANDBOX_CLAUDE_BIN` | Override path to the `claude` binary |

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
