package main

import (
	"flag"
	"fmt"
	"strings"
)

type ListCommand struct {
	*App
}

func (cmd *ListCommand) Run() int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	decrypt := fs.Bool("d", false, "show decrypted values")

	fs.Usage = cmd.usage

	if err := fs.Parse(cmd.args); err != nil {
		cmd.usage()
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

	var err error

	parts := strings.Split(cmd.args[0], "/")

	if len(parts) == 1 {
		err = cmd.List(cmd.stdout, cmd.args[0], "", *decrypt)
	} else {
		err = cmd.List(cmd.stdout, parts[0], parts[1], *decrypt)
	}

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}
