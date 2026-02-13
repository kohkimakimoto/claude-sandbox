# claude-sandbox

[![test](https://github.com/kohkimakimoto/claude-sandbox/actions/workflows/test.yml/badge.svg)](https://github.com/kohkimakimoto/claude-sandbox/actions/workflows/test.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kohkimakimoto/claude-sandbox/blob/main/LICENSE)

A wrapper around Claude Code (`claude` command) to run it in a sandboxed environment using macOS's `sandbox-exec`.

Table of Contents:
- [Why Not the Built-in Sandbox?](#why-not-the-built-in-sandbox)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration File](#configuration-file)
  - [Creating a Configuration File](#creating-a-configuration-file)
  - [Example](#example)
  - [`[sandbox]` Section](#sandbox-section)
  - [`[unboxexec]` Section](#unboxexec-section)
  - [Sandbox Profile Parameters](#sandbox-profile-parameters)
  - [Viewing the Sandbox Profile](#viewing-the-sandbox-profile)
- [Sandbox-External Command Execution](#sandbox-external-command-execution)
  - [The `unboxexec` Subcommand](#the-unboxexec-subcommand)
    - [Options](#options)
    - [Examples](#examples)
  - [Command Restrictions](#command-restrictions)
- [Environment Variables](#environment-variables)
- [License](#license)

## Why Not the Built-in Sandbox?

Claude Code provides a [built-in sandboxing feature](https://code.claude.com/docs/en/sandboxing) with filesystem and network isolation. I tried it, but in my workflow and environment it wasn't the best fit:

- Unexpected restrictions kept blocking legitimate operations, and I spent a lot of time troubleshooting and working around them.
- I didn't need network isolation at all, so it only added complexity without benefit.

What I actually needed was simpler: **restrict file writes to the current directory** and **explicitly allow exceptions** when needed. So I built this tool — minimal, predictable sandboxing with straightforward configuration.

> [!NOTE]
> v2 is a full rewrite from a shell script to Go. This enables single binary distribution and supports command execution outside the sandbox.
> **v2 is currently under active development and may introduce breaking changes without notice.**

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

## Configuration File

All settings are managed through a single TOML configuration file: `.claude/sandbox.toml`. The configuration file is resolved in the following order:

1. `.claude/sandbox.toml` in the current working directory (project-specific)
2. `~/.claude/sandbox.toml` (global)

The project-specific configuration takes precedence over the global configuration. If neither file exists, built-in defaults are used.

### Creating a Configuration File

Create a project-specific configuration:

```bash
claude-sandbox init
```

This creates `.claude/sandbox.toml` in your current directory.

Create a global configuration:

```bash
claude-sandbox init-global
```

This creates `~/.claude/sandbox.toml`.

### Example

```toml
# .claude/sandbox.toml or ~/.claude/sandbox.toml

[sandbox]
# Sandbox profile for sandbox-exec.
# If not set, the built-in default profile is used.
profile = '''
(version 1)
(allow default)
(deny file-write*)
(allow file-write*
    (subpath (param "WORKDIR"))
    (regex (string-append "^" (param "HOME") "/\\.claude"))
    (subpath "/tmp")
)
'''

# Override working directory (optional).
# workdir = "/path/to/workdir"

# Override claude binary path (optional).
# claude_bin = "/path/to/claude"

[unboxexec]
# Regex patterns for allowed commands.
# The command and its arguments are joined by spaces, and the resulting string
# is matched against each pattern. If any pattern matches, the command is allowed.
# If empty or not configured, all commands are rejected.
allowed_commands = [
    "^playwright-cli",
]
```

### `[sandbox]` Section

| Key | Type | Description |
|-----|------|-------------|
| `profile` | String | The sandbox-exec profile content. If not set, a built-in default profile is used. Use TOML multiline literal strings (`'''`) for readability. |
| `workdir` | String | Override the working directory for sandbox execution. If not set, the current directory is used. |
| `claude_bin` | String | Override the path to the `claude` binary. If not set, it is resolved from PATH. |

### `[unboxexec]` Section

| Key | Type | Description |
|-----|------|-------------|
| `allowed_commands` | Array of strings | Regex patterns that define which commands are allowed to execute via `unboxexec`. The command and arguments are joined with spaces and matched against each pattern. If any pattern matches, the command is permitted. |

### Sandbox Profile Parameters

The sandbox profile uses parameters that are passed from claude-sandbox automatically:

- `WORKDIR`: The current working directory where claude-sandbox is executed
- `HOME`: The user's home directory

You can use these parameters in your sandbox profile like this:

```scheme
(allow file-write*
    (subpath (param "WORKDIR"))
    (subpath (string-append (param "HOME") "/.claude"))
)
```

### Viewing the Sandbox Profile

You can view the actual profile being used:

```bash
claude-sandbox profile
```

The sandbox uses macOS's `sandbox-exec` (Apple Seatbelt) technology. Even if Claude Code tried to execute a command like `rm -rf /usr/bin` or modify system configuration files, the sandbox would block these operations.

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

By default, all commands executed via `unboxexec` are **rejected** unless explicitly allowed by the `[unboxexec]` section in the configuration file.

## Environment Variables

The following environment variables are set by claude-sandbox and available to the Claude Code process running inside the sandbox.

| Variable | Description |
|---|---|
| `CLAUDE_SANDBOX` | Set to `1` inside the sandbox |
| `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` | Path to the unboxexec daemon socket |
| `CLAUDE_SANDBOX_WORKDIR` | Working directory used for sandbox execution |
| `CLAUDE_SANDBOX_CLAUDE_BIN` | Path to the claude binary used |

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
