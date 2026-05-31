package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	leftStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63"))

	rightStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("204"))
)

type session struct {
	Default    bool   `json:"default"`
	Name       string `json:"name"`
	Running    bool   `json:"running"`
	SessionDir string `json:"session_dir"`
	SocketPath string `json:"socket_path"`
}

type sessionsMsg struct {
	sessions []session
	err      error
}

type model struct {
	sessions []session
	err      error
	cursor   int
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return loadSessions
}

func loadSessions() tea.Msg {
	out, err := exec.Command("herdr", "session", "list", "--json").Output()
	if err != nil {
		return sessionsMsg{err: err}
	}

	var result struct {
		Sessions []session `json:"sessions"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return sessionsMsg{err: err}
	}

	return sessionsMsg{sessions: result.Sessions}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}
		default:
			if n, ok := keyToIndex(msg.String()); ok && n >= 0 && n < len(m.sessions) {
				m.cursor = n
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case sessionsMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.sessions = msg.sessions
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	colWidth := m.width / 2

	left := leftStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(m.renderLeft())

	right := rightStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(fmt.Sprintf("Right Panel\n\nResize to see layout."))

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func keyToIndex(key string) (int, bool) {
	if len(key) != 1 {
		return 0, false
	}
	switch key[0] {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return int(key[0] - '1'), true
	case '0':
		return 9, true
	}
	return 0, false
}

func indexToLabel(i int) string {
	if i == 9 {
		return "0"
	}
	return string(rune('1' + i))
}

func (m model) renderLeft() string {
	s := "Herdr:\n\n"

	if m.err != nil {
		return s + fmt.Sprintf("Error: %v", m.err)
	}
	if m.sessions == nil {
		return s + "Loading..."
	}

	for i, sess := range m.sessions {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s += fmt.Sprintf("%s[%s] - %s\n", cursor, indexToLabel(i), sess.Name)
	}

	return s
}
