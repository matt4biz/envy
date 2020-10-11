package internal

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os/user"

	"github.com/zalando/go-keyring"
)

const (
	defaultService = "matt4biz-envy-secret-key"
)

type Ring interface {
	GetSecret() ([]byte, error)
	GetUsername() string
}

type Keychain struct {
	service string
	user    string
	keyer   KeyGenerator
}

func NewKeychain() (*Keychain, error) {
	user, err := user.Current()

	if err != nil {
		return nil, err
	}

	k := Keychain{
		service: defaultService,
		user:    user.Username,
		keyer:   &realGenerator{},
	}

	return &k, nil
}

func (k *Keychain) GetSecret() ([]byte, error) {
	secret, err := keyring.Get(k.service, k.user)

	if err != nil {
		if !errors.Is(err, keyring.ErrNotFound) {
			return nil, err
		}

		key, err := k.keyer.MakeKey()

		if err != nil {
			return nil, err
		}

		secret = base64.StdEncoding.EncodeToString(key)

		if err := keyring.Set(k.service, k.user, secret); err != nil {
			return nil, err
		}
	}

	return base64.StdEncoding.DecodeString(secret)
}

func (k *Keychain) GetUsername() string {
	return k.user
}

type mockRing string

func (m mockRing) GetSecret() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(m))
}

func (m mockRing) GetUsername() string {
	return "test-user"
}

var testRing = mockRing("/k/N/2581jXuUnkflj3RXVK6oMzfTR/rBlWeWAE3ewA=")

type KeyGenerator interface {
	MakeKey() ([]byte, error)
}

type realGenerator struct{}

func (r realGenerator) MakeKey() ([]byte, error) {
	key := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return key, nil
}

type mockGenerator struct{}

func (m *mockGenerator) MakeKey() ([]byte, error) {
	return testRing.GetSecret()
}
