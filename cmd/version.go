package main

import (
	"fmt"
)

type VersionCommand struct {
	*App
}

func (cmd *VersionCommand) Run(app *App) int {
	cmd.App = app

	fmt.Fprintln(cmd.stdout, cmd.version)
	return 0
}

func init() {
	commands["version"] = &VersionCommand{}
}
