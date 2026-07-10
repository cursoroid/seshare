package ids

import (
	"regexp"
	"testing"
)

func TestNewUUID(t *testing.T) {
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	a, b := NewUUID(), NewUUID()
	if !re.MatchString(a) {
		t.Errorf("NewUUID() = %q, not a v4 uuid", a)
	}
	if a == b {
		t.Error("NewUUID() returned identical values")
	}
}

func TestNewCode(t *testing.T) {
	c := NewCode()
	if len(c) < 12 {
		t.Errorf("NewCode() = %q, too short to be a secret", c)
	}
	if c == NewCode() {
		t.Error("NewCode() not random")
	}
}
