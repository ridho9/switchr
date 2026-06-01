package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
)

var (
	leftStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Magenta)

	rightStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Cyan)

	attachedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Green).
			Bold(true)

	detachedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack)

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Green)

	stoppedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack)

	highlightStyle = lipgloss.NewStyle().
			Background(lipgloss.Blue)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Cyan)
)

func swHelp() help.Model {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Cyan)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Cyan)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	return h
}
