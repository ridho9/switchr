package main

import (
	"fmt"
	"io"
	"os/exec"
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
