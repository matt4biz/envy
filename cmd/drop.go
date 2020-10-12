package main

import (
	"fmt"
	"strings"
)

type DropCommand struct {
	*App
}

func (cmd *DropCommand) Run() int {
	if len(cmd.args) < 1 {
		cmd.usage()
		return 1
	}

	var err error

	parts := strings.Split(cmd.args[0], "/")

	if len(parts) == 1 {
		err = cmd.Purge(cmd.args[0])
	} else {
		err = cmd.Drop(parts[0], parts[1])
	}

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}
