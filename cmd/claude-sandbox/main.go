package main

import (
	"fmt"
	"os"

	"github.com/kohkimakimoto/claude-sandbox/internal"
)

func main() {
	if err := internal.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
