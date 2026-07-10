package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolve(t *testing.T) {
	claudeProjectsDir = t.TempDir()
	cwd := "/some/proj"
	dir := filepath.Join(claudeProjectsDir, mungedCwd(cwd))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	older := filepath.Join(dir, "older.jsonl")
	newer := filepath.Join(dir, "newer.jsonl")
	for _, f := range []string{older, newer} {
		if err := os.WriteFile(f, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(older, old, old); err != nil {
		t.Fatal(err)
	}

	if got, err := Resolve(cwd, ""); err != nil || got != newer {
		t.Errorf("Resolve(newest) = %q,%v want %q", got, err, newer)
	}
	if got, err := Resolve(cwd, "older"); err != nil || got != older {
		t.Errorf("Resolve(id) = %q,%v want %q", got, err, older)
	}
	if _, err := Resolve(cwd, "nope"); err == nil {
		t.Error("Resolve(missing id) should error")
	}
}
