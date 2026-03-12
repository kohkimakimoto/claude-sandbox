package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/kohkimakimoto/claude-sandbox/v2/internal/version"
)

func TestVersionCommand(t *testing.T) {
	t.Run("via version subcommand", func(t *testing.T) {
		buf := &bytes.Buffer{}
		app := newApp()
		app.Writer = buf

		if err := app.Run(context.Background(), []string{"claude-sandbox", "version"}); err != nil {
			t.Fatalf("version command failed: %v", err)
		}

		expected := "claude-sandbox version " + version.Version + " (commit: " + version.CommitHash + ")\n"
		if got := buf.String(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("via -v flag", func(t *testing.T) {
		buf := &bytes.Buffer{}
		app := newApp()
		app.Writer = buf

		if err := app.Run(context.Background(), []string{"claude-sandbox", "-v"}); err != nil {
			t.Fatalf("-v flag failed: %v", err)
		}

		expected := "claude-sandbox version " + version.Version + " (commit: " + version.CommitHash + ")\n"
		if got := buf.String(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("via --version flag", func(t *testing.T) {
		buf := &bytes.Buffer{}
		app := newApp()
		app.Writer = buf

		if err := app.Run(context.Background(), []string{"claude-sandbox", "--version"}); err != nil {
			t.Fatalf("--version flag failed: %v", err)
		}

		expected := "claude-sandbox version " + version.Version + " (commit: " + version.CommitHash + ")\n"
		if got := buf.String(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}
