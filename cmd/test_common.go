package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/matt4biz/envy"
	"github.com/matt4biz/envy/internal"
)

func NewTestApp(t *testing.T, so, se *bytes.Buffer) *App {
	// once we've set the mock we're done in this
	// test process -- even for other unit tests

	keyring.MockInit()

	dname, err := ioutil.TempDir("", "envy")

	if err != nil {
		t.Fatal("tempdir", err)
	}

	defer os.RemoveAll(dname)

	e, err := envy.NewWithSealer(dname, internal.NewTestSealer())

	if err != nil {
		t.Fatal("new-sealer", err)
	}

	app := App{
		Envy:   e,
		args:   []string{},
		path:   dname,
		stdout: so,
		stderr: se,
	}

	return &app
}
