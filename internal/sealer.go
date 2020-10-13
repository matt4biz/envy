package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Sealed struct {
	Data string `json:"data"`
	Meta string `json:"metadata"`
}

type Unsealed struct {
	Data string
	Meta metadata
}

type Stored map[string]Sealed
type Loaded map[string]Unsealed

type metadata struct {
	Size     int    `json:"size"`      // size of the stored.Data before encryption
	Hash     string `json:"hash"`      // hash of the stored.Data  ""     ""
	Modified int64  `json:"timestamp"` // Unix time we make this data
}

func (md metadata) ToString(w int) string {
	if len(md.Hash) == 0 {
		return "unprepared"
	}

	return fmt.Sprintf("%s  %*d  %s", time.Unix(md.Modified, 0).Format(time.RFC3339), w, md.Size, md.Hash[:7])
}

type Sealer struct {
	Ring

	key    []byte
	noncer Noncer
}

func NewDefaultSealer() (*Sealer, error) {
	r, err := NewKeychain()

	if err != nil {
		return nil, err
	}

	k, err := r.GetSecret()

	if err != nil {
		return nil, err
	}

	s := Sealer{r, k, &realNonce{}}

	return &s, nil
}

func NewSealer(r Ring, n Noncer) (*Sealer, error) {
	k, err := r.GetSecret()

	if err != nil {
		return nil, err
	}

	return &Sealer{r, k, n}, nil
}

func (u *Unsealed) prep() ([]byte, []byte, error) {
	b, err := json.Marshal(u.Data)

	if err != nil {
		return nil, nil, err
	}

	hash := md5.New()

	if _, err = hash.Write(b); err != nil {
		return nil, nil, err
	}

	tag := hash.Sum(nil)

	u.Meta.Size = len(u.Data)
	u.Meta.Modified = time.Now().Unix()
	u.Meta.Hash = hex.EncodeToString(tag)

	return b, tag, nil
}

func (s Sealer) Seal(ud Unsealed) (Sealed, error) {
	var sd Sealed

	pt, aad, err := ud.prep()

	if err != nil {
		return sd, err
	}

	ct, err := s.encrypt(pt, aad)

	if err != nil {
		return sd, err
	}

	md, err := json.Marshal(ud.Meta)

	if err != nil {
		return sd, err
	}

	sd.Data = ct
	sd.Meta = base64.StdEncoding.EncodeToString(md)

	return sd, nil
}

func (s Sealer) Unseal(sd Sealed) (Unsealed, error) {
	var ud Unsealed

	md, err := base64.StdEncoding.DecodeString(sd.Meta)

	if err != nil {
		return ud, err
	}

	if err = json.Unmarshal(md, &ud.Meta); err != nil {
		return ud, err
	}

	pt, err := s.decrypt(sd.Data, ud.Meta.Hash)

	if err != nil {
		return ud, err
	}

	err = json.Unmarshal(pt, &ud.Data)
	return ud, err
}

func (s Sealer) encrypt(pt, aad []byte) (string, error) {
	nonce, err := s.noncer.GetNonce()

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.key)

	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	ct := aesgcm.Seal(nonce, nonce, pt, aad)

	return base64.StdEncoding.EncodeToString(ct), nil
}

func (s Sealer) decrypt(data, tag string) ([]byte, error) {
	mixed, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, err
	}

	nonce := mixed[0:12]
	ct := mixed[12:]
	block, err := aes.NewCipher(s.key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	aad, err := hex.DecodeString(tag)

	if err != nil {
		return nil, err
	}

	pt, err := aesgcm.Open(nil, nonce, ct, aad)

	if err != nil {
		return nil, err
	}

	return pt, nil
}

func NewTestSealer() *Sealer {
	s, _ := NewSealer(testRing, testNonce)

	return s
}

type Noncer interface {
	GetNonce() ([]byte, error)
}

type realNonce struct{}

func (r realNonce) GetNonce() ([]byte, error) {
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

type mockNonce string

func (m mockNonce) GetNonce() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(m))
}

var testNonce = mockNonce("tRKG2M0EpzwiyQfc")
