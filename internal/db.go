package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/boltdb/bolt"
)

// DB exists in case we want to mock the DB later.
type DB interface {
	SetKey(realm, key string, data Sealed) error
	GetKey(realm, key string) (Sealed, error)
	DropKey(realm, key string) error
	ListKeys(realm string) ([]string, error)
	ListRealms() ([]string, error)
	Purge(realm string) error
	GetAllKeys(realm string) (Stored, error)
	SetKeys(realm string, keys Stored) error
	Close() error
}

var ErrNotFound = errors.New("not found")

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB(fpath string) (*BoltDB, error) {
	if err := ensureDir(path.Dir(fpath)); err != nil {
		return nil, err
	}

	db, err := bolt.Open(fpath, 0600, nil)

	if err != nil {
		return nil, err
	}

	return &BoltDB{db}, nil
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}

func (b *BoltDB) GetKey(realm, key string) (s Sealed, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(realm))

		if bk == nil {
			return fmt.Errorf("realm %s: %w", realm, ErrNotFound)
		}

		v := bk.Get([]byte(key))

		if v == nil {
			return fmt.Errorf("%s/%s: %w", realm, key, ErrNotFound)
		}

		return json.Unmarshal(v, &s)
	})

	return
}

func (b *BoltDB) SetKey(realm, key string, s Sealed) error {
	v, err := json.Marshal(s)

	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(realm))

		if err != nil {
			return err
		}

		return bk.Put([]byte(key), v)
	})
}

func (b *BoltDB) DropKey(realm, key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(realm))

		if bk == nil {
			return fmt.Errorf("realm %s: %w", realm, ErrNotFound)
		}

		return bk.Delete([]byte(key))
	})
}

func (b *BoltDB) ListKeys(realm string) (s []string, err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(realm))

		if bk == nil {
			return fmt.Errorf("realm %s: %w", realm, ErrNotFound)
		}

		s = make([]string, 0)

		// must copy any byte slice to avoid invalid
		// memory references later; k & v are volatile

		return bk.ForEach(func(k, v []byte) error {
			c := make([]byte, len(k))
			copy(c, k)
			s = append(s, string(c))
			return nil
		})
	})

	return
}

func (b *BoltDB) ListRealms() (s []string, err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		s = make([]string, 0)

		// must copy any byte slice to avoid invalid
		// memory references later; k & v are volatile

		return tx.ForEach(func(k []byte, v *bolt.Bucket) error {
			c := make([]byte, len(k))
			copy(c, k)
			s = append(s, string(c))
			return nil
		})
	})

	return
}

func (b *BoltDB) Purge(realm string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(realm))
	})
}

func (b *BoltDB) GetAllKeys(realm string) (s Stored, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(realm))

		if bk == nil {
			return fmt.Errorf("realm %s: %w", realm, ErrNotFound)
		}

		s = make(Stored)

		// must copy any byte slice to avoid invalid
		// memory references later; k & v are volatile

		return bk.ForEach(func(k, v []byte) error {
			c := make([]byte, len(k))
			copy(c, k)

			var sd Sealed

			if err := json.Unmarshal(v, &sd); err != nil {
				return err
			}

			s[string(c)] = sd
			return nil
		})
	})

	return
}

func (b *BoltDB) SetKeys(realm string, s Stored) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(realm))

		if err != nil {
			return err
		}

		for k, sd := range s {
			v, err := json.Marshal(sd)
			if err != nil {
				return err
			}
			if err = bk.Put([]byte(k), v); err != nil {
				return err
			}
		}

		return nil
	})
}

func ensureDir(path string) error {
	fi, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, os.ModeDir|0700); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s exists, but is not a directory", path)
	}

	return nil
}
