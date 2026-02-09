package main

import (
	"fmt"
	"os"

	"github.com/kohkimakimoto/claude-sandbox/v2/internal/command"
)

func main() {
	if err := command.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
