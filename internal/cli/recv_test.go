package cli

import (
	"testing"

	"github.com/cursoroid/seshare/internal/contacts"
)

func TestRecvCode(t *testing.T) {
	contacts.SeshareDir = t.TempDir()
	if err := contacts.Add("bob", "BOBCODE"); err != nil {
		t.Fatal(err)
	}

	// @name -> contact code
	if code, via, err := recvCode("@bob"); err != nil || code != "BOBCODE" || !via {
		t.Errorf(`recvCode("@bob") = %q,%v,%v want BOBCODE,true,nil`, code, via, err)
	}
	// bare name that IS a contact -> contact code (the footgun fix)
	if code, via, err := recvCode("bob"); err != nil || code != "BOBCODE" || !via {
		t.Errorf(`recvCode("bob") = %q,%v,%v want BOBCODE,true,nil`, code, via, err)
	}
	// bare word that is NOT a contact -> literal one-time code
	if code, via, err := recvCode("a1b2c3d4"); err != nil || code != "a1b2c3d4" || via {
		t.Errorf(`recvCode("a1b2c3d4") = %q,%v,%v want a1b2c3d4,false,nil`, code, via, err)
	}
	// @unknown -> error (explicit contact reference that doesn't exist)
	if _, _, err := recvCode("@nobody"); err == nil {
		t.Error(`recvCode("@nobody") should error`)
	}
}
