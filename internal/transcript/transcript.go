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
		Usage   Usage           `json:"usage"`
	} `json:"message"`
}

// Usage mirrors message.usage from a Claude Code assistant entry.
type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
}

// Context returns the current context window fill (input + cache_read +
// cache_creation tokens) from the most recent assistant message. Returns
// zero if the transcript has no usage data.
func Context(path string) int {
	u, ok := lastUsage(path)
	if !ok {
		return 0
	}
	return u.InputTokens + u.CacheReadInputTokens + u.CacheCreationInputTokens
}

// LastUsage returns the most recent assistant message's usage data.
func LastUsage(path string) (Usage, bool) {
	return lastUsage(path)
}

// TotalTokens returns the cumulative (sum across all assistant messages)
// input and output token counts for the transcript.
func TotalTokens(path string) (input, output int) {
	if path == "" {
		return 0, 0
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var e entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		if e.Type != "assistant" {
			continue
		}
		input += e.Message.Usage.InputTokens
		output += e.Message.Usage.OutputTokens
	}
	return input, output
}

func lastUsage(path string) (Usage, bool) {
	if path == "" {
		return Usage{}, false
	}
	f, err := os.Open(path)
	if err != nil {
		return Usage{}, false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	var last Usage
	found := false
	for scanner.Scan() {
		var e entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		if e.Type != "assistant" {
			continue
		}
		if e.Message.Usage.InputTokens == 0 &&
			e.Message.Usage.CacheReadInputTokens == 0 &&
			e.Message.Usage.CacheCreationInputTokens == 0 {
			continue
		}
		last = e.Message.Usage
		found = true
	}
	return last, found
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
