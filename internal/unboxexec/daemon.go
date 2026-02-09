package unboxexec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ExecRequest represents a command execution request from inside the sandbox.
type ExecRequest struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
	Dir     string            `json:"dir"`
	Timeout int               `json:"timeout"`
}

// ExecResponse represents the result of a command execution.
type ExecResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error"`
}

const defaultTimeout = 60 // seconds

// StartDaemon starts a Unix Domain Socket server that accepts command execution
// requests. It runs in a goroutine and stops when ctx is cancelled.
// The socket file is cleaned up on shutdown.
// allowedCommands specifies regex patterns that the command string must match.
// If allowedCommands is empty, all commands are rejected.
func StartDaemon(ctx context.Context, sockPath string, allowedCommands []*regexp.Regexp) error {
	// Remove stale socket file if it exists
	os.Remove(sockPath)

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", sockPath, err)
	}

	// Start accept loop in a goroutine
	go func() {
		defer listener.Close()
		defer os.Remove(sockPath)

		for {
			conn, err := listener.Accept()
			if err != nil {
				// Check if the context was cancelled (normal shutdown)
				if ctx.Err() != nil {
					return
				}
				// Check if the listener was closed
				if errors.Is(err, net.ErrClosed) {
					return
				}
				continue
			}
			go handleConnection(ctx, conn, allowedCommands)
		}
	}()

	// Close the listener when ctx is cancelled
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	return nil
}

// handleConnection processes a single request on the connection.
func handleConnection(ctx context.Context, conn net.Conn, allowedCommands []*regexp.Regexp) {
	defer conn.Close()

	var req ExecRequest
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		resp := ExecResponse{Error: fmt.Sprintf("failed to decode request: %v", err)}
		json.NewEncoder(conn).Encode(resp)
		return
	}

	resp := executeCommand(ctx, &req, allowedCommands)
	json.NewEncoder(conn).Encode(resp)
}

// validateCommand checks whether the command is allowed by the configured patterns.
// It joins the command and args with spaces, then checks against each pattern.
func validateCommand(req *ExecRequest, allowedCommands []*regexp.Regexp) error {
	if len(allowedCommands) == 0 {
		return fmt.Errorf("command not allowed: no allowed_commands configured")
	}

	cmdStr := req.Command
	if len(req.Args) > 0 {
		cmdStr = cmdStr + " " + strings.Join(req.Args, " ")
	}

	for _, re := range allowedCommands {
		if re.MatchString(cmdStr) {
			return nil
		}
	}

	return fmt.Errorf("command not allowed: %q does not match any allowed pattern", cmdStr)
}

// executeCommand runs the requested command and returns the response.
func executeCommand(ctx context.Context, req *ExecRequest, allowedCommands []*regexp.Regexp) ExecResponse {
	if req.Command == "" {
		return ExecResponse{Error: "command is required"}
	}

	if err := validateCommand(req, allowedCommands); err != nil {
		return ExecResponse{Error: err.Error()}
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, req.Command, req.Args...)

	// Set environment variables
	if len(req.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range req.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	if req.Dir != "" {
		cmd.Dir = req.Dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	resp := ExecResponse{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
		} else {
			resp.ExitCode = -1
			resp.Error = err.Error()
		}
	}

	return resp
}
