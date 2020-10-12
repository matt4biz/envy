package envy

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/matt4biz/envy/internal"
	"github.com/zalando/go-keyring"
)

func TestEnvy(t *testing.T) {
	// once we've set the mock we're done in this
	// test process -- even for other unit tests

	keyring.MockInit()

	dname, err := ioutil.TempDir("", "envy")

	if err != nil {
		t.Fatal("tempdir", err)
	}

	t.Log(dname)

	defer os.RemoveAll(dname)

	e, err := New(dname)

	if err != nil {
		t.Fatal("new", err)
	}

	t.Log(e.CurrentUser())

	m := map[string]string{
		"name": "matt",
		"age":  "58",
	}

	if err := e.Add("data", m); err != nil {
		t.Fatal("add", err)
	}

	m2, err := e.Fetch("data")

	if err != nil {
		t.Fatal("fetch", err)
	} else if !reflect.DeepEqual(m, m2) {
		t.Errorf("invalid data: %#v", m2)
	}

	if j, err := e.FetchAsJSON("data"); err != nil {
		t.Errorf("json: %s", err)
	} else {
		t.Logf("json: %s", j)
	}

	if err := e.Set("data", "key", "0x77e3"); err != nil {
		t.Errorf("set: %s", err)
	}

	m["key"] = "0x77e3"

	m2, err = e.Fetch("data")

	if err != nil {
		t.Fatal("fetch", err)
	} else if !reflect.DeepEqual(m, m2) {
		t.Errorf("invalid data: %#v", m2)
	}

	t.Log(m2)

	b := new(bytes.Buffer)

	if err := e.List(b, "data", "", true); err != nil {
		t.Errorf("list: %s", err)
	} else {
		t.Log("\n", b.String())
	}

	if s, err := e.Get("data", "key"); err != nil || s != m["key"] {
		t.Errorf("get: %s", err)
	}

	if err := e.Drop("data", "key"); err != nil {
		t.Errorf("drop: %s", err)
	}

	if _, err := e.Get("data", "key"); !errors.Is(err, internal.ErrNotFound) {
		t.Errorf("after drop: %s", err)
	}

	if err := e.Purge("data"); err != nil {
		t.Fatal("purge", err)
	}

	if _, err = e.Fetch("data"); !errors.Is(err, internal.ErrNotFound) {
		t.Errorf("after purge: %s", err)
	}
}
