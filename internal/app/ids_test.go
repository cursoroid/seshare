package app

import (
	"regexp"
	"testing"
)

func TestNewUUID(t *testing.T) {
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	a, b := newUUID(), newUUID()
	if !re.MatchString(a) {
		t.Errorf("newUUID() = %q, not a v4 uuid", a)
	}
	if a == b {
		t.Error("newUUID() returned identical values")
	}
}

func TestNewCode(t *testing.T) {
	c := newCode()
	if len(c) < 12 {
		t.Errorf("newCode() = %q, too short to be a secret", c)
	}
	if c == newCode() {
		t.Error("newCode() not random")
	}
}
