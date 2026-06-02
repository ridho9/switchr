//go:build ignore

package main

import "path/filepath"

// Screenshot-only mock data for README captures.
// Copy this file to the repo root and wire it in if you need the screenshot fixture again.
const screenshotMockEnabled = true

func mockSessionListMsg() sessionListMsg {
	return sessionListMsg{sessions: mockSessions()}
}

func mockSessions() []session {
	home := mockHome()
	sessionDir := func(name string) string {
		return filepath.Join(home, ".config", "herdr", "sessions", name)
	}

	return []session{
		{
			Default:    true,
			Name:       "switchr",
			Running:    true,
			SessionDir: sessionDir("switchr"),
			SocketPath: filepath.Join(sessionDir("switchr"), "herdr.sock"),
			Attached:   true,
		},
		{
			Name:       "build-systems",
			Running:    true,
			SessionDir: sessionDir("build-systems"),
			SocketPath: filepath.Join(sessionDir("build-systems"), "herdr.sock"),
		},
		{
			Name:       "docs",
			Running:    false,
			SessionDir: sessionDir("docs"),
			SocketPath: filepath.Join(sessionDir("docs"), "herdr.sock"),
		},
		{
			Name:       "infra",
			Running:    true,
			SessionDir: sessionDir("infra"),
			SocketPath: filepath.Join(sessionDir("infra"), "herdr.sock"),
		},
	}
}

func mockSessionTree(sessionName string) *sessionTreeData {
	home := mockHome()
	path := func(parts ...string) string {
		elems := append([]string{home}, parts...)
		return filepath.Join(elems...)
	}

	pfx := sessionName + "-"

	switch sessionName {
	case "build-systems":
		return &sessionTreeData{
			Workspaces: []workspaceData{
				{WorkspaceID: pfx + "ws-1", Label: "build-systems", Focused: true},
				{WorkspaceID: pfx + "ws-2", Label: "api", Focused: false},
				{WorkspaceID: pfx + "ws-3", Label: "tests", Focused: false},
				{WorkspaceID: pfx + "ws-4", Label: "infra", Focused: false},
			},
			Tabs: []tabData{
				{TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Label: "server", Focused: true},
				{TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Label: "notes", Focused: false},
				{TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Label: "routes", Focused: false},
				{TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Label: "qa", Focused: false},
				{TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Label: "deploy", Focused: false},
			},
			Panes: []paneData{
				{PaneID: pfx + "pane-1", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Work", "build-systems"), Focused: true},
				{PaneID: pfx + "pane-2", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Work", "build-systems", "internal", "server"), Focused: false},
				{PaneID: pfx + "pane-3", TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Cwd: path("Notes", "build-systems-plan.md"), Focused: false},
				{PaneID: pfx + "pane-4", TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Cwd: path("Work", "build-systems", "internal", "api"), Focused: false},
				{PaneID: pfx + "pane-5", TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Cwd: path("Work", "build-systems", "tests"), Focused: false},
				{PaneID: pfx + "pane-6", TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Cwd: path("Library", "Logs", "herdr"), Focused: false},
			},
		}

	case "docs":
		return &sessionTreeData{
			Workspaces: []workspaceData{
				{WorkspaceID: pfx + "ws-1", Label: "docs", Focused: true},
				{WorkspaceID: pfx + "ws-2", Label: "preview", Focused: false},
				{WorkspaceID: pfx + "ws-3", Label: "archive", Focused: false},
			},
			Tabs: []tabData{
				{TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Label: "write", Focused: true},
				{TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Label: "review", Focused: false},
				{TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Label: "screenshots", Focused: false},
				{TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Label: "notes", Focused: false},
			},
			Panes: []paneData{
				{PaneID: pfx + "pane-1", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Notes", "switchr-readme.md"), Focused: true},
				{PaneID: pfx + "pane-2", TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Cwd: path("Project", "switchr", "README.md"), Focused: false},
				{PaneID: pfx + "pane-3", TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Cwd: path("Project", "switchr", "assets", "screenshot.png"), Focused: false},
				{PaneID: pfx + "pane-4", TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Cwd: path("Notes", "release-v0.2.0.md"), Focused: false},
			},
		}

	case "infra":
		return &sessionTreeData{
			Workspaces: []workspaceData{
				{WorkspaceID: pfx + "ws-1", Label: "infra", Focused: true},
				{WorkspaceID: pfx + "ws-2", Label: "logs", Focused: false},
				{WorkspaceID: pfx + "ws-3", Label: "deploy", Focused: false},
				{WorkspaceID: pfx + "ws-4", Label: "alerts", Focused: false},
			},
			Tabs: []tabData{
				{TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Label: "ops", Focused: true},
				{TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Label: "runbook", Focused: false},
				{TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Label: "journal", Focused: false},
				{TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Label: "deploy", Focused: false},
				{TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Label: "watch", Focused: false},
			},
			Panes: []paneData{
				{PaneID: pfx + "pane-1", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Work", "infra"), Focused: true},
				{PaneID: pfx + "pane-2", TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Cwd: path("Work", "infra", "runbook.md"), Focused: false},
				{PaneID: pfx + "pane-3", TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Cwd: path("Library", "Logs", "herdr", "server.log"), Focused: false},
				{PaneID: pfx + "pane-4", TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Cwd: path("Work", "infra", "deploy"), Focused: false},
				{PaneID: pfx + "pane-5", TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Cwd: path("Library", "Logs", "herdr", "alerts.log"), Focused: false},
			},
		}

	default:
		return &sessionTreeData{
			Workspaces: []workspaceData{
				{WorkspaceID: pfx + "ws-1", Label: "switchr", Focused: true},
				{WorkspaceID: pfx + "ws-2", Label: "build-systems", Focused: false},
				{WorkspaceID: pfx + "ws-3", Label: "docs", Focused: false},
				{WorkspaceID: pfx + "ws-4", Label: "ops", Focused: false},
			},
			Tabs: []tabData{
				{TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Label: "main", Focused: true},
				{TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Label: "notes", Focused: false},
				{TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Label: "backend", Focused: false},
				{TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Label: "drafts", Focused: false},
				{TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Label: "logs", Focused: false},
			},
			Panes: []paneData{
				{PaneID: pfx + "pane-1", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Project", "switchr"), Focused: true},
				{PaneID: pfx + "pane-2", TabID: pfx + "ws-1:1", WorkspaceID: pfx + "ws-1", Cwd: path("Project", "switchr", "README.md"), Focused: false},
				{PaneID: pfx + "pane-3", TabID: pfx + "ws-1:2", WorkspaceID: pfx + "ws-1", Cwd: path("Notes", "release-v0.2.0.md"), Focused: false},
				{PaneID: pfx + "pane-4", TabID: pfx + "ws-2:1", WorkspaceID: pfx + "ws-2", Cwd: path("Work", "build-systems"), Focused: false},
				{PaneID: pfx + "pane-5", TabID: pfx + "ws-3:1", WorkspaceID: pfx + "ws-3", Cwd: path("Notes", "switchr-readme.md"), Focused: false},
				{PaneID: pfx + "pane-6", TabID: pfx + "ws-4:1", WorkspaceID: pfx + "ws-4", Cwd: path("Library", "Logs", "herdr"), Focused: false},
			},
		}
	}
}

func mockHome() string {
	return "/Users/you"
}
