package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultProfile is the built-in sandbox profile used when no custom profile is found.
const DefaultProfile = `;; This is a default built-in sandbox profile for claude-sandbox.
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code
    (regex (string-append "^" (param "HOME") "/\\.claude"))

    ;; Keychain access for Claude Code credentials
    (subpath (string-append (param "HOME") "/Library/Keychains"))

    ;; Temporary directories and files
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")

    ;; Home directory
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))
    (subpath (string-append (param "HOME") "/Library/Caches"))
    (regex (string-append "^" (param "HOME") "/\\.viminfo"))

    ;; devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
    (literal "/dev/dtracehelper")
    (regex #"^/dev/tty*")
)
`

// ProjectProfileTemplate is the template for project-specific sandbox profiles.
const ProjectProfileTemplate = `;; This is a project-specific sandbox profile for claude-sandbox.
;; You can customize this file to suit your project's needs.
;; see https://github.com/kohkimakimoto/claude-sandbox
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code
    (regex (string-append "^" (param "HOME") "/\\.claude"))

    ;; Keychain access for Claude Code credentials
    (subpath (string-append (param "HOME") "/Library/Keychains"))

    ;; Temporary directories and files
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")

    ;; Home directory
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))
    (subpath (string-append (param "HOME") "/Library/Caches"))
    (regex (string-append "^" (param "HOME") "/\\.viminfo"))

    ;; devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
    (literal "/dev/dtracehelper")
    (regex #"^/dev/tty*")
)
`

// GlobalProfileTemplate is the template for global sandbox profiles.
const GlobalProfileTemplate = `;; This is a global sandbox profile for claude-sandbox.
;; You can customize this file to suit your needs.
;; see https://github.com/kohkimakimoto/claude-sandbox
(version 1)

(allow default)

(deny file-write*)
(allow file-write*
    ;; Working directory
    (subpath (param "WORKDIR"))

    ;; Claude Code
    (regex (string-append "^" (param "HOME") "/\\.claude"))

    ;; Keychain access for Claude Code credentials
    (subpath (string-append (param "HOME") "/Library/Keychains"))

    ;; Temporary directories and files
    (subpath "/tmp")
    (subpath "/var/folders")
    (subpath "/private/tmp")
    (subpath "/private/var/folders")

    ;; Home directory
    (subpath (string-append (param "HOME") "/.npm"))
    (subpath (string-append (param "HOME") "/.cache"))
    (subpath (string-append (param "HOME") "/Library/Caches"))
    (regex (string-append "^" (param "HOME") "/\\.viminfo"))

    ;; devices
    (literal "/dev/stdout")
    (literal "/dev/stderr")
    (literal "/dev/null")
    (literal "/dev/dtracehelper")
    (regex #"^/dev/tty*")
)
`

// BuildProfile creates a temporary file with the sandbox profile and returns
// its path and a cleanup function. The profile is resolved in this order:
// 1. .claude/sandbox.sb in the working directory (project-specific)
// 2. $HOME/.claude/sandbox.sb (global)
// 3. Built-in default profile
func BuildProfile() (profilePath string, cleanup func(), err error) {
	workdir := GetWorkdir()
	home, _ := os.UserHomeDir()

	// Determine profile content
	var content []byte

	projectProfile := filepath.Join(workdir, ".claude", "sandbox.sb")
	globalProfile := filepath.Join(home, ".claude", "sandbox.sb")

	if _, err := os.Stat(projectProfile); err == nil {
		content, err = os.ReadFile(projectProfile)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read profile %s: %w", projectProfile, err)
		}
	} else if _, err := os.Stat(globalProfile); err == nil {
		content, err = os.ReadFile(globalProfile)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read profile %s: %w", globalProfile, err)
		}
	} else {
		content = []byte(DefaultProfile)
	}

	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "claude-sandbox-profile-*.sb")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("failed to write profile: %w", err)
	}
	tmpFile.Close()

	cleanup = func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup, nil
}
