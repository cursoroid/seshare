// Package session locates Claude Code session transcripts on disk.
package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LocalClaudeVersion returns the installed Claude Code version, or "" if it
// can't be determined. Best-effort, used only for a mismatch warning.
func LocalClaudeVersion() string {
	out, err := exec.Command("claude", "--version").Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

// claudeProjectsDir is the root of Claude Code's per-project transcripts.
// Overridable in tests.
var claudeProjectsDir = defaultClaudeProjectsDir()

func defaultClaudeProjectsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".claude", "projects")
	}
	return filepath.Join(home, ".claude", "projects")
}

// Dir returns the project transcript directory for a working directory.
func Dir(cwd string) string {
	return filepath.Join(claudeProjectsDir, mungedCwd(cwd))
}

// Resolve returns the transcript path for a session in cwd's project.
// Empty id -> newest *.jsonl by mtime; otherwise the file named <id>.jsonl.
func Resolve(cwd, id string) (string, error) {
	dir := Dir(cwd)
	if id != "" {
		p := filepath.Join(dir, id+".jsonl")
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf("session %q not found in %s", id, dir)
		}
		return p, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("no sessions for this directory (%s): %w", dir, err)
	}
	var newest string
	var newestMod int64 = -1
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".jsonl" {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		if m := fi.ModTime().UnixNano(); m > newestMod {
			newestMod, newest = m, filepath.Join(dir, e.Name())
		}
	}
	if newest == "" {
		return "", fmt.Errorf("no .jsonl sessions in %s", dir)
	}
	return newest, nil
}

// mungedCwd converts a working directory into Claude Code's project folder
// name: every char that is not [a-zA-Z0-9-] becomes '-', with no collapsing of
// consecutive dashes. Verified against Claude Code 2.1.x.
func mungedCwd(path string) string {
	b := []byte(path)
	for i, c := range b {
		ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-'
		if !ok {
			b[i] = '-'
		}
	}
	return string(b)
}
