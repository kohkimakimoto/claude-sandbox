package command

import (
	"os"
	"path/filepath"
	"testing"
)

// testSetupFakeHome sets HOME to a fresh temp directory and returns its path.
// The .claude subdirectory is created so os.UserHomeDir()-based lookups work cleanly.
// The original HOME is restored on cleanup.
func testSetupFakeHome(t *testing.T) string {
	t.Helper()
	fakeHome := t.TempDir()
	if err := os.MkdirAll(filepath.Join(fakeHome, ".claude"), 0755); err != nil {
		t.Fatalf("failed to create fake home .claude dir: %v", err)
	}
	orig := os.Getenv("HOME")
	if err := os.Setenv("HOME", fakeHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Setenv("HOME", orig); err != nil {
			t.Errorf("failed to restore HOME: %v", err)
		}
	})
	return fakeHome
}

// testChdirTemp changes the working directory to a new temporary directory
// for the duration of the test, restoring it on cleanup.
func testChdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})
	return dir
}
