package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func (m *model) loadTreeIfNeeded() tea.Cmd {
	if len(m.sessions) == 0 || m.cursor >= len(m.sessions) {
		return nil
	}
	sess := m.sessions[m.cursor]
	if !sess.Running {
		return nil
	}
	if m.treeLoading && m.treeFor == sess.Name {
		return nil
	}
	if m.treeData != nil && m.treeFor == sess.Name && m.treeErr == nil {
		return nil
	}
	m.treeLoading = true
	m.treeFor = sess.Name
	m.treeErr = nil
	m.treeData = nil
	return loadSessionTree(sess.Name)
}

func (m model) renderRight() string {
	if len(m.sessions) == 0 || m.cursor >= len(m.sessions) {
		return "Session Info\n\nNo session available."
	}

	sess := m.sessions[m.cursor]
	header := fmt.Sprintf("Session: %s\n", sess.Name)

	if !sess.Running {
		return header + "\nNot running."
	}

	if m.treeLoading || (m.treeData == nil && m.treeErr == nil) {
		return header + "\nLoading..."
	}

	if m.treeErr != nil {
		return header + fmt.Sprintf("\nError: %v", m.treeErr)
	}

	contentWidth := m.width/2 - 2 // normal border: 1 char each side
	return header + "\n" + m.buildTree(contentWidth)
}

func (m model) buildTree(width int) string {
	data := m.treeData
	if data == nil {
		return ""
	}

	tabsByWS := make(map[string][]tabData)
	for _, tab := range data.Tabs {
		tabsByWS[tab.WorkspaceID] = append(tabsByWS[tab.WorkspaceID], tab)
	}

	panesByTab := make(map[string][]paneData)
	for _, pane := range data.Panes {
		panesByTab[pane.TabID] = append(panesByTab[pane.TabID], pane)
	}

	t := tree.New()
	if width > 0 {
		t.Width(width)
	}

	t.EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.BrightBlack))
	t.IndenterStyle(lipgloss.NewStyle().Foreground(lipgloss.BrightBlack))
	t.ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style {
		if strings.HasPrefix(children.At(i).Value(), "* ") {
			return lipgloss.NewStyle().Foreground(lipgloss.Cyan).Bold(true)
		}
		return lipgloss.NewStyle()
	})

	for _, ws := range data.Workspaces {
		wsTree := tree.Root(focusLabel(ws.Focused, ws.Label))

		for _, tab := range tabsByWS[ws.WorkspaceID] {
			tabTree := tree.Root(focusLabel(tab.Focused, fmt.Sprintf("Tab %s", tab.Label)))

			for _, pane := range panesByTab[tab.TabID] {
				tabTree.Child(focusLabel(pane.Focused, shortenPath(pane.Cwd)))
			}

			wsTree.Child(tabTree)
		}

		t.Child(wsTree)
	}

	return t.String()
}

func focusLabel(focused bool, label string) string {
	if focused {
		return "* " + label
	}
	return "  " + label
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
