package main

import (
	"fmt"

	"github.com/matt4biz/envy/internal"
)

type AddCommand struct {
	*App
}

func (cmd *AddCommand) Run() int {
	e, err := internal.NewExtractor(cmd.args)

	if err != nil {
		fmt.Fprintf(cmd.stderr, "extract: %s\n", err)
		return -1
	}

	cmd.args = e.Args()

	if err = cmd.Add(e.Realm(), e.Values()); err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}
