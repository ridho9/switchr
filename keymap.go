package main

import "charm.land/bubbles/v2/key"

// keyMap defines the keybindings for the switcher UI.
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

var defaultKeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Quit},
	}
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
