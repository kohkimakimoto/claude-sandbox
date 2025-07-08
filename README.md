# claude-sandbox

A wrapper around the `claude` command to run it in a sandboxed environment using macOS's sandbox-exec.

## Installation

`claude-sandbox` is a simple, single-file Bash script. You can download the [`claude-sandbox`](https://github.com/kohkimakimoto/claude-sandbox/raw/main/claude-sandbox) file from the repository and make the file executable.

The following command will download the script and install it to `/usr/local/bin/claude-sandbox`:

```bash
curl -sSL https://raw.githubusercontent.com/kohkimakimoto/claude-sandbox/refs/heads/main/claude-sandbox | sudo tee /usr/local/bin/claude-sandbox > /dev/null && sudo chmod +x /usr/local/bin/claude-sandbox
```

To check if the installation was successful, you can run the following command to see the help message:

```bash
claude-sandbox -h
```

## Usage

`claude-sandbox` can be used as a drop-in replacement for the `claude` command, but runs in a sandboxed environment that restricts access to the file system, network, and other resources.

The simplest way to use `claude-sandbox` is as a direct replacement for the `claude` command:

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

By default, `claude-sandbox` uses a built-in sandbox profile that restricts file system access to the current working directory, Claude Code configuration, and temporary directories. 
You can view the actual profile being used with `claude-sandbox profile` command.

```bash
claude-sandbox profile
```

You will get output as the following:

```scheme
;; This is a default built-in sandbox profile for claude-sandbox.
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code
    (subpath (string-append (param "HOME") "/.claude"))
    (literal (string-append (param "HOME") "/.claude.json"))
    (literal (string-append (param "HOME") "/.claude.json.lock"))
    (literal (string-append (param "HOME") "/.claude.json.backup"))

    ;; Temporary directories and files
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")

    ;; Home directory
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))

    ;; devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
)
```

The sandbox uses macOS's `sandbox-exec` (Apple Seatbelt) technology. 
This sandbox profile protects you from accidentally breaking your system by preventing Claude Code from modifying files outside the allowed areas. For example, even if Claude Code tried to execute a command like `rm -rf /usr/bin` or modify system configuration files, the sandbox would block these operations:

### Configuring Sandbox Profiles

To customize the sandbox environment, you need to create a sandbox profile. There are two types of profiles:

#### Project-Specific Profile

Create a project-specific sandbox profile that applies only to the current project:

```bash
claude-sandbox init
```

This creates `.claude/sandbox.sb` in your current directory.
You can then edit this file to customize the sandbox permissions for your project.

#### Global Profile

Create a global sandbox profile that applies to all projects:

```bash
claude-sandbox init-global
```

This creates `~/.claude/sandbox.sb`.

**Profile Priority**: Project-specific profiles take precedence over global profiles. If neither exists, a built-in default profile is used.

#### Parameters

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

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
