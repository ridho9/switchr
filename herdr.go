package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type session struct {
	Default    bool   `json:"default"`
	Name       string `json:"name"`
	Running    bool   `json:"running"`
	SessionDir string `json:"session_dir"`
	SocketPath string `json:"socket_path"`
	Attached   bool
}

type sessionListMsg struct {
	sessions []session
	err      error
}

func loadSessions() tea.Msg {
	command := exec.Command("herdr", "session", "list", "--json")
	out, err := command.Output()
	if err != nil {
		return sessionListMsg{err: err}
	}

	var result struct {
		Sessions []session `json:"sessions"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return sessionListMsg{err: err}
	}

	for i := range result.Sessions {
		result.Sessions[i].Attached = isSessionAttached(result.Sessions[i].SocketPath)
	}

	return sessionListMsg{sessions: result.Sessions}
}

func isSessionAttached(serverSocket string) bool {
	clientSocket := strings.Replace(serverSocket, "herdr.sock", "herdr-client.sock", 1)
	cmd := exec.Command("lsof", clientSocket)
	out, _ := cmd.Output()
	// Heuristic: >2 lines (header + single server fd = 2 lines) means attached.
	return len(out) > 0 && bytes.Count(out, []byte("\n")) > 2
}
