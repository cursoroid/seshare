package app

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRewriteTranscript(t *testing.T) {
	in := strings.Join([]string{
		`{"type":"user","cwd":"/old/path","sessionId":"OLD-ID","uuid":"u1","message":{"role":"user"}}`,
		`{"type":"assistant","cwd":"/old/path","sessionId":"OLD-ID","session_id":"OLD-ID","uuid":"u2"}`,
		`{"type":"file-history-snapshot","messageId":"m1"}`, // no cwd/sessionId
		``, // blank line preserved
	}, "\n")

	out := rewriteTranscript([]byte(in), "/new/pwd", "NEW-ID")

	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("line not valid json after rewrite: %v\n%s", err, line)
		}
		if cwd, ok := m["cwd"]; ok && cwd != "/new/pwd" {
			t.Errorf("cwd = %v, want /new/pwd", cwd)
		}
		if sid, ok := m["sessionId"]; ok && sid != "NEW-ID" {
			t.Errorf("sessionId = %v, want NEW-ID", sid)
		}
		if sid, ok := m["session_id"]; ok && sid != "NEW-ID" {
			t.Errorf("session_id = %v, want NEW-ID", sid)
		}
	}

	// the snapshot line had no cwd; rewrite must not invent one
	if strings.Contains(string(out), `"file-history-snapshot"`) &&
		strings.Contains(strings.Split(string(out), "\n")[2], `"cwd"`) {
		t.Error("cwd should not be added to a line that lacked it")
	}
}
