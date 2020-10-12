package main

import (
	"bytes"
	"testing"
)

func TestDrop(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"top/a"}

	cmd := DropCommand{app}
	data := map[string]string{"a": "XX", "b": "YY"}

	if err := app.Add("top", data); err != nil {
		t.Fatal("setup", err)
	}

	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 1st return: %d", o)
	}

	m, err := cmd.Fetch("top")

	if err != nil {
		t.Fatal("fetch", err)
	}

	if m["b"] != "YY" {
		t.Errorf("invalid data: %#v", m)
	}

	app.args = []string{"top"}

	o = cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid 2nd return: %d", o)
	}

	_, err = cmd.Fetch("top")

	if err == nil {
		t.Errorf("no error after purge!")
	}
}
