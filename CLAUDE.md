# claude-sandbox

## Overview

A wrapper tool to safely run Claude Code in a sandboxed environment on macOS.
Implemented in Go for single binary distribution and integrated process management.

## Features

### Sandboxed Execution

- Execute Claude Code using macOS `sandbox-exec`
- Sandbox profile resolution: project-specific (`.claude/sandbox.sb`) → global (`$HOME/.claude/sandbox.sb`) → built-in default
- Transparent argument passing to Claude Code

### Sandbox-External Command Execution (unboxexec)

- Built-in daemon to execute commands outside the sandbox
- Communication via Unix Domain Socket (`/tmp/claude-sandbox-unboxexec-{PID}.sock`)
- Socket path is passed to Claude Code via `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` environment variable

## Architecture

```
claude-sandbox (single process)
│
├─ unboxexec daemon (goroutine)
│  ├─ Listen on Unix socket: /tmp/claude-sandbox-unboxexec-{PID}.sock
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

### Communication Protocol

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

## Directory Structure

```
claude-sandbox/
├── cmd/
│   └── claude-sandbox/
│       └── main.go                # Entry point
├── internal/
│   ├── command/
│   │   ├── app.go                 # Run(), newApp() — CLI application setup
│   │   ├── claude.go              # ClaudeCommand, RunClaudeAction
│   │   ├── init.go                # InitCommand — create project sandbox profile
│   │   ├── init_global.go         # InitGlobalCommand — create global sandbox profile
│   │   ├── profile.go             # ProfileCommand — print evaluated profile
│   │   └── help.go                # RootHelpTemplate, HelpTemplate
│   ├── sandbox/
│   │   ├── env.go                 # GetWorkdir, GetClaudeBin, SocketPath
│   │   └── profile.go             # BuildProfile, profile templates
│   ├── unboxexec/
│   │   └── daemon.go              # StartDaemon, ExecRequest/Response, command execution
│   └── version/
│       └── version.go             # Version, CommitHash (set via ldflags)
├── Makefile
├── go.mod
└── go.sum
```

### Package Dependencies

```
cmd/claude-sandbox  → command
                      command  → sandbox, unboxexec, version
                      sandbox  (no internal dependencies)
                      unboxexec (no internal dependencies)
                      version  (no internal dependencies)
```

No circular dependencies.

## Environment Variables

| Variable | Description |
|---|---|
| `CLAUDE_SANDBOX` | Set to `1` when running inside claude-sandbox (set automatically) |
| `CLAUDE_SANDBOX_WORKDIR` | Override working directory for sandbox execution |
| `CLAUDE_SANDBOX_CLAUDE_BIN` | Override path to claude binary |
| `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` | Unix socket path for unboxexec daemon (set automatically) |

## Development

- macOS only (`sandbox-exec` is macOS-specific)
- Requires Go 1.24+
- External dependency: `github.com/urfave/cli/v3`
- `make build` — Build dev binary to `.dev/build/dev/claude-sandbox`
- `make test` — Run tests
- `make format` — Format source code
- Version and commit hash are injected via `-ldflags` at build time
