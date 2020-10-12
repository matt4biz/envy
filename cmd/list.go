package main

import (
	"flag"
	"fmt"
)

type ListCommand struct {
	*App
}

func (cmd *ListCommand) Run(app *App) int {
	cmd.App = app

	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	decrypt := fs.Bool("d", false, "")

	if err := fs.Parse(cmd.args); err != nil {
		fmt.Fprintln(cmd.stderr, usage())
		return 1
	}

	cmd.args = fs.Args()

	if len(cmd.args) < 1 {
		realms, err := cmd.Realms()

		if err != nil {
			fmt.Fprintln(cmd.stderr, err)
			return -1
		}

		for _, r := range realms {
			fmt.Fprintln(cmd.stdout, r)
		}

		return 0
	}

	err := cmd.List(cmd.stdout, cmd.args[0], *decrypt)

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}

func init() {
	commands["list"] = &ListCommand{}
}
