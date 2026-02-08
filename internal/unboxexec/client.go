package unboxexec

import (
	"encoding/json"
	"fmt"
	"net"
)

// SendRequest connects to the Unix socket and sends an ExecRequest.
// Returns ExecResponse or an error if connection/protocol fails.
func SendRequest(sockPath string, req *ExecRequest) (*ExecResponse, error) {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var resp ExecResponse
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &resp, nil
}
