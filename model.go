package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
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

	highlightStyle = lipgloss.NewStyle().
			Background(lipgloss.Blue)
)

type sessionFinishedMsg struct{}

// titledCmd wraps an exec.Cmd to set the terminal title before running.
type titledCmd struct {
	*exec.Cmd
	title string
}

func (c *titledCmd) Run() error {
	fmt.Fprintf(c.Stdout, "\033]0;%s\007", c.title)
	return c.Cmd.Run()
}

func (c *titledCmd) SetStdin(r io.Reader) {
	if c.Cmd.Stdin == nil {
		c.Cmd.Stdin = r
	}
}

func (c *titledCmd) SetStdout(w io.Writer) {
	c.Cmd.Stdout = w
}

func (c *titledCmd) SetStderr(w io.Writer) {
	if c.Cmd.Stderr == nil {
		c.Cmd.Stderr = w
	}
}

type model struct {
	sessions      []session
	err           error
	cursor        int
	width         int
	height        int
	selectedIndex int
	printMode     bool

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
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.selectedIndex = -1
			return m, tea.Quit
		case "enter":
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
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, m.loadTreeIfNeeded()
		case "down", "j":
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

	left := leftStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(m.renderLeft(colWidth))

	right := rightStyle.
		Width(colWidth).
		Height(m.height - 2).
		Render(m.renderRight())

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = m.windowTitle()
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

func (m model) windowTitle() string {
	return "switcher"
}
