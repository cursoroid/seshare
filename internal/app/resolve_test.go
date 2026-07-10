package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveSession(t *testing.T) {
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
	// make `older` genuinely older
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(older, old, old); err != nil {
		t.Fatal(err)
	}

	// no id -> newest by mtime
	got, err := resolveSession(cwd, "")
	if err != nil || got != newer {
		t.Errorf("resolveSession(newest) = %q,%v want %q", got, err, newer)
	}

	// explicit id -> that file
	got, err = resolveSession(cwd, "older")
	if err != nil || got != older {
		t.Errorf("resolveSession(id) = %q,%v want %q", got, err, older)
	}

	// missing id -> error
	if _, err := resolveSession(cwd, "nope"); err == nil {
		t.Error("resolveSession(missing id) should error")
	}
}
