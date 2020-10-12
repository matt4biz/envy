package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"top", "../test/test.sh"}

	cmd := ExecCommand{app}
	data := map[string]string{"a": "b", "b": "1"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 1st return: %d", o)
	} else {
		t.Log(strings.TrimSpace(stdout.String()))
	}

	lines := strings.Split(stdout.String(), "\n")

	// trailing '\n' makes an extra (blank) line

	if len(lines) != 2 {
		t.Fatalf("invalid 1st count: %v", lines)
	}

	if lines[0] != "b 1" {
		t.Fatalf("invalid 1st output: %q", lines)
	}
}

func TestExecOne(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"top/b", "../test/test.sh"}

	cmd := ExecCommand{app}
	data := map[string]string{"a": "b", "b": "1"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 2nd return: %d", o)
	} else {
		t.Log(strings.TrimSpace(stdout.String()))
	}

	lines := strings.Split(stdout.String(), "\n")

	// trailing '\n' makes an extra (blank) line

	if len(lines) != 2 {
		t.Fatalf("invalid 2nd count: %v", lines)
	}

	if lines[0] != "1" {
		t.Fatalf("invalid 2nd output: %q", lines)
	}
}
