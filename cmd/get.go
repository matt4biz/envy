package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

type GetCommand struct {
	*App
}

func (cmd *GetCommand) Run() int {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	raw := fs.Bool("n", false, "remove trailing newline")

	fs.Usage = cmd.usage

	if err := fs.Parse(cmd.args); err != nil {
		cmd.usage()
		return 1
	}

	cmd.args = fs.Args()

	if len(cmd.args) < 1 {
		cmd.usage()
		return 1
	}

	var err error

	parts := strings.Split(cmd.args[0], "/")

	if len(parts) == 1 {
		var m json.RawMessage

		if m, err = cmd.FetchAsJSON(cmd.args[0]); err == nil {
			if *raw {
				fmt.Fprint(cmd.stdout, string(m))
			} else {
				fmt.Fprintln(cmd.stdout, string(m))
			}

		}
	} else {
		var s string

		if s, err = cmd.Get(parts[0], parts[1]); err == nil {
			if *raw {
				fmt.Fprint(cmd.stdout, s)
			} else {
				fmt.Fprintln(cmd.stdout, s)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}
