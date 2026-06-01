package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// restartModal wraps the main model when herdr server is incompatible.
// It blocks all interaction except y/n/q and renders a centered modal.
type restartModal struct {
	inner     model
	serverVer string
	clientVer string
}

func (m restartModal) Init() tea.Cmd {
	return nil
}

func (m restartModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch strings.ToLower(msg.String()) {
		case "y":
			return m, restartHerderDaemon
		case "n":
			return m.inner, m.inner.loadTreeIfNeeded()
		case "q", "ctrl+c":
			m.inner.selectedIndex = -1
			return m.inner, tea.Quit
		}
		return m, nil

	case sessionListMsg:
		next, cmd := m.inner.Update(msg)
		if inner, ok := next.(model); ok {
			m.inner = inner
			if !inner.server.restartNeeded {
				inner.server.justRestarted = true
				return inner, cmd
			}
			return m, cmd
		}
		return next, cmd

	default:
		next, cmd := m.inner.Update(msg)
		if inner, ok := next.(model); ok {
			m.inner = inner
			return m, cmd
		}
		return next, cmd
	}
}

func (m restartModal) View() tea.View {
	v := tea.NewView(m.renderModal())
	v.AltScreen = true
	v.WindowTitle = "switchr"
	return v
}

func (m restartModal) renderModal() string {
	w, h := m.inner.width, m.inner.height
	if w < 40 {
		w = 80
	}
	if h < 10 {
		h = 24
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Yellow).
		Padding(1, 2)

	lines := []string{
		warnStyle.Render("⚠ Herdr server is incompatible"),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.BrightBlack).Render(
			"server: v" + m.serverVer + "    client: v" + m.clientVer),
		"",
		"[y] Restart server",
		"[n] Dismiss",
		"[q] Quit switchr",
	}

	content := strings.Join(lines, "\n")
	styled := box.Render(content)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, styled)
}
