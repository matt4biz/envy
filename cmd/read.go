package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

type ReadCommand struct {
	*App
}

func (cmd *ReadCommand) Run() int {
	fs := flag.NewFlagSet("read", flag.ContinueOnError)
	unquote := fs.Bool("q", false, "unquote embedded JSON")

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

	m, err := cmd.FetchAsJSON(cmd.args[0])

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	if *unquote {
		// we're going to have to do this the hard way

		v := string(m)

		v = strings.Trim(v, "\"")
		v = strings.ReplaceAll(v, `\"`, `"`)
		v = strings.ReplaceAll(v, `"{`, `{`)
		v = strings.ReplaceAll(v, `}"`, `}`)
		v = strings.ReplaceAll(v, `\\n`, ``)

		m = []byte(v)
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
