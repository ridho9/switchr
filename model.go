package main

import (
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
				return m, attachSessionCmd(m.sessions[m.cursor].Name)
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
		return m, loadSessions

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
