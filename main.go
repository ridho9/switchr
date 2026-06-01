package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
)

type selection struct {
	Command     string `json:"command"`
	SessionName string `json:"session_name"`
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	printFlag := flag.Bool("print", false, "Print selection as JSON and exit (non-interactive mode)")
	flag.Parse()

	if *versionFlag {
		fmt.Println("switcher", versionString())
		os.Exit(0)
	}

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot open /dev/tty: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	printMode := *printFlag || !term.IsTerminal(os.Stdout.Fd())

	p := tea.NewProgram(
		model{selectedIndex: -1, printMode: printMode, help: swHelp()},
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
