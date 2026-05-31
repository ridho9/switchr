package main

import (
	"fmt"
	"os/exec"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	sessions      []session
	err           error
	cursor        int
	width         int
	height        int
	selectedIndex int
	printMode     bool
	help          help.Model

	treeData    *sessionTreeData
	treeLoading bool
	treeErr     error
	treeFor     string
}

func (m model) Init() tea.Cmd {
	return loadSessions
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.Quit):
			m.selectedIndex = -1
			return m, tea.Quit
		case key.Matches(msg, defaultKeyMap.Select):
			if len(m.sessions) > 0 && m.cursor < len(m.sessions) {
				if m.printMode {
					m.selectedIndex = m.cursor
					return m, tea.Quit
				}
				sess := m.sessions[m.cursor]
				cmd := &titledCmd{
					Cmd:   exec.Command("herdr", "session", "attach", sess.Name),
					title: fmt.Sprintf("herdr: %s", sess.Name),
				}
				return m, tea.Exec(cmd, func(err error) tea.Msg {
					return sessionFinishedMsg{}
				})
			}
		case key.Matches(msg, defaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, m.loadTreeIfNeeded()
		case key.Matches(msg, defaultKeyMap.Down):
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}
			return m, m.loadTreeIfNeeded()
		default:
			if n, ok := keyToIndex(msg.String()); ok && n >= 0 && n < len(m.sessions) {
				m.cursor = n
				return m, m.loadTreeIfNeeded()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case sessionListMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.sessions = msg.sessions
		}
		return m, m.loadTreeIfNeeded()

	case sessionFinishedMsg:
		// Session ended; TUI resumes automatically.
		return m, nil

	case treeDataMsg:
		m.treeLoading = false
		if msg.sessionName == m.treeFor {
			if msg.err != nil {
				m.treeErr = msg.err
				m.treeData = nil
			} else {
				m.treeErr = nil
				m.treeData = msg.data
			}
		}
	}

	return m, nil
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

func (m model) View() tea.View {
	colWidth := m.width / 2
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(defaultKeyMap))
	helpHeight := lipgloss.Height(helpView)
	panelHeight := m.height - 2 - helpHeight
	if panelHeight < 1 {
		panelHeight = 1
	}

	m.help.SetWidth(m.width)

	left := leftStyle.
		Width(colWidth).
		Height(panelHeight).
		Render(m.renderLeft(colWidth))

	right := rightStyle.
		Width(colWidth).
		Height(panelHeight).
		Render(m.renderRight())

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	full := lipgloss.JoinVertical(lipgloss.Left, content, helpView)

	v := tea.NewView(full)
	v.AltScreen = true
	v.WindowTitle = "switcher"
	return v
}

func (m model) renderLeft(colWidth int) string {
	contentWidth := colWidth - 2 // normal border takes 1 char each side
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

		status := "(detached)"
		statusStyle := detachedStyle
		if sess.Attached {
			status = "(attached)"
			statusStyle = attachedStyle
		}

		leftPart := cursor + "[" + indexToLabel(i) + "] " + sess.Name
		remaining := max(contentWidth-lipgloss.Width(leftPart), 1)

		if m.cursor == i {
			// Plain style (no color) so ANSI reset doesn't clear highlight bg.
			line := leftPart + lipgloss.NewStyle().
				Width(remaining).
				Align(lipgloss.Right).
				Render(status)
			line = highlightStyle.Width(contentWidth).Render(line)
			s += line + "\n"
		} else {
			s += leftPart + statusStyle.
				Width(remaining).
				Align(lipgloss.Right).
				Render(status) + "\n"
		}
	}

	return s
}
