package app

import (
	"os"
	"testing"
)

func TestContactsRoundTrip(t *testing.T) {
	seshareDir = t.TempDir()

	if _, err := getContact("nobody"); err == nil {
		t.Error("getContact on missing contact should error")
	}

	if err := addContact("alice", "code-abc"); err != nil {
		t.Fatalf("addContact: %v", err)
	}
	if err := addContact("bob", "code-xyz"); err != nil {
		t.Fatalf("addContact: %v", err)
	}

	if code, err := getContact("alice"); err != nil || code != "code-abc" {
		t.Errorf("getContact(alice) = %q,%v want code-abc,nil", code, err)
	}

	names, err := listContactNames()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("listContactNames = %v, want 2 entries", names)
	}

	fi, err := os.Stat(contactsFile())
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("contacts file perm = %o, want 600", fi.Mode().Perm())
	}
}
