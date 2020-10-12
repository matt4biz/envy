package main

import (
	"fmt"

	"github.com/matt4biz/envy/internal"
)

type AddCommand struct {
	*App
}

func (cmd *AddCommand) Run(app *App) int {
	cmd.App = app

	e, err := internal.NewExtractor(cmd.args)

	if err != nil {
		fmt.Fprintf(cmd.stderr, "extract: %s\n", err)
		return -1
	}

	cmd.args = e.Args()
	//fmt.Fprintln(cmd.stdout, "realm:", e.Realm(), "vals:", e.Values(), "args", cmd.args)

	err = cmd.Add(e.Realm(), e.Values())

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}

func init() {
	commands["add"] = &AddCommand{}
}
