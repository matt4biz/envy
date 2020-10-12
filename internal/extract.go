package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	assignRe = regexp.MustCompile(`^(?P<key>[\w]+)=(?P<value>.+)$`)
	realmRe  = regexp.MustCompile(`^[-:/\w]+$`)

	ErrNoArguments = errors.New("not enough arguments")
	ErrBadRealm    = errors.New("invalid realm: non-word characters")
)

type Extractor struct {
	args   []string
	values map[string]string
	realm  string
	index  int
}

func NewExtractor(args []string) (*Extractor, error) {
	e := Extractor{args: args, values: make(map[string]string)}
	err := e.prepare()

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *Extractor) Args() []string {
	return e.args
}

func (e *Extractor) Realm() string {
	return e.realm
}

func (e *Extractor) Values() map[string]string {
	return e.values
}

func (e *Extractor) prepare() error {
	if len(e.args) == 0 {
		return ErrNoArguments
	}

	e.realm = e.args[0]
	e.args = e.args[1:]

	if !realmRe.MatchString(e.realm) {
		return ErrBadRealm
	}

	for _, a := range e.args {
		// FindAllStringSubmatch is going to return [][]string in some
		// form like [["a=b" "a" "b"]], so the first-level slice needs
		// to have length 1, and the next length 3; we want m[0][1:2]

		m := assignRe.FindAllStringSubmatch(a, -1)

		if len(m) == 0 {
			break
		}

		if len(m) == 1 && len(m[0]) == 3 {
			k := strings.TrimSpace(m[0][1])
			v := strings.TrimSpace(m[0][2])

			e.values[k] = v
			e.index++
			continue
		}

		return fmt.Errorf("invalid pair %s = %v", a, m)
	}

	e.args = e.args[e.index:]
	return nil
}
