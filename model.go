package main

import (
	"fmt"

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

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() tea.View {
	colWidth := m.width / 2

	left := leftStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(fmt.Sprintf("Left Panel\n\nPress q to quit."))

	right := rightStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(fmt.Sprintf("Right Panel\n\nResize to see layout."))

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
