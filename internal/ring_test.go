package internal

import (
	"encoding/hex"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestMockKeyring(t *testing.T) {
	// once we've set the mock we're done in this
	// test process -- even for other unit tests

	keyring.MockInit()

	k := &Keychain{
		service: "envy-test",
		user:    "test-user",
		keyer:   &mockGenerator{},
	}

	s, err := k.GetSecret()

	if err != nil {
		t.Fatal("first-get", err)
	} else if hex.EncodeToString(s) != "fe4fcdff6e7cd635ee52791f963dd15d52baa0ccdf4d1feb06559e5801377b00" {
		t.Errorf("first-get wrong: %s", hex.EncodeToString(s))
	}

	s2, err := k.GetSecret()

	if err != nil {
		t.Fatal("second-get", err)
	} else if string(s) != string(s2) {
		t.Errorf("second-get wrong: %s", hex.EncodeToString(s2))
	}
}
