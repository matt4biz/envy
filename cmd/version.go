package main

import (
	"fmt"
)

type VersionCommand struct {
	App
}

func (cmd *VersionCommand) Run(a App) int {
	cmd.App = a

	fmt.Fprintln(cmd.stdout, cmd.version)
	return 0
}

func init() {
	commands["version"] = &VersionCommand{}
}
