package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	dname, err := ioutil.TempDir("", "scratch")

	if err != nil {
		t.Fatal("tempdir", err)
	}

	defer os.RemoveAll(dname)

	file, err := ioutil.TempFile(dname, "test")

	if err != nil {
		t.Fatal("tempfile", err)
	}

	if _, err = file.WriteString(`{"x":"21", "y":"14"}`); err != nil {
		t.Fatal("write-temp", err)
	}

	file.Close()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)

	app.args = []string{"test", file.Name()}

	cmd := WriteCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	m, err := cmd.Fetch("test")

	if err != nil {
		t.Fatalf("can't fetch: %s", err)
	}

	exp := map[string]string{"x": "21", "y": "14"}

	if !reflect.DeepEqual(m, exp) {
		t.Errorf("invalid values: %#v", m)
	}
}

func TestReadOverwrite(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)
	data := map[string]string{"a": "xxx", "b": "yyy"}

	if err := app.Add("test", data); err != nil {
		t.Fatal("setup", err)
	}

	app.stdin = bytes.NewBufferString(`{"x":"21", "y":"14"}`)
	app.args = []string{"-clear", "test", "-"}

	cmd := WriteCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	m, err := cmd.Fetch("test")

	if err != nil {
		t.Fatalf("can't fetch: %s", err)
	}

	exp := map[string]string{"x": "21", "y": "14"}

	if !reflect.DeepEqual(m, exp) {
		t.Errorf("invalid values: %#v", m)
	}
}
