package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type selection struct {
	Command     string `json:"command"`
	SessionName string `json:"session_name"`
}

func main() {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot open /dev/tty: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	p := tea.NewProgram(
		model{selectedIndex: -1},
		tea.WithInput(tty),
		tea.WithOutput(tty),
	)

	m, runErr := p.Run()
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", runErr)
		os.Exit(1)
	}

	final := m.(model)
	if final.selectedIndex >= 0 && final.selectedIndex < len(final.sessions) {
		out, _ := json.Marshal(selection{
			Command:     "open-herdr",
			SessionName: final.sessions[final.selectedIndex].Name,
		})
		fmt.Println(string(out))
	} else {
		os.Exit(1)
	}
}
