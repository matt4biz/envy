package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"top", "a=b"}

	cmd := AddCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	} else {
		t.Logf("output: %s", stdout.String())
	}

	m, err := cmd.Fetch("top")

	if err != nil {
		t.Fatalf("can't fetch: %s", err)
	}

	exp := map[string]string{"a": "b"}

	if !reflect.DeepEqual(m, exp) {
		t.Errorf("invalid values: %#v", m)
	}
}
