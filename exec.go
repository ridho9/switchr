package main

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
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

func attachSessionCmd(name string) tea.Cmd {
	cmd := &titledCmd{
		Cmd:   exec.Command("herdr", "session", "attach", name),
		title: fmt.Sprintf("herdr: %s", name),
	}
	return tea.Exec(cmd, func(err error) tea.Msg {
		return sessionFinishedMsg{}
	})
}

func restartHerderDaemon() tea.Msg {
	exec.Command("herdr", "server", "stop").Run()
	cmd := exec.Command("herdr", "server")
	cmd.Start() // background: runs headless, doesn't attach

	// Poll until the new server is ready to accept commands.
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		if err := exec.Command("herdr", "status", "--json").Run(); err == nil {
			break
		}
	}

	return loadSessions()
}
