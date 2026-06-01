package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	sessions      []session
	err           error
	cursor        int
	scrollOffset  int
	width         int
	height        int
	selectedIndex int
	printMode     bool
	help          help.Model
	spinner       spinner.Model

	refreshing bool

	treeData    *sessionTreeData
	treeLoading bool
	treeErr     error
	treeFor     string
}

func (m model) Init() tea.Cmd {
	return m.refreshSessions()
}

// refreshSessions starts the spinner and reloads the session list.
func (m *model) refreshSessions() tea.Cmd {
	m.err = nil
	m.refreshing = true
	return tea.Sequence(m.spinner.Tick, loadSessions)
}

func (m model) Update(msg tea.Msg) (next tea.Model, cmd tea.Cmd) {
	skipSpinner := false
	defer func() {
		if skipSpinner {
			return
		}

		var spinCmd tea.Cmd
		if t, ok := msg.(spinner.TickMsg); ok {
			m.spinner, spinCmd = m.spinner.Update(t)
		}
		if m.sessions != nil && !m.refreshing {
			spinCmd = nil
		}

		next = m
		cmd = tea.Batch(cmd, spinCmd)
	}()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.Quit):
			m.selectedIndex = -1
			skipSpinner = true
			return m, tea.Quit
		case key.Matches(msg, defaultKeyMap.Refresh):
			return m, m.refreshSessions()
		case key.Matches(msg, defaultKeyMap.Select):
			if len(m.sessions) > 0 && m.cursor < len(m.sessions) {
				if m.printMode {
					m.selectedIndex = m.cursor
					skipSpinner = true
					return m, tea.Quit
				}
				return m, attachSessionCmd(m.sessions[m.cursor].Name)
			}
		case key.Matches(msg, defaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
				m.scrollToCursor()
			}
			return m, m.loadTreeIfNeeded()
		case key.Matches(msg, defaultKeyMap.Down):
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
				m.scrollToCursor()
			}
			return m, m.loadTreeIfNeeded()
		default:
			if n, ok := keyToIndex(msg.String()); ok && n >= 0 && n < len(m.sessions) {
				m.cursor = n
				m.scrollToCursor()
				return m, m.loadTreeIfNeeded()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.snapScrollUp()
		m.scrollToCursor()

	case sessionListMsg:
		m.refreshing = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.sessions = msg.sessions
			m.clampBounds()
		}
		return m, m.loadTreeIfNeeded()

	case sessionFinishedMsg:
		m.treeData = nil
		m.treeFor = ""
		m.treeLoading = false
		m.treeErr = nil
		m.scrollOffset = 0
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

func (m model) panelContentHeight() int {
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(defaultKeyMap))
	helpHeight := lipgloss.Height(helpView)
	ph := m.height - helpHeight
	if ph < 1 {
		ph = 1
	}
	return ph
}

// scrollToCursor adjusts scrollOffset so the cursor is visible.
func (m *model) scrollToCursor() {
	vc := max(m.visibleSessionCount()-2, 1) // indicators may take up to 2 lines
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+vc {
		m.scrollOffset = m.cursor - vc + 1
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// snapScrollUp reduces scrollOffset after the terminal grows, so the top
// indicator disappears when there's now enough room to show earlier sessions.
func (m *model) snapScrollUp() {
	vc := max(m.visibleSessionCount()-2, 1)
	potential := max(0, m.cursor-vc+1)
	if potential < m.scrollOffset {
		m.scrollOffset = potential
	}
}

// visibleSessionCount returns how many session lines fit in the left panel.
func (m model) visibleSessionCount() int {
	n := m.panelContentHeight() - borderWidth - 3 // header = "Herdr:\n\n"
	if n < 1 {
		n = 1
	}
	return n
}

// clampBounds resets cursor and scrollOffset if the session list shrunk.
func (m *model) clampBounds() {
	if len(m.sessions) == 0 {
		m.cursor = 0
		m.scrollOffset = 0
		return
	}
	if m.cursor >= len(m.sessions) {
		m.cursor = len(m.sessions) - 1
	}
	maxOffset := max(0, len(m.sessions)-1)
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
}

func (m model) View() tea.View {
	colWidth := m.width / 2
	m.help.SetWidth(m.width)
	panelHeight := m.panelContentHeight()
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(defaultKeyMap))

	left := leftStyle.
		Width(colWidth).
		Height(panelHeight).
		Render(m.renderLeft(colWidth, panelHeight))

	right := rightStyle.
		Width(colWidth).
		Height(panelHeight).
		Render(m.renderRight(panelHeight))

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	full := lipgloss.JoinVertical(lipgloss.Left, content, helpView)

	v := tea.NewView(full)
	v.AltScreen = true
	v.WindowTitle = "switchr"
	return v
}
