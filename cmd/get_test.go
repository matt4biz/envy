package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"top/a"}

	cmd := GetCommand{app}
	data := map[string]string{"a": "XX", "b": "YY"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 1st return: %d", o)
	}

	if s := strings.TrimSpace(stdout.String()); s != "XX" {
		t.Errorf("wrong data, got %q", s)
	}

	stdout.Reset()

	app.args = []string{"top"}
	o = cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 1st return: %d", o)
	}

	if s := strings.TrimSpace(stdout.String()); s != `{"a":"XX","b":"YY"}` {
		t.Errorf("wrong data, got %s", s)
	}
}
