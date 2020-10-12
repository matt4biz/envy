package internal

import (
	"reflect"
	"testing"
)

type extractTest struct {
	args     []string
	leftover []string
	expected map[string]string
	name     string
	realm    string
	err      error
	index    int
}

func (tc extractTest) run(t *testing.T) {
	e, err := NewExtractor(tc.args)

	if err != nil {
		if err != tc.err {
			t.Fatalf("invalid err: %s", err)
		}

		return
	}

	if e == nil {
		t.Fatalf("no extractor, no err")
	}

	if e.realm != tc.realm {
		t.Errorf("invalid realm: %s", e.realm)
	}

	if e.index != tc.index {
		t.Errorf("invalid index: %d", e.index)
	}

	if tc.leftover != nil {
		if len(e.args) != len(tc.leftover) || !reflect.DeepEqual(tc.leftover, e.args) {
			t.Errorf("invalid leftover: %#v; wanted %#v", e.args, tc.leftover)
		}

	} else if len(e.args) != 0 {
		t.Errorf("invalid leftover: %#v; wanted none", e.args)
	}

	if tc.expected != nil {
		if len(e.values) != len(tc.expected) || !reflect.DeepEqual(tc.expected, e.values) {
			t.Errorf("invalid values: %#v; wanted %#v", e.values, tc.expected)
		}
	} else if len(e.values) != 0 {
		t.Errorf("invalid values: %#v; wanted none", e.values)
	}
}

func TestExtractor(t *testing.T) {
	table := []extractTest{
		{
			name: "none",
			err:  ErrNoArguments,
		},
		{
			name: "bad realm",
			args: []string{","},
			err:  ErrBadRealm,
		},
		{
			name:  "only realm",
			args:  []string{"top"},
			realm: "top",
		},
		{
			name:     "one pair",
			args:     []string{"top", "a=b"},
			realm:    "top",
			index:    1,
			expected: map[string]string{"a": "b"},
		},
		{
			name:     "one pair plus",
			args:     []string{"top", "a=b", "c", "d"},
			realm:    "top",
			index:    1,
			expected: map[string]string{"a": "b"},
			leftover: []string{"c", "d"},
		},
		{
			name:  "two pair",
			args:  []string{"top", "a=b", "c={a:b}"},
			realm: "top",
			index: 2,
			expected: map[string]string{
				"a": "b",
				"c": "{a:b}",
			},
		},
		{
			name:  "quoted pair",
			args:  []string{"top", `c={"a":"b"}`},
			realm: "top",
			index: 1,
			expected: map[string]string{
				"c": `{"a":"b"}`,
			},
		},
		{
			name:     "pair with space",
			args:     []string{"top", "a =b"},
			realm:    "top",
			index:    0,
			leftover: []string{"a =b"},
		},
		{
			name:     "pair with space after",
			args:     []string{"top", "a= b"},
			realm:    "top",
			index:    1,
			expected: map[string]string{"a": "b"},
		},
		{
			name:     "pair with 2 spaces",
			args:     []string{"top", "a = b"},
			realm:    "top",
			index:    0,
			leftover: []string{"a = b"},
		},
		{
			name:  "two pair with space after",
			args:  []string{"top", "a=b", `c= {"a":"b"}`},
			realm: "top",
			index: 2,
			expected: map[string]string{
				"a": "b",
				"c": `{"a":"b"}`,
			},
		},
		{
			name:  "two pair with space after plus",
			args:  []string{"top", "a=b", `c= {"a":"b"}`, "d = e"},
			realm: "top",
			index: 2,
			expected: map[string]string{
				"a": "b",
				"c": `{"a":"b"}`,
			},
			leftover: []string{"d = e"},
		},
	}

	for _, tc := range table {
		t.Run(tc.name, tc.run)
	}
}
