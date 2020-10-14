package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	dname, err := ioutil.TempDir("", "scratch")

	if err != nil {
		t.Fatal("tempdir", err)
	}

	defer os.RemoveAll(dname)

	fname := path.Join(dname, "test.json")

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)
	expData := map[string]string{"a": "12", "b": "21"}

	if err := app.Add("test", expData); err != nil {
		t.Fatal("setup", err)
	}

	app.args = []string{"test", fname}

	cmd := ReadCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	raw, err := ioutil.ReadFile(fname)

	if err != nil {
		t.Fatal("read", err)
	}

	t.Log(string(raw))

	var readData map[string]string

	if err = json.Unmarshal(raw, &readData); err != nil {
		t.Fatal("decode", err)
	}

	if !reflect.DeepEqual(readData, expData) {
		t.Errorf("invalid data: %#v", readData)
	}
}

func TestReadStdout(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)
	expData := map[string]string{"a": "12", "b": "21"}

	if err := app.Add("test", expData); err != nil {
		t.Fatal("setup", err)
	}

	app.args = []string{"test", "-"}

	cmd := ReadCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	var readData map[string]string

	if err := json.NewDecoder(stdout).Decode(&readData); err != nil {
		t.Fatal("decode", err)
	}

	if !reflect.DeepEqual(readData, expData) {
		t.Errorf("invalid data: %#v", readData)
	}
}

func TestReadUnquoted(t *testing.T) {
	type X struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	type wrap struct {
		One X `json:"one"`
		Two X `json:"two"`
	}

	expData := map[string]wrap{
		"a": wrap{
			One: X{"1", "2"},
			Two: X{"5", "6"},
		},
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	app := NewTestApp(t, stdout, stderr)
	input := map[string]string{"a": `{"one":{"a":"1","b":"2"},\n "two":{"a":"5","b":"6"}}`}

	if err := app.Add("test", input); err != nil {
		t.Fatal("setup", err)
	}

	app.args = []string{"-q", "test", "-"}

	cmd := ReadCommand{app}
	o := cmd.Run()

	if o != 0 {
		t.Errorf("errors: %s", stderr.String())
		t.Fatalf("invalid return: %d", o)
	}

	t.Log("output:", stdout.String())

	var readData map[string]wrap

	if err := json.NewDecoder(stdout).Decode(&readData); err != nil {
		t.Fatal("decode", err)
	}

	if !reflect.DeepEqual(readData, expData) {
		t.Errorf("invalid data: %+v (should be %+v)", readData, expData)
	}
}
