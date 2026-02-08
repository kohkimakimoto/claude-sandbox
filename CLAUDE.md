# claude-sandbox (Go Reimplementation)

## Overview

A wrapper tool to safely run Claude Code in a sandboxed environment on macOS.
Reimplemented in Go from shell script version to achieve single binary distribution and feature integration.

## Background

### Existing Shell Script Version

I created a shell script [`claude-sandbox`](claude-sandbox) to run Claude Code in a macOS sandbox using `sandbox-exec`.
But it has limitations:

- Uses macOS `sandbox-exec` to run Claude Code in a restricted environment
- Controls filesystem access via configuration file (`sandbox.sb`)
- Implemented as a simple shell script

### Problems

**Browser-based tools limitation**:
MCP/Skills that execute browsers (like Playwright) don't work inside sandbox.
Playwright itself uses sandbox, so nested sandboxes cannot run.

### Purpose of Go Reimplementation

- Integrate sandbox-external command execution capability
- Easy distribution and management with single binary
- Centralized process lifecycle management

## Features

### 1. Sandboxed Execution (Existing Feature Migration)

- Execute Claude Code using `sandbox-exec`
- Apply sandbox configuration
- Transparent argument passing to Claude Code

### 2. Sandbox-External Command Execution (New Feature)

- Built-in daemon to execute commands outside sandbox
- Communication via Unix Domain Socket
- Accessible from Claude Code running inside sandbox

## Architecture
```
claude-sandbox
│
├─ Main Process
│  ├─ Start internal daemon
│  ├─ Set environment variables (CLAUDE_SANDBOX_PID)
│  ├─ Execute claude code via sandbox-exec
│  └─ Cleanup on exit
│
└─ Internal Daemon (child process)
   ├─ Create Unix socket: /tmp/claude-sandbox-{PID}.sock
   ├─ Accept command execution requests
   └─ Execute outside sandbox
```

## Technical Specifications

### Command Structure
```bash
# Basic usage (same as existing)
claude-sandbox [claude code options...]

# Internal behavior
# 1. Start internal daemon in background
# 2. Set CLAUDE_SANDBOX_PID={daemon_pid} environment variable
# 3. sandbox-exec -f sandbox.sb claude code [options...]
# 4. Shutdown daemon on exit
```

### Communication Protocol

**Unix Socket**: `/tmp/claude-sandbox-{PID}.sock`

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

### Configuration Files

Use existing `sandbox.sb` as-is (loaded and applied from Go)

Optional: `~/.claude-sandbox/config.json` for allowed commands
```json
{
  "allowed_commands": ["playwright", "chromium", "firefox"]
}
```

## Implementation Plan

### Phase 1: Minimum Implementation

1. ✅ Basic sandbox execution (migrate existing features)
2. ✅ Internal daemon startup/shutdown management
3. ✅ Unix socket communication

### Phase 2: Extensions (as needed)

- Command restriction via config file
- Logging
- Debug mode

## Directory Structure
```
claude-sandbox/
├── cmd/
│   └── claude-sandbox/
│       └── main.go          # Entry point
├── internal/
│   ├── sandbox/
│   │   └── executor.go      # Sandbox execution
│   └── daemon/
│       └── server.go         # Unix socket server & command executor
├── sandbox.sb                # Sandbox config (existing)
├── go.mod
└── README.md
```

## Notes

- macOS only (`sandbox-exec` is macOS-specific)
- Requires Go 1.21+
- Keep it simple, prefer standard library
- Explicit error handling