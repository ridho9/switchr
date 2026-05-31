package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// --- session tree data ---

type workspaceData struct {
	WorkspaceID string `json:"workspace_id"`
	Label       string `json:"label"`
	Focused     bool   `json:"focused"`
}

type tabData struct {
	TabID       string `json:"tab_id"`
	WorkspaceID string `json:"workspace_id"`
	Label       string `json:"label"`
	Focused     bool   `json:"focused"`
}

type paneData struct {
	PaneID      string `json:"pane_id"`
	TabID       string `json:"tab_id"`
	WorkspaceID string `json:"workspace_id"`
	Cwd         string `json:"cwd"`
	Focused     bool   `json:"focused"`
}

type sessionTreeData struct {
	Workspaces []workspaceData
	Tabs       []tabData
	Panes      []paneData
}

type treeDataMsg struct {
	sessionName string
	data        *sessionTreeData
	err         error
}

func fetchHerdrList[T any](sessionName, subcommand, jsonKey string) ([]T, error) {
	out, err := exec.Command("herdr", "--session", sessionName, subcommand, "list").Output()
	if err != nil {
		return nil, fmt.Errorf("%s list: %w", subcommand, err)
	}
	var outer struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(out, &outer); err != nil {
		return nil, fmt.Errorf("parse %s: %w", subcommand, err)
	}
	var inner map[string]json.RawMessage
	if err := json.Unmarshal(outer.Result, &inner); err != nil {
		return nil, fmt.Errorf("parse %s result: %w", subcommand, err)
	}
	raw, ok := inner[jsonKey]
	if !ok {
		return nil, fmt.Errorf("key %q not found in %s result", jsonKey, subcommand)
	}
	var items []T
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse %s items: %w", subcommand, err)
	}
	return items, nil
}

func loadSessionTree(sessionName string) tea.Cmd {
	return func() tea.Msg {
		type fetchPair[T any] struct {
			items []T
			err   error
		}

		wsCh := make(chan fetchPair[workspaceData], 1)
		tabCh := make(chan fetchPair[tabData], 1)
		paneCh := make(chan fetchPair[paneData], 1)

		go func() {
			items, err := fetchHerdrList[workspaceData](sessionName, "workspace", "workspaces")
			wsCh <- fetchPair[workspaceData]{items, err}
		}()
		go func() {
			items, err := fetchHerdrList[tabData](sessionName, "tab", "tabs")
			tabCh <- fetchPair[tabData]{items, err}
		}()
		go func() {
			items, err := fetchHerdrList[paneData](sessionName, "pane", "panes")
			paneCh <- fetchPair[paneData]{items, err}
		}()

		ws := <-wsCh
		tab := <-tabCh
		pane := <-paneCh

		var errs []string
		if ws.err != nil {
			errs = append(errs, ws.err.Error())
		}
		if tab.err != nil {
			errs = append(errs, tab.err.Error())
		}
		if pane.err != nil {
			errs = append(errs, pane.err.Error())
		}

		if len(errs) > 0 {
			return treeDataMsg{sessionName: sessionName, err: fmt.Errorf("%s", strings.Join(errs, "; "))}
		}

		return treeDataMsg{
			sessionName: sessionName,
			data: &sessionTreeData{
				Workspaces: ws.items,
				Tabs:       tab.items,
				Panes:      pane.items,
			},
		}
	}
}
