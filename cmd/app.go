package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/matt4biz/envy"
)

type App struct {
	*envy.Envy

	args    []string
	version string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

type Command interface {
	Run() int
	NeedsDB() bool
}

var (
	ErrUsage          = errors.New("usage")
	ErrUnknownCommand = errors.New("unknown command")
)

// NeedsDB operates on the opt-out theory; a subcommand
// should override this if it doesn't need a DB set up.
func (a *App) NeedsDB() bool {
	return true
}

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

	a.args = fs.Args()
	return nil
}

func (a *App) getCommand() (Command, error) {
	if len(a.args) == 0 {
		return nil, ErrUnknownCommand
	}

	s := a.args[0]
	a.args = a.args[1:]

	switch s {
	case "add":
		return &AddCommand{a}, nil
	case "drop":
		return &DropCommand{a}, nil
	case "exec":
		return &ExecCommand{a}, nil
	case "get":
		return &GetCommand{a}, nil
	case "list":
		return &ListCommand{a}, nil
	case "read":
		return &ReadCommand{a}, nil
	case "version":
		return &VersionCommand{a}, nil
	case "write":
		return &WriteCommand{a}, nil
	}

	return nil, fmt.Errorf("%s: %w", s, ErrUnknownCommand)
}

func (a *App) usage() {
	msg := strings.TrimSpace(`
envy: a tool to securely store and retrieve environment variables.

Variables are key-value pairs stored in a "realm" (or "namespace") of which 
there may be one or more. All data is stored in a DB within the user's "config" 
directory, encrypted with a per-user secret key stored in the system keychain.

All operations take place in one of the subcommands. Add will create a realm
if it doesn't exist, or overwrite keys in a realm that already exists. Drop
may be used to delete one key or an entire realm. Exec will execute a command
with arguments, with value(s) from the realm injected as environment variables.
Get will return the stored value (string for a key, JSON for an entire realm).
Read and write allow a realm's data to be exported or imported in JSON format.

Usage: envy [opts] subcommand
  -h  show this help message and exit

  add          realm       key=value [key=value ...]
  get          realm[/key]
    -n	don't add a trailing newline
  drop         realm[/key]
  exec         realm[/key] command [args ...]
  list  [opts] [realm[/key]]
    -d  show decrypted secrets also
  read  [opts] realm       file ('-' for stdout)
    -q  unquote embedded JSON in values
  write [opts] realm       file ('-' for stdin)
    -clear  overwrite contents
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

	if cmd.NeedsDB() {
		a.Envy, err = envy.New()

		if err != nil {
			fmt.Fprintln(stderr, err)
			return -1
		}

		defer a.Envy.Close()
	}

	return cmd.Run()
}
