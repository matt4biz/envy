package main

import (
	"fmt"
)

type VersionCommand struct {
	*App
}

func (cmd *VersionCommand) Run() int {
	fmt.Fprintln(cmd.stdout, cmd.version)
	return 0
}
