package main

import (
	"flag"
	"fmt"
	"os"
)

type ReadCommand struct {
	*App
}

func (cmd *ReadCommand) Run() int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	clear := fs.Bool("clear", false, "overwrite contents")

	fs.Usage = cmd.usage

	if err := fs.Parse(cmd.args); err != nil {
		cmd.usage()
		return 1
	}

	cmd.args = fs.Args()

	if len(cmd.args) < 2 {
		cmd.usage()
		return 1
	}

	reader := cmd.stdin

	if cmd.args[1] != "-" {
		file, err := os.Open(cmd.args[1])

		if err != nil {
			fmt.Fprintf(cmd.stderr, "read: %s\n", err)
			return -1
		}

		defer file.Close()

		reader = file
	}

	if *clear {
		if err := cmd.Purge(cmd.args[0]); err != nil {
			fmt.Fprintf(cmd.stderr, "read: %s\n", err)
			return -1
		}
	}

	err := cmd.Read(reader, cmd.args[0])

	if err != nil {
		fmt.Fprintf(cmd.stderr, "read: %s\n", err)
		return -1
	}

	return 0
}
