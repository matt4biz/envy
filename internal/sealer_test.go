package internal

import (
	"testing"
)

func TestSealer(t *testing.T) {
	s := NewTestSealer()
	ud := Unsealed{Data: "matt"}

	t.Log("unsealed", ud)

	sd, err := s.Seal(ud)

	if err != nil {
		t.Fatal("seal", err)
	}

	t.Log("sealed", sd)

	ud2, err := s.Unseal(sd)

	if err != nil {
		t.Fatal("unseal", err)
	} else if ud2.Data != ud.Data {
		t.Errorf("invalid data: %#v", ud2)
	}

	t.Log("unsealed", ud2)
}
