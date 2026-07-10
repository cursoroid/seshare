package contacts

import (
	"os"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	SeshareDir = t.TempDir()

	if _, err := Get("nobody"); err == nil {
		t.Error("Get on missing contact should error")
	}
	if err := Add("alice", "code-abc"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := Add("bob", "code-xyz"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if code, err := Get("alice"); err != nil || code != "code-abc" {
		t.Errorf("Get(alice) = %q,%v want code-abc,nil", code, err)
	}
	names, err := ListNames()
	if err != nil {
		t.Fatalf("ListNames: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("ListNames = %v, want 2 entries", names)
	}
	fi, err := os.Stat(contactsFile())
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("contacts file perm = %o, want 600", fi.Mode().Perm())
	}
}
