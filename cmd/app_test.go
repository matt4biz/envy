package main

import (
	"bytes"
	"strings"
	"testing"
)

// All these tests are safe because they should exit
// the app before a database may be created

func TestAppNoCommand(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	o := runApp(nil, "666", nil, stdout, stderr)

	t.Log(stderr.String())

	if o != -1 {
		t.Errorf("invalid code: %d", o)
	}
}

func TestAppInvalidCommand(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	o := runApp([]string{"empty"}, "666", nil, stdout, stderr)

	t.Log(stderr.String())

	if o != -1 {
		t.Errorf("invalid code: %d", o)
	}
}

func TestAppUsage(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	o := runApp([]string{"-h"}, "666", nil, stdout, stderr)

	if o != 0 {
		t.Errorf("invalid code: %d", o)
	}

	lines := strings.Split(stderr.String(), "\n")

	if len(lines) < 1 || !strings.HasPrefix(lines[0], "envy: a tool") {
		t.Errorf("invalid stderr: %s", stderr.String())
	}
}

func TestAppVersion(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	o := runApp([]string{"version"}, "666", nil, stdout, stderr)

	if o != 0 {
		t.Errorf("invalid code: %d", o)
	}

	lines := strings.Split(stdout.String(), "\n")

	if len(lines) < 1 || lines[0] != "666" {
		t.Errorf("invalid stdout: %s", stdout.String())
	}
}
