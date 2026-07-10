package app

import "testing"

func TestPairRotateAndClobberGuard(t *testing.T) {
	seshareDir = t.TempDir()

	if err := cmdPair([]string{"alice"}); err != nil {
		t.Fatalf("first pair: %v", err)
	}
	first, _ := getContact("alice")
	if first == "" {
		t.Fatal("pair did not store a code")
	}

	// re-pairing without --rotate must refuse (no silent clobber)
	if err := cmdPair([]string{"alice"}); err == nil {
		t.Error("re-pair without --rotate should error")
	}
	if code, _ := getContact("alice"); code != first {
		t.Error("refused pair must not change the code")
	}

	// --rotate replaces the code
	if err := cmdPair([]string{"alice", "--rotate"}); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if code, _ := getContact("alice"); code == first {
		t.Error("--rotate did not change the code")
	}
}
