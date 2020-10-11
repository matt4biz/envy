package main

import (
	"os"
)

var version string // do not remove or change

func main() {
	os.Exit(runApp(os.Args[1:], version, os.Stdin, os.Stdout, os.Stderr))
}
