package unboxexec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSendRequest(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := StartDaemon(ctx, sockPath); err != nil {
		t.Fatalf("failed to start daemon: %v", err)
	}

	resp, err := SendRequest(sockPath, &ExecRequest{
		Command: "echo",
		Args:    []string{"hello"},
	})
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}

	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Stdout != "hello\n" {
		t.Errorf("expected stdout %q, got %q", "hello\n", resp.Stdout)
	}
	if resp.Error != "" {
		t.Errorf("unexpected error: %s", resp.Error)
	}
}

func TestSendRequestWithDir(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := StartDaemon(ctx, sockPath); err != nil {
		t.Fatalf("failed to start daemon: %v", err)
	}

	resp, err := SendRequest(sockPath, &ExecRequest{
		Command: "pwd",
		Dir:     os.TempDir(),
	})
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}

	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Error != "" {
		t.Errorf("unexpected error: %s", resp.Error)
	}
}

func TestSendRequestWithEnv(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := StartDaemon(ctx, sockPath); err != nil {
		t.Fatalf("failed to start daemon: %v", err)
	}

	resp, err := SendRequest(sockPath, &ExecRequest{
		Command: "sh",
		Args:    []string{"-c", "echo $TEST_VAR"},
		Env:     map[string]string{"TEST_VAR": "myvalue"},
	})
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}

	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Stdout != "myvalue\n" {
		t.Errorf("expected stdout %q, got %q", "myvalue\n", resp.Stdout)
	}
}

func TestSendRequestNoCommand(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := StartDaemon(ctx, sockPath); err != nil {
		t.Fatalf("failed to start daemon: %v", err)
	}

	resp, err := SendRequest(sockPath, &ExecRequest{})
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected error for empty command")
	}
}

func TestSendRequestConnectionError(t *testing.T) {
	_, err := SendRequest("/tmp/nonexistent.sock", &ExecRequest{
		Command: "echo",
		Args:    []string{"hello"},
	})
	if err == nil {
		t.Error("expected error for nonexistent socket")
	}
}
