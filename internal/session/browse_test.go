package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListAllAndPreview(t *testing.T) {
	claudeProjectsDir = t.TempDir()
	proj := filepath.Join(claudeProjectsDir, "-Users-me-repo")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}
	lines := strings.Join([]string{
		`{"type":"user","cwd":"/Users/me/repo","sessionId":"S","message":{"role":"user","content":"first thing i asked"}}`,
		`{"type":"assistant","cwd":"/Users/me/repo","message":{"role":"assistant","content":[{"type":"text","text":"ok"}]}}`,
		`{"type":"ai-title","aiTitle":"Fix the parser"}`,
		`{"type":"user","message":{"role":"user","content":"second question"}}`,
		`{"type":"last-prompt","lastPrompt":"second question"}`,
	}, "\n")
	sess := filepath.Join(proj, "abc.jsonl")
	if err := os.WriteFile(sess, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := ListAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Path != sess {
		t.Fatalf("ListAll = %+v, want the one session", entries)
	}

	m, err := Preview(sess)
	if err != nil {
		t.Fatal(err)
	}
	if m.Cwd != "/Users/me/repo" {
		t.Errorf("Cwd = %q", m.Cwd)
	}
	if m.Title != "Fix the parser" {
		t.Errorf("Title = %q", m.Title)
	}
	if m.FirstPrompt != "first thing i asked" {
		t.Errorf("FirstPrompt = %q", m.FirstPrompt)
	}
	if m.LastPrompt != "second question" {
		t.Errorf("LastPrompt = %q", m.LastPrompt)
	}
	if m.MsgCount != 3 { // 2 user + 1 assistant
		t.Errorf("MsgCount = %d, want 3", m.MsgCount)
	}
}
