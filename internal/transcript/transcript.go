// Package transcript reads claude-code transcript JSONL files. We only need
// the first user message right now; the implementation streams the file so
// large transcripts don't pay a full-parse cost.
package transcript

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"time"
)

type entry struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Message   struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var commandTagRE = regexp.MustCompile(`<[^>]*>`)

// FirstUserMessage returns the first user message text from path, with any
// XML-style command tags stripped. Returns "" if path is empty, missing, or
// has no user messages.
func FirstUserMessage(path string) string {
	if path == "" {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Some transcript lines are large; raise the buffer ceiling.
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		var e entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		if e.Type != "user" {
			continue
		}
		text := extractText(e.Message.Content)
		if text == "" {
			continue
		}
		return commandTagRE.ReplaceAllString(text, "")
	}
	return ""
}

// LastUserMessageTime returns the timestamp of the most recent user message
// in the transcript, or zero if unavailable. Used by wellness to know when
// the current coding burst started.
func LastUserMessageTime(path string) time.Time {
	if path == "" {
		return time.Time{}
	}
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024)

	var last time.Time
	for scanner.Scan() {
		line := scanner.Bytes()
		var e entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		if e.Type != "user" || e.Timestamp == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339, e.Timestamp)
		if err != nil {
			continue
		}
		last = t
	}
	return last
}

// extractText handles both string and array forms of message.content.
func extractText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	// Try string form.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	// Try array of content blocks.
	var blocks []contentBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				return b.Text
			}
		}
	}
	return ""
}
