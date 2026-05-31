package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	leftBorderColor  = lipgloss.Color("63")  // blue
	rightBorderColor = lipgloss.Color("204") // pink
	cursorColor      = lipgloss.Color("212") // light pink
)

var (
	leftPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(leftBorderColor)

	rightPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(rightBorderColor)

	cursorStyle = lipgloss.NewStyle().
			Foreground(cursorColor).
			Bold(true)
)

type model struct {
	items    []string
	details  []string
	cursor   int
	selected int
	width    int
	height   int
}

func initialModel() model {
	return model{
		items: []string{
			"Project Alpha",
			"Design System",
			"API Gateway",
			"Auth Service",
			"Dashboard",
		},
		details: []string{
			"Frontend rewrite in React\nStatus: In Progress",
			"Shared component library\nStatus: Review",
			"Rate limiting & routing\nStatus: Planning",
			"OAuth2 + JWT tokens\nStatus: Done",
			"Admin analytics panel\nStatus: Backlog",
		},
		selected: -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			m.selected = m.cursor
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() tea.View {
	colWidth := m.width / 2

	left := leftPanelStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(m.renderList())

	right := rightPanelStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(m.renderDetail())

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m model) renderList() string {
	s := "Contexts\n\n"
	for i, item := range m.items {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		line := fmt.Sprintf("%s%s\n", cursor, item)
		if m.cursor == i {
			line = cursorStyle.Render(line)
		}
		s += line
	}
	s += "\n↑↓ navigate  ↵ switch  q quit"
	return s
}

func (m model) renderDetail() string {
	if m.selected == -1 {
		return "Select a context\nto see details."
	}
	return fmt.Sprintf("%s\n\n%s", m.items[m.selected], m.details[m.selected])
}
