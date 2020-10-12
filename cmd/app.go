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

	fs.Usage = a.usage

	if err := fs.Parse(args); err != nil {
		return ErrUsage
	} else if *help {
		a.usage()
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

func (a *App) usage() {
	msg := strings.TrimSpace(`
envy: a tool to securely store and retrieve environment variables.

Variables are key-value pairs stored in a "realm" (or "namespace") of which 
there may be one or more. All data is stored in a DB within the user's "config" 
directory, encrypted with a per-user secret key stored in the system keychain.

All operations take place in one of five subcommands. Add will create a realm
if it doesn't exist, or overwrite keys in a realm that already exists. Drop
may be used to delete one key or an entire realm. Exec will execute a command
with arguments, with value(s) from the realm injected as environment variables.

Usage:
  -h show this help message

  add       realm key=value [key=value ...]
  drop      realm[/key]
  exec      realm[/key] command [args...]
  list [-d] realm
    -d decrypt & show secrets also
  version

Listing a realm displays a timestamp, size, and hash for each key-value pair.
	`)

	fmt.Fprintln(a.stderr, msg)
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
