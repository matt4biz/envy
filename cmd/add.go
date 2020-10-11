package main

import (
	"fmt"
)

type AddCommand struct {
	App
}

func (cmd *AddCommand) Run(a App) int {
	cmd.App = a

	fmt.Fprintln(cmd.stdout, "user dir =", cmd.path)
	fmt.Fprintln(cmd.stderr, "not adding today")
	return 0
}

func init() {
	commands["add"] = &AddCommand{}
}
