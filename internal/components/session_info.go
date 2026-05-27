package components

import (
	"strings"

	"github.com/andhikapraa/goccline/internal/transcript"
)

const (
	defaultIDLength    = 8
	defaultTitleLength = 40
)

func init() {
	Register("session_info", renderSessionInfo)
}

// renderSessionInfo emits "🔗 abc12345 first user message..." — the short
// session ID followed by the first user prompt read from the transcript,
// truncated to defaultTitleLength chars.
func renderSessionInfo(ctx *Context) string {
	id := ctx.Input.SessionID
	if id == "" {
		return ""
	}
	if len(id) > defaultIDLength {
		id = id[:defaultIDLength]
	}

	title := strings.TrimSpace(transcript.FirstUserMessage(ctx.Input.TranscriptPath))
	if title != "" {
		title = collapseWhitespace(title)
		if len(title) > defaultTitleLength {
			title = title[:defaultTitleLength-3] + "..."
		}
	}

	out := "🔗 " + id
	if title != "" {
		out += " " + title
	}
	return out
}

// collapseWhitespace flattens any run of whitespace (including newlines) to
// a single space so multiline first prompts don't break the statusline.
func collapseWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		b.WriteRune(r)
		prevSpace = false
	}
	return b.String()
}
