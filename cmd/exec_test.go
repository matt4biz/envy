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
	cmd := ExecCommand{}
	data := map[string]string{"a": "XX"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	app.args = []string{"top", "../test/test.sh"}

	o := cmd.Run(app)

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	} else {
		t.Log(strings.TrimSpace(stdout.String()))
	}

	lines := strings.Split(stdout.String(), "\n")

	// trailing '\n' makes an extra (blank) line

	if len(lines) != 2 {
		t.Fatalf("invalid count: %v", lines)
	}

	if lines[0] != "XX" {
		t.Fatalf("invalid output: %q", lines)
	}
}
