package cli

import (
	"testing"

	"github.com/cursoroid/seshare/internal/contacts"
)

func TestParseSend(t *testing.T) {
	contacts.SeshareDir = t.TempDir()
	if err := contacts.Add("bob", "CODE"); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		args     []string
		wantID   string
		wantName string
		wantYes  bool
	}{
		{[]string{"@bob"}, "", "bob", false},               // explicit contact
		{[]string{"bob"}, "", "bob", false},                // bare name resolves to contact
		{[]string{"6a9ad174-xxxx"}, "6a9ad174-xxxx", "", false}, // non-contact -> session id
		{[]string{"6a9ad174-xxxx", "bob"}, "6a9ad174-xxxx", "bob", false},
		{[]string{"bob", "--yes"}, "", "bob", true},
	}
	for _, c := range cases {
		id, name, yes, err := parseSend(c.args)
		if err != nil {
			t.Errorf("parseSend(%v) err = %v", c.args, err)
			continue
		}
		if id != c.wantID || name != c.wantName || yes != c.wantYes {
			t.Errorf("parseSend(%v) = (%q,%q,%v) want (%q,%q,%v)", c.args, id, name, yes, c.wantID, c.wantName, c.wantYes)
		}
	}

	if _, _, _, err := parseSend([]string{"--bogus"}); err == nil {
		t.Error("parseSend should reject unknown flags")
	}
}
