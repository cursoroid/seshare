package cli

import (
	"testing"

	"github.com/cursoroid/seshare/internal/contacts"
)

func TestPairRotateAndClobberGuard(t *testing.T) {
	contacts.SeshareDir = t.TempDir()

	if err := cmdPair([]string{"alice"}); err != nil {
		t.Fatalf("first pair: %v", err)
	}
	first, _ := contacts.Get("alice")
	if first == "" {
		t.Fatal("pair did not store a code")
	}

	// re-pairing without --rotate must refuse (no silent clobber)
	if err := cmdPair([]string{"alice"}); err == nil {
		t.Error("re-pair without --rotate should error")
	}
	if code, _ := contacts.Get("alice"); code != first {
		t.Error("refused pair must not change the code")
	}

	// --rotate replaces the code
	if err := cmdPair([]string{"alice", "--rotate"}); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if code, _ := contacts.Get("alice"); code == first {
		t.Error("--rotate did not change the code")
	}
}
