package transcript

import (
	"strings"
	"testing"
)

func TestStripSnapshots(t *testing.T) {
	in := strings.Join([]string{
		`{"type":"user","uuid":"u1"}`,
		`{"type":"file-history-snapshot","messageId":"m1"}`,
		`{"type":"assistant","uuid":"u2","isSnapshotUpdate":true}`,
		`{"type":"assistant","uuid":"u3"}`,
	}, "\n")

	out := string(StripSnapshots([]byte(in)))

	if strings.Contains(out, "file-history-snapshot") {
		t.Error("file-history-snapshot line not stripped")
	}
	if strings.Contains(out, "isSnapshotUpdate") {
		t.Error("isSnapshotUpdate line not stripped")
	}
	for _, keep := range []string{`"u1"`, `"u3"`} {
		if !strings.Contains(out, keep) {
			t.Errorf("non-snapshot line %s was dropped", keep)
		}
	}
}
