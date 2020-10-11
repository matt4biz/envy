package internal

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestBoltDBOps(t *testing.T) {
	dname, err := ioutil.TempDir("", "envy")

	if err != nil {
		t.Fatal("tempdir", err)
	}

	t.Log(dname)

	defer os.RemoveAll(dname)

	db, err := NewBoltDB(path.Join(dname, "/enty"))

	if err != nil {
		t.Fatal("newdb", err)
	}

	s := Sealed{Data: "data", Meta: "metadata"}

	if err := db.SetKey("top", "key", s); err != nil {
		t.Fatal("set", err)
	}

	db.Close()

	db2, err := NewBoltDB(path.Join(dname, "/enty"))

	if err != nil {
		t.Fatal("newdb 2", err)
	}

	if s2, err := db2.GetKey("top", "key"); err != nil {
		t.Fatal("get", err)
	} else if !reflect.DeepEqual(s, s2) {
		t.Errorf("bad data: %#v", s2)
	}

	if _, err := db2.GetKey("pot", "key"); err != nil {
		t.Log(err)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("wrong error for bad realm: %s", err)
		}
	} else {
		t.Errorf("no err for bad key!")
	}

	if err := db2.SetKey("top", "key2", s); err != nil {
		t.Fatal("set", err)
	}

	l, err := db2.ListKeys("top")
	t.Log("list", l)
	if err != nil {
		t.Fatal("list", err)
	} else if len(l) < 2 || l[0] != "key" || l[1] != "key2" {
		t.Errorf("invalid list: %#v", l)
	}

	x, err := db2.GetAllKeys("top")
	t.Log(x)
	if err != nil {
		t.Fatal("get-all", err)
	} else if len(x) != 2 || x["key"].Data != "data" {
		t.Errorf("invalid key set %#v", x)
	}

	if err := db2.SetKey("top2", "key", s); err != nil {
		t.Fatal("set", err)
	}

	l2, err := db2.ListRealms()
	t.Log("realms", l2)
	if err != nil {
		t.Fatal("realms", err)
	} else if len(l2) < 2 || l2[0] != "top" || l2[1] != "top2" {
		t.Errorf("invalid realms: %#v", l2)
	}

	if err := db2.DropKey("top", "key"); err != nil {
		t.Fatal("drop", err)
	}

	if _, err := db2.GetKey("top", "key"); err != nil {
		t.Log(err)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("wrong error for bad key: %s", err)
		}
	} else {
		t.Errorf("no err for bad key!")
	}

	if err := db2.DropKey("pot", "key"); err != nil {
		t.Log(err)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("wrong error for bad realm: %s", err)
		}
	} else {
		t.Errorf("no err for bad key!")
	}

	if err := db2.Purge("top2"); err != nil {
		t.Fatal("purge", err)
	}

	if _, err := db2.GetKey("top2", "key"); err != nil {
		t.Log(err)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("wrong error for bad key: %s", err)
		}
	} else {
		t.Errorf("no err for bad key!")
	}

	if err := db2.SetKeys("top3", x); err != nil {
		t.Fatal("set-all", err)
	}

	x2, err := db2.GetAllKeys("top3")
	t.Log(x2)
	if err != nil {
		t.Fatal("get-all", err)
	} else if !reflect.DeepEqual(x, x2) {
		t.Errorf("invalid key set %#v", x2)
	}

	l3, err := db2.ListRealms()
	t.Log("realms", l3)
	if err != nil {
		t.Fatal("realms", err)
	} else if len(l3) < 2 || l3[0] != "top" || l3[1] != "top3" {
		t.Errorf("invalid realms: %#v", l3)
	}

	db2.Close()
}
