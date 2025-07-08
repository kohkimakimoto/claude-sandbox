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
The default sandbox profile restricts write access to only safe locations:
your current working directory, Claude Code configuration files, temporary directories, and common cache directories like `~/.npm` and `~/.cache`.

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

### Viewing the Current Profile

To see the sandbox profile that will be used:

```bash
claude-sandbox profile
```

This shows the complete profile with all parameters resolved, helping you understand what permissions are currently configured.

### Understanding Sandbox Profiles

Sandbox profiles are written in Apple's Seatbelt language and define what resources Claude Code can access. The profile uses parameters that are passed from the environment:

- `WORKDIR`: The current working directory where claude-sandbox is executed
- `HOME`: The user's home directory

#### Basic Profile Structure

```scheme
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code configuration
    (subpath (string-append (param "HOME") "/.claude"))
    (literal (string-append (param "HOME") "/.claude.json"))
    
    ;; Temporary directories
    (subpath "/tmp")
    (subpath "/var/folders")
    
    ;; Output devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
)
```

### Customizing Sandbox Profiles

After creating a profile with `claude-sandbox init` or `claude-sandbox init-global`, you can edit the `.claude/sandbox.sb` file to customize permissions.

#### Common Customizations

**Allow access to additional directories:**

```scheme
(allow file-write*
    ;; Allow access to a specific project directory
    (subpath "/path/to/your/project")
    
    ;; Allow access to Documents folder
    (subpath (string-append (param "HOME") "/Documents"))
)
```

**Allow network access:**

```scheme
;; Allow outbound network connections
(allow network-outbound*)

;; Allow specific network destinations
(allow network-outbound*
    (remote tcp "api.anthropic.com:443")
    (remote tcp "github.com:443")
)
```

**Allow access to system tools:**

```scheme
(allow file-read* file-write*
    ;; Allow access to development tools
    (subpath "/usr/local/bin")
    (subpath "/opt/homebrew/bin")
)
```

#### Using Parameters in Profiles

Always use `(param "NAME")` syntax instead of environment variables:

```scheme
;; ✅ Correct: Use parameters
(subpath (param "WORKDIR"))
(subpath (string-append (param "HOME") "/.config"))

;; ❌ Incorrect: Don't use environment variables directly
(subpath "$CLAUDE_SANDBOX_WORKDIR")
(subpath "$HOME/.config")
```

### Environment Variables

You can override the working directory by setting the `CLAUDE_SANDBOX_WORKDIR` environment variable:

```bash
# Run claude-sandbox in a specific directory
CLAUDE_SANDBOX_WORKDIR=/path/to/project claude-sandbox

# Or export it for the session
export CLAUDE_SANDBOX_WORKDIR=/path/to/project
claude-sandbox
```

### Advanced Usage Examples

#### Development Environment Setup

For a typical development environment, you might want to allow access to common development directories:

```scheme
(version 1)
(allow default)

(deny file-write*)
(allow file-write*
    ;; Project workspace
    (subpath (param "WORKDIR"))
    
    ;; Development tools and caches
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))
    (subpath (string-append (param "HOME") "/.local"))
    
    ;; Claude Code
    (subpath (string-append (param "HOME") "/.claude"))
    (literal (string-append (param "HOME") "/.claude.json"))
    (literal (string-append (param "HOME") "/.claude.json.lock"))
    (literal (string-append (param "HOME") "/.claude.json.backup"))
    
    ;; System directories
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")
    
    ;; Output devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
)

;; Allow network access for package managers and APIs
(allow network-outbound*)
```

#### Restricted Environment

For a more restricted environment that only allows access to the current project:

```scheme
(version 1)
(allow default)

(deny file-write*)
(allow file-write*
    ;; Only current project
    (subpath (param "WORKDIR"))
    
    ;; Minimal Claude Code access
    (subpath (string-append (param "HOME") "/.claude"))
    (literal (string-append (param "HOME") "/.claude.json"))
    
    ;; Essential system access
    (subpath "/tmp")
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
)

;; No network access
(deny network*)
```

### Troubleshooting

#### Check Profile Syntax

Use the `profile` command to validate your profile syntax:

```bash
claude-sandbox profile
```

If there are syntax errors, they will be displayed when the profile is evaluated.

#### Common Issues

1. **Permission Denied Errors**: Add the required paths to your sandbox profile
2. **Network Access Issues**: Add network permissions if Claude Code needs internet access
3. **File Access Issues**: Ensure the working directory and necessary paths are included in `file-write*` rules

#### Debug Mode

To see the exact sandbox-exec command being executed, you can uncomment the debug line in the script or run:

```bash
# This will show the sandbox-exec command that would be run
claude-sandbox profile > /tmp/debug-profile.sb
echo "sandbox-exec -D WORKDIR=$PWD -D HOME=$HOME -f /tmp/debug-profile.sb claude"
```

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
