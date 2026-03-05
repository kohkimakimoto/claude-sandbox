# claude-sandbox

## Overview

A wrapper tool to safely run Claude Code in a sandboxed environment on macOS.
Implemented in Go.
Runs a background daemon (goroutine) to support sandbox-bypass execution.

## Features

### Sandboxed Execution

- Execute Claude Code using macOS `sandbox-exec`
- Sandbox profile is configured via `[sandbox].profile` in `sandbox.toml`; if not set, built-in default is used
- Transparent argument passing to Claude Code

### Sandbox-External Command Execution (unboxexec)

- Built-in daemon to execute commands outside the sandbox
- Communication via Unix Domain Socket (`{TMPDIR}/claude-sandbox-unboxexec-{PID}.sock`)
- Socket path is passed to Claude Code via `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` environment variable

### Configuration File

- Three-tier TOML configuration with layered merging
- Config resolution order (each level overrides the previous for any field that is explicitly set):
  1. User: `~/.claude/sandbox.toml`
  2. Project: `.claude/sandbox.toml` in working directory
  3. Local: `.claude/sandbox.local.toml` in working directory (gitignore-friendly overrides)
- If no config files exist, built-in defaults are used

```toml
# ~/.claude/sandbox.toml         (user-level)
# .claude/sandbox.toml           (project-level)
# .claude/sandbox.local.toml     (local overrides)

[sandbox]
# Sandbox profile for sandbox-exec (multiline literal string).
# If not set, the built-in default profile is used.
# profile = '''
# (version 1)
# (allow default)
# ...
# '''

# Override working directory (optional).
# workdir = "/path/to/workdir"

# Override claude binary path (optional).
# claude_bin = "/path/to/claude"

[unboxexec]
# Regex patterns for allowed commands.
# The command + args joined by spaces is matched against each pattern.
# If any pattern matches, the command is allowed.
# If empty or not configured, all commands are rejected.
allowed_commands = [
    "^playwright-cli",
]
```

## Architecture

```
claude-sandbox (single process)
│
├─ unboxexec daemon (goroutine)
│  ├─ Listen on Unix socket: {TMPDIR}/claude-sandbox-unboxexec-{PID}.sock
│  ├─ Accept JSON command execution requests
│  └─ Execute commands outside sandbox, return JSON responses
│
└─ sandbox-exec → claude (child process)
   ├─ Runs inside macOS sandbox
   ├─ Receives CLAUDE_SANDBOX_UNBOXEXEC_SOCK env var
   └─ Can request sandbox-external execution via Unix socket
```

The `claude-sandbox` process starts the unboxexec daemon as a goroutine, then
spawns `sandbox-exec` as a child process. When claude exits, the context is
cancelled and the daemon goroutine shuts down.

### Unboxexec Communication Protocol

The unboxexec daemon communicates with its clients via JSON over the Unix Domain Socket.
This protocol is only used for sandbox-bypass command execution (unboxexec); the main sandboxed execution of Claude Code does not involve this protocol.

**Request**:
```json
{
  "command": "playwright",
  "args": ["install", "chromium"],
  "env": {"KEY": "value"},
  "dir": "/path/to/workdir",
  "timeout": 300
}
```

**Response**:
```json
{
  "stdout": "...",
  "stderr": "...",
  "exit_code": 0,
  "error": ""
}
```

## Package Structure

| Package | Description |
|---|---|
| `cmd/claude-sandbox` | Entry point (`main.go`) |
| `internal/command` | CLI application setup, subcommand definitions (claude, init, init-user, profile, unboxexec, etc.) |
| `internal/config` | TOML configuration loading and allowed-command compilation |
| `internal/sandbox` | Sandbox profile building, environment variable helpers |
| `internal/unboxexec` | Unboxexec daemon (server) and client |
| `internal/version` | Version and commit hash (set via `-ldflags` at build time) |

`cmd/claude-sandbox` depends on `internal/command`, which depends on all other `internal/*` packages. The `internal/*` packages have no circular dependencies among themselves.

## Environment Variables

| Variable | Description |
|---|---|
| `CLAUDE_SANDBOX` | Set to `1` to indicate the process is running inside the sandbox (set by claude-sandbox for child process) |
| `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` | Unix socket path for communicating with the unboxexec daemon (set by claude-sandbox for child process) |
| `CLAUDE_SANDBOX_WORKDIR` | Working directory used for sandbox execution (set by claude-sandbox for child process) |
| `CLAUDE_SANDBOX_CLAUDE_BIN` | Path to claude binary (set by claude-sandbox for child process) |

## Development

- macOS only (`sandbox-exec` is macOS-specific)
- Go version and external dependencies are defined in `go.mod`
- `make build` — Build dev binary to `.dev/build/dev/claude-sandbox`
- `make test` — Run tests
- `make format` — Format source code
- Version and commit hash are injected via `-ldflags` at build time
