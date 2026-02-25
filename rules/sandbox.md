# Sandbox Environment

For security purposes and to prevent unexpected destructive behavior, Claude Code is typically executed within a sandboxed environment. This uses Apple's `sandbox-exec` to restrict access to unnecessary system resources.
You can check whether Claude Code is running inside a sandbox environment by verifying if the `CLAUDE_SANDBOX` environment variable is set.

```bash
echo $CLAUDE_SANDBOX
# => 1
```

## Checking the Profile

The `sandbox-exec` profile that defines the sandbox environment can be printed to standard output with the following command. This is useful when troubleshooting issues caused by the sandbox environment.

```bash
claude-sandbox profile
```

## Executing Commands Outside the Sandbox

In some use cases, you may need to bypass the sandbox restrictions to execute commands.
For example, browser automation uses its own sandbox, which can cause sandbox nesting and may not work correctly.

The `claude-sandbox unboxexec` subcommand is provided as a workaround for such cases, allowing you to execute commands outside the sandbox.

**Important!:** `claude-sandbox unboxexec` bypasses sandbox protections and may therefore expose the system to risk. For this reason, you MUST ask for explicit user approval before using this command. You MUST NOT execute arbitrary commands at your own discretion.

### Basic Usage

```bash
claude-sandbox unboxexec [<options>] -- <command> [<args...>]
```

### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-C` | Specify the working directory for the command |
| `--timeout` | `-t` | Timeout in seconds (default: 60 seconds) |
| `--env` | `-e` | Specify environment variables in `KEY=VALUE` format (can be specified multiple times) |

### Examples

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

### Permitted Use Cases without User Approval

You can use `claude-sandbox unboxexec` without asking for user approval in the following specific scenarios:

#### playwright-cli (Browser Automation)

`playwright-cli` is a command-line tool for browser automation.
This command does not work inside the sandbox. It must be executed through `unboxexec`.

Usage example:

```bash
claude-sandbox unboxexec playwright-cli open https://example.com/
```
