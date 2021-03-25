package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GetCommand struct {
	*App
}

func (cmd *GetCommand) Run() int {
	if len(cmd.args) < 1 {
		cmd.usage()
		return 1
	}

	var err error

	parts := strings.Split(cmd.args[0], "/")

	if len(parts) == 1 {
		var m json.RawMessage

		if m, err = cmd.FetchAsJSON(cmd.args[0]); err == nil {
			fmt.Fprintln(cmd.stdout, string(m))
		}
	} else {
		var s string

		if s, err = cmd.Get(parts[0], parts[1]); err == nil {
			fmt.Fprintln(cmd.stdout, s)
		}
	}

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	return 0
}
