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
