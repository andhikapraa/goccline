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

// renderSessionInfo emits "🔗 abc12345" by default. When
// config.session_info.show_first_message is true, also appends the first
// user prompt read from the transcript file — opt-in because transcript I/O
// can add ~130ms per render on large transcripts.
func renderSessionInfo(ctx *Context) string {
	id := ctx.Input.SessionID
	if id == "" {
		return ""
	}
	idLen := ctx.Config.SessionInfo.IDLength
	if idLen <= 0 {
		idLen = defaultIDLength
	}
	if len(id) > idLen {
		id = id[:idLen]
	}

	t := ctx.Theme
	out := "🔗 " + t.Session + id + t.Reset

	if !ctx.Config.SessionInfo.ShowFirstMessage {
		return out
	}

	title := strings.TrimSpace(transcript.FirstUserMessage(ctx.Input.TranscriptPath))
	if title == "" {
		return out
	}
	title = collapseWhitespace(title)
	maxLen := ctx.Config.SessionInfo.TitleLength
	if maxLen <= 0 {
		maxLen = defaultTitleLength
	}
	if len(title) > maxLen {
		title = title[:maxLen-3] + "..."
	}
	return out + " " + t.SessionText + title + t.Reset
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
