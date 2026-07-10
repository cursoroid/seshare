package session

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry is one session transcript found on disk. Cheap to build (stat only).
type Entry struct {
	Path    string
	Project string // munged project dir name
	ModTime time.Time
}

// ListAll returns every session across all projects, newest first.
func ListAll() ([]Entry, error) {
	projects, err := os.ReadDir(claudeProjectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Entry
	for _, p := range projects {
		if !p.IsDir() {
			continue
		}
		dir := filepath.Join(claudeProjectsDir, p.Name())
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || filepath.Ext(f.Name()) != ".jsonl" {
				continue
			}
			fi, err := f.Info()
			if err != nil {
				continue
			}
			out = append(out, Entry{Path: filepath.Join(dir, f.Name()), Project: p.Name(), ModTime: fi.ModTime()})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime.After(out[j].ModTime) })
	return out, nil
}

// Meta is the preview of a single session, parsed from its transcript.
type Meta struct {
	Cwd         string
	Title       string
	FirstPrompt string
	LastPrompt  string
	MsgCount    int
}

// Preview parses one transcript for display fields. One pass, whole file.
func Preview(path string) (Meta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Meta{}, err
	}
	var m Meta
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var l struct {
			Type       string `json:"type"`
			Cwd        string `json:"cwd"`
			AiTitle    string `json:"aiTitle"`
			LastPrompt string `json:"lastPrompt"`
			Message    *struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(line, &l) != nil {
			continue
		}
		if m.Cwd == "" && l.Cwd != "" {
			m.Cwd = l.Cwd
		}
		if l.AiTitle != "" {
			m.Title = l.AiTitle
		}
		if l.LastPrompt != "" {
			m.LastPrompt = l.LastPrompt
		}
		if l.Type == "user" || l.Type == "assistant" {
			m.MsgCount++
		}
		if l.Type == "user" && m.FirstPrompt == "" && l.Message != nil {
			m.FirstPrompt = messageText(l.Message.Content)
		}
	}
	return m, nil
}

// messageText pulls display text from a message content field, which is either
// a plain string or an array of content blocks with {"type":"text","text":…}.
func messageText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) == nil {
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				return b.Text
			}
		}
	}
	return ""
}
