# claude-sandbox

A wrapper around the `claude` command to run it in a sandboxed environment using macOS's `sandbox-exec`.

It also includes a built-in daemon ([unboxexec](#sandbox-external-command-execution)) that allows Claude Code running inside the sandbox to execute specific commands outside the sandbox via a Unix Domain Socket.

## Installation

Download the binary from the [GitHub Releases](https://github.com/kohkimakimoto/claude-sandbox/releases) page and place it in your `PATH`.

Or build from source:

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

`claude-sandbox` includes a built-in daemon called **unboxexec** that runs outside the sandbox as a goroutine. When `claude-sandbox` starts, the daemon listens on a Unix Domain Socket and the socket path is passed to Claude Code via the `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` environment variable.

Claude Code running inside the sandbox can send JSON requests to this socket to execute commands outside the sandbox.

### Protocol

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

## Environment Variables

| Variable | Description |
|---|---|
| `CLAUDE_SANDBOX_WORKDIR` | Override working directory for sandbox execution |
| `CLAUDE_SANDBOX_CLAUDE_BIN` | Override path to the `claude` binary |
| `CLAUDE_SANDBOX_UNBOXEXEC_SOCK` | Unix socket path for the unboxexec daemon (set automatically) |

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
