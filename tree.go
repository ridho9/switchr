package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

// borderWidth is the space taken by NormalBorder (1 cell each side).
const borderWidth = 2

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

func statusParts(sess session) (string, lipgloss.Style) {
	if sess.Attached {
		return "(attached)", attachedStyle
	}
	return "(detached)", detachedStyle
}

func (m model) renderLeft(colWidth, panelHeight int) string {
	contentWidth := colWidth - borderWidth
	s := "Herdr:\n\n"

	if m.err != nil {
		return s + fmt.Sprintf("Error: %v", m.err)
	}

	if m.sessions == nil {
		return "Herdr: " + spinnerStyle.Render(m.spinner.View()) + "\n\n"
	}
	if m.refreshing {
		s = "Herdr: " + spinnerStyle.Render(m.spinner.View()) + "\n\n"
	}
	if len(m.sessions) == 0 {
		return s + "No sessions."
	}

	contentHeight := panelHeight - borderWidth
	maxVisible := contentHeight - 3 // header = 3 lines
	if maxVisible < 1 {
		maxVisible = 1
	}

	so := m.scrollOffset
	maxOffset := max(0, len(m.sessions)-1)
	if so > maxOffset {
		so = maxOffset
	}
	if so < 0 {
		so = 0
	}

	hasAbove := so > 0
	hasBelow := so+maxVisible < len(m.sessions)

	visibleCount := maxVisible
	if hasAbove {
		visibleCount--
	}
	if hasBelow {
		visibleCount--
	}
	if visibleCount < 1 {
		visibleCount = 1
	}

	end := min(so+visibleCount, len(m.sessions))

	if hasAbove {
		s += "  \u00b7\u00b7\u00b7\n"
	}

	for i := so; i < end; i++ {
		sess := m.sessions[i]
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		icon := "\u25cb"
		iconStyle := stoppedStyle
		if sess.Running {
			icon = "\u25cf"
			iconStyle = runningStyle
		}

		prefix := cursor + "[" + indexToLabel(i) + "] "
		styledIcon := iconStyle.Render(icon)
		leftPart := prefix + styledIcon + " " + sess.Name

		statusText, statusStyle := statusParts(sess)

		if m.cursor == i {
			// Plain text (no color) so ANSI reset doesn't clear highlight bg.
			plainLeft := prefix + icon + " " + sess.Name
			line := lipgloss.NewStyle().Width(contentWidth).Render(plainLeft)
			if sess.Running {
				remaining := max(contentWidth-lipgloss.Width(plainLeft), 1)
				line = plainLeft + lipgloss.NewStyle().
					Width(remaining).
					Align(lipgloss.Right).
					Render(statusText)
			}
			line = highlightStyle.Width(contentWidth).Render(line)
			s += line + "\n"
		} else {
			if sess.Running {
				remaining := max(contentWidth-lipgloss.Width(leftPart), 1)
				s += leftPart + statusStyle.
					Width(remaining).
					Align(lipgloss.Right).
					Render(statusText) + "\n"
			} else {
				s += lipgloss.NewStyle().Width(contentWidth).Render(leftPart) + "\n"
			}
		}
	}

	if hasBelow {
		s += "  \u00b7\u00b7\u00b7\n"
	}

	return s
}

func (m model) renderRight(panelHeight int) string {
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

	contentWidth := m.width/2 - borderWidth
	treeStr := m.buildTree(contentWidth)

	// Truncate tree if it exceeds available height.
	headerLines := 2 // "Session: name\n" + blank line
	contentHeight := panelHeight - borderWidth
	maxTreeLines := contentHeight - headerLines
	if maxTreeLines < 1 {
		maxTreeLines = 1
	}

	lines := strings.Split(treeStr, "\n")
	if len(lines) > maxTreeLines {
		lines = lines[:maxTreeLines-1]
		lines = append(lines, "\u00b7\u00b7\u00b7")
		treeStr = strings.Join(lines, "\n")
	}

	return header + "\n" + treeStr
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
