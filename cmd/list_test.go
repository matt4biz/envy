package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)
	cmd := ListCommand{}
	data := map[string]string{"a": "xxx", "b": "yyy"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	app.args = []string{"top"}

	o := cmd.Run(app)

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	lines := strings.Split(stdout.String(), "\n")

	// trailing '\n' makes an extra (blank) line

	if len(lines) != 3 {
		t.Fatalf("invalid count: %v", lines)
	}

	if lines[0][0] != 'a' || lines[1][0] != 'b' {
		t.Fatalf("invalid output: %q", lines)
	}
}
