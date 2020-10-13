package main

import (
	"fmt"
)

type VersionCommand struct {
	*App
}

func (a *VersionCommand) NeedsDB() bool {
	return false
}

func (cmd *VersionCommand) Run() int {
	fmt.Fprintln(cmd.stdout, cmd.version)
	return 0
}
