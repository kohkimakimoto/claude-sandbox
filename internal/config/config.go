package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
)

// Config represents the claude-sandbox configuration.
type Config struct {
	Sandbox   SandboxConfig   `toml:"sandbox"`
	Unboxexec UnboxexecConfig `toml:"unboxexec"`
}

// SandboxConfig holds settings for the sandbox environment.
type SandboxConfig struct {
	// Profile is the sandbox-exec profile content.
	// If empty, the built-in default profile is used.
	Profile string `toml:"profile"`
	// Workdir overrides the working directory for sandbox execution.
	// If empty, the current directory is used.
	Workdir string `toml:"workdir"`
	// ClaudeBin overrides the path to the claude binary.
	// If empty, PATH search is used.
	ClaudeBin string `toml:"claude_bin"`
}

// UnboxexecConfig holds settings for the unboxexec daemon.
type UnboxexecConfig struct {
	AllowedCommands []string `toml:"allowed_commands"`
}

// ResolveConfigPath resolves the config file path in the following order:
// 1. .claude/sandbox.toml in the working directory (project-specific)
// 2. ~/.claude/sandbox.toml (global)
// Returns empty string if neither exists.
func ResolveConfigPath() string {
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()

	projectConfig := filepath.Join(wd, ".claude", "sandbox.toml")
	if _, err := os.Stat(projectConfig); err == nil {
		return projectConfig
	}

	globalConfig := filepath.Join(home, ".claude", "sandbox.toml")
	if _, err := os.Stat(globalConfig); err == nil {
		return globalConfig
	}

	return ""
}

// Load reads and parses a TOML config file at the given path.
// If the path is empty or the file does not exist, it returns an empty Config without error.
func Load(path string) (*Config, error) {
	cfg := &Config{}

	if path == "" {
		return cfg, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config %s: %w", path, err)
	}

	return cfg, nil
}

// CompileAllowedCommands compiles a list of regex pattern strings into []*regexp.Regexp.
func CompileAllowedCommands(patterns []string) ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_commands pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}
