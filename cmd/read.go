package main

import (
	"fmt"
	"io/ioutil"
)

type ReadCommand struct {
	*App
}

func (cmd *ReadCommand) Run() int {
	if len(cmd.args) < 2 {
		cmd.usage()
		return 1
	}

	m, err := cmd.FetchAsJSON(cmd.args[0])

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	// it's nice to have the file (or stdout)
	// have a trailing newline

	m = append(m, '\n')

	if cmd.args[1] != "-" {
		if err := ioutil.WriteFile(cmd.args[1], m, 0600); err != nil {
			fmt.Fprintln(cmd.stderr, err)
			return -1
		}
	} else {
		if _, err := cmd.stdout.Write(m); err != nil {
			fmt.Fprintln(cmd.stderr, err)
			return -1
		}
	}

	return 0
}
