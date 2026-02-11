package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
[unboxexec]
allowed_commands = [
    "^playwright",
    "^echo hello",
]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Unboxexec.AllowedCommands) != 2 {
		t.Fatalf("expected 2 allowed_commands, got %d", len(cfg.Unboxexec.AllowedCommands))
	}
	if cfg.Unboxexec.AllowedCommands[0] != "^playwright" {
		t.Errorf("expected %q, got %q", "^playwright", cfg.Unboxexec.AllowedCommands[0])
	}
	if cfg.Unboxexec.AllowedCommands[1] != "^echo hello" {
		t.Errorf("expected %q, got %q", "^echo hello", cfg.Unboxexec.AllowedCommands[1])
	}
}

func TestLoadConfigWithSandboxSection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
[sandbox]
profile = "(version 1)\n(allow default)"
workdir = "/tmp/myworkdir"
claude_bin = "/usr/local/bin/claude"

[unboxexec]
allowed_commands = [
    "^playwright-cli",
]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Sandbox.Profile != "(version 1)\n(allow default)" {
		t.Errorf("expected profile %q, got %q", "(version 1)\n(allow default)", cfg.Sandbox.Profile)
	}
	if cfg.Sandbox.Workdir != "/tmp/myworkdir" {
		t.Errorf("expected workdir %q, got %q", "/tmp/myworkdir", cfg.Sandbox.Workdir)
	}
	if cfg.Sandbox.ClaudeBin != "/usr/local/bin/claude" {
		t.Errorf("expected claude_bin %q, got %q", "/usr/local/bin/claude", cfg.Sandbox.ClaudeBin)
	}
	if len(cfg.Unboxexec.AllowedCommands) != 1 {
		t.Fatalf("expected 1 allowed_commands, got %d", len(cfg.Unboxexec.AllowedCommands))
	}
}

func TestLoadConfigWithMultilineProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
[sandbox]
profile = '''
(version 1)
(allow default)
(deny file-write*)
'''
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// TOML multiline literal strings (''') strip the first newline
	expected := "(version 1)\n(allow default)\n(deny file-write*)\n"
	if cfg.Sandbox.Profile != expected {
		t.Errorf("expected profile %q, got %q", expected, cfg.Sandbox.Profile)
	}
}

func TestLoadConfigSandboxDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
[unboxexec]
allowed_commands = ["^echo"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Sandbox fields should be zero values when not specified
	if cfg.Sandbox.Profile != "" {
		t.Errorf("expected empty profile, got %q", cfg.Sandbox.Profile)
	}
	if cfg.Sandbox.Workdir != "" {
		t.Errorf("expected empty workdir, got %q", cfg.Sandbox.Workdir)
	}
	if cfg.Sandbox.ClaudeBin != "" {
		t.Errorf("expected empty claude_bin, got %q", cfg.Sandbox.ClaudeBin)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected no error for empty path, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Unboxexec.AllowedCommands) != 0 {
		t.Errorf("expected empty allowed_commands, got %d", len(cfg.Unboxexec.AllowedCommands))
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Unboxexec.AllowedCommands) != 0 {
		t.Errorf("expected empty allowed_commands, got %d", len(cfg.Unboxexec.AllowedCommands))
	}
}

func TestLoadInvalidTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("invalid [[[ toml"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}

func TestCompileAllowedCommands(t *testing.T) {
	patterns := []string{"^playwright", "^echo .*"}
	compiled, err := CompileAllowedCommands(patterns)
	if err != nil {
		t.Fatalf("CompileAllowedCommands failed: %v", err)
	}
	if len(compiled) != 2 {
		t.Fatalf("expected 2 compiled patterns, got %d", len(compiled))
	}

	// Test matching
	if !compiled[0].MatchString("playwright install chromium") {
		t.Error("expected pattern to match 'playwright install chromium'")
	}
	if compiled[0].MatchString("echo playwright") {
		t.Error("expected pattern not to match 'echo playwright'")
	}
	if !compiled[1].MatchString("echo hello world") {
		t.Error("expected pattern to match 'echo hello world'")
	}
}

func TestCompileAllowedCommandsEmpty(t *testing.T) {
	compiled, err := CompileAllowedCommands(nil)
	if err != nil {
		t.Fatalf("CompileAllowedCommands failed: %v", err)
	}
	if len(compiled) != 0 {
		t.Errorf("expected 0 compiled patterns, got %d", len(compiled))
	}
}

func TestCompileAllowedCommandsInvalidPattern(t *testing.T) {
	patterns := []string{"[invalid"}
	_, err := CompileAllowedCommands(patterns)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}
