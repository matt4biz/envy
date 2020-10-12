package envy

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"path"

	"github.com/matt4biz/envy/internal"
)

// Envy provides the interface to the local secure
// variable store.
type Envy struct {
	db     internal.DB
	sealer *internal.Sealer
}

// New returns a secure variable store whose DB
// lives in the provided directory (typically the
// user's "config" directory).
func New(dir string) (*Envy, error) {
	s, err := internal.NewDefaultSealer()

	if err != nil {
		return nil, err
	}

	return NewWithSealer(dir, s)
}

// NewWithSealer is really for UTs, so we can pass in
// a fake sealer that's deterministic.
func NewWithSealer(dir string, s *internal.Sealer) (*Envy, error) {
	db, err := internal.NewBoltDB(path.Join(dir, "/envy.db"))

	if err != nil {
		return nil, err
	}

	e := Envy{
		db:     db,
		sealer: s,
	}

	return &e, nil
}

// CurrentUser returns the user's login name.
func (e *Envy) CurrentUser() string {
	return e.sealer.GetUsername()
}

// Add writes a map of {variable, value} pairs to the secure
// store, possibly creating it and/or overwriting variables
// that are already there.
func (e *Envy) Add(realm string, vars map[string]string) error {
	m := make(internal.Stored)

	for k, v := range vars {
		ud := internal.Unsealed{Data: v}
		sd, err := e.sealer.Seal(ud)

		if err != nil {
			return fmt.Errorf("sealing %s/%s: %w", realm, k, err)
		}

		m[k] = sd
	}

	return e.db.SetKeys(realm, m)
}

func (e *Envy) fetchRaw(realm string) (internal.Loaded, error) {
	m, err := e.db.GetAllKeys(realm)

	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", realm, err)
	}

	result := make(internal.Loaded, len(m))

	for k, sd := range m {
		ud, err := e.sealer.Unseal(sd)

		if err != nil {
			return nil, fmt.Errorf("unsealing %s/%s: %w", realm, k, err)
		}

		result[k] = ud
	}

	return result, nil
}

// Fetch returns a map of {variable, value} pairs from the
// secure store for the given realm, if present.
func (e *Envy) Fetch(realm string) (map[string]string, error) {
	m, err := e.fetchRaw(realm)

	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", realm, err)
	}

	result := make(map[string]string, len(m))

	for k, ud := range m {
		result[k] = ud.Data
	}

	return result, nil
}

// FetchAsJSON returns the variables for a given realm as a
// JSON object. This is handy for other tools using Envy as a
// library, e.g., to store secure login credentials / tokens.
func (e *Envy) FetchAsJSON(realm string) (json.RawMessage, error) {
	m, err := e.Fetch(realm)

	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

// FetchAsVarList returns the variables in a realm as a list of
// key=value expressions that can be appended to a command's
// list of environment variables.
func (e *Envy) FetchAsVarList(realm string) ([]string, error) {
	m, err := e.Fetch(realm)

	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(m))

	for k, v := range m {
		result = append(result, k+"="+v)
	}

	return result, nil
}

// Set adds a single key and value to the realm in the
// secure store, possibly creating it and/or overwriting
// and existing key.
func (e *Envy) Set(realm, key, data string) error {
	ud := internal.Unsealed{Data: data}
	sd, err := e.sealer.Seal(ud)

	if err != nil {
		return fmt.Errorf("sealing %s/%s: %w", realm, key, err)
	}

	return e.db.SetKey(realm, key, sd)
}

// Get returns a single key's value from the realm, if it
// is present.
func (e *Envy) Get(realm, key string) (string, error) {
	sd, err := e.db.GetKey(realm, key)

	if err != nil {
		return "", fmt.Errorf("fetching %s/%s: %w", realm, key, err)
	}

	ud, err := e.sealer.Unseal(sd)

	if err != nil {
		return "", fmt.Errorf("unsealing %s/%s: %w", realm, key, err)
	}

	return ud.Data, nil
}

// Drop removes a single key from the realm's secure store.
func (e *Envy) Drop(realm, key string) error {
	return e.db.DropKey(realm, key)
}

// Purge removes an entire realm from the secure store.
// Use with caution.
func (e *Envy) Purge(realm string) error {
	return e.db.Purge(realm)
}

// Realms returns a list of the realms in the secure store.
func (e *Envy) Realms() ([]string, error) {
	return e.db.ListRealms()
}

// List writes to its destination a single realm's variables and
// their metadata, and optionally their values (use with caution).
func (e *Envy) List(w io.Writer, realm string, decrypt bool) error {
	m, err := e.fetchRaw(realm)

	if err != nil {
		return err
	}

	var maxWidth int
	var maxSize int

	for k, v := range m {
		if l := len(k); l > maxWidth {
			maxWidth = l
		}

		if l := v.Meta.Size; l > maxSize {
			maxSize = l
		}
	}

	maxSize = (int)(math.Log10(float64(maxSize)) + 1)

	for k, ud := range m {
		if decrypt {
			fmt.Fprintf(w, "%-*s   %s   %s\n", maxWidth, k, ud.Meta.ToString(maxSize), ud.Data)
		} else {
			fmt.Fprintf(w, "%-*s   %s\n", maxWidth, k, ud.Meta.ToString(maxSize))
		}
	}

	return nil
}

func (e *Envy) Close() {
	_ = e.db.Close()
}
