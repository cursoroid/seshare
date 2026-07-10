package app

import (
	"bytes"
	"encoding/json"
)

// rewriteTranscript rewrites each JSONL line so a recipient can resume the
// session: cwd -> newCwd and sessionId/session_id -> newID, only on lines that
// already carry those fields. Other fields and blank lines are preserved.
//
// ponytail: full map round-trip reorders keys and floats large ints; fine for
// JSONL (order-independent) and Claude's small numeric fields. Switch to a
// targeted field edit if a future field needs exact byte fidelity.
func rewriteTranscript(data []byte, newCwd, newID string) []byte {
	lines := bytes.Split(data, []byte("\n"))
	for i, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var m map[string]any
		if json.Unmarshal(line, &m) != nil {
			continue // leave anything unparseable untouched
		}
		if _, ok := m["cwd"]; ok {
			m["cwd"] = newCwd
		}
		if _, ok := m["sessionId"]; ok {
			m["sessionId"] = newID
		}
		if _, ok := m["session_id"]; ok {
			m["session_id"] = newID
		}
		if b, err := json.Marshal(m); err == nil {
			lines[i] = b
		}
	}
	return bytes.Join(lines, []byte("\n"))
}

// stripSnapshots drops file-history-snapshot lines and any line marked
// isSnapshotUpdate — they reference the sender's local files/checkpoints and
// can trip up resume. Fallback used by `recv --strip-snapshots`.
func stripSnapshots(data []byte) []byte {
	var kept [][]byte
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			kept = append(kept, line)
			continue
		}
		var m struct {
			Type             string `json:"type"`
			IsSnapshotUpdate bool   `json:"isSnapshotUpdate"`
		}
		if json.Unmarshal(line, &m) == nil && (m.Type == "file-history-snapshot" || m.IsSnapshotUpdate) {
			continue
		}
		kept = append(kept, line)
	}
	return bytes.Join(kept, []byte("\n"))
}
