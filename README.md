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

`claude-sandbox` can be used as a drop-in replacement for the `claude` command.
However, it runs in a sandboxed environment, which restricts the command's access to the file system, network, and other resources.

## License

The MIT License (MIT)

Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
