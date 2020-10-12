package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/matt4biz/envy"
)

type App struct {
	*envy.Envy

	args    []string
	version string
	path    string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

type Command interface {
	Run(*App) int
}

var (
	ErrUsage          = errors.New("usage")
	ErrUnknownCommand = errors.New("unknown command")

	commands = map[string]Command{} // must be ready for init
)

func (a *App) fromArgs(args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	help := fs.Bool("h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	} else if *help {
		fmt.Fprintln(a.stderr, usage())
		return ErrUsage
	}

	// we need e.g.
	// [envy] [opts] command [command-opts] [path] [other args ...]
	// where if the command is set
	// the other args need to be in key=value pairs that we regex

	a.args = fs.Args()
	return nil
}

func (a *App) getCommand() (Command, error) {
	if len(a.args) == 0 {
		return nil, ErrUnknownCommand
	}

	s := a.args[0]
	a.args = a.args[1:]

	c, ok := commands[s]

	if !ok {
		return nil, fmt.Errorf("%s: %w", s, ErrUnknownCommand)
	}

	return c, nil
}

func usage() string {
	return strings.TrimSpace(`
envy: a tool to securely store and retrieve environment variables.

usage:
  -h help
	`)
}

func runApp(args []string, version string, stdin io.Reader, stdout, stderr io.Writer) int {
	a := App{version: version, stdin: stdin, stdout: stdout, stderr: stderr}

	if err := a.fromArgs(args); err != nil {
		if err == ErrUsage {
			return 0
		}

		fmt.Fprintln(stderr, err)
		return 1
	}

	cmd, err := a.getCommand()

	if err != nil {
		fmt.Fprintln(stderr, err)
		return -1
	}

	p, err := os.UserConfigDir()

	if err != nil {
		fmt.Fprintln(stderr, err)
		return -1
	}

	a.path = path.Join(p, "envy")
	a.Envy, err = envy.New(a.path)

	if err != nil {
		fmt.Fprintln(stderr, err)
		return -1
	}

	defer a.Envy.Close()
	return cmd.Run(&a)
}
