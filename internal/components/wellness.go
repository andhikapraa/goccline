package components

import (
	"fmt"
	"time"

	"github.com/andhikapraa/goccline/internal/transcript"
)

func init() {
	Register("wellness", renderWellness)
}

// renderWellness shows how long the user has been in their current coding
// burst, defined as time since the most recent user message. The threshold
// is 45 minutes — beyond that, the icon flips to suggest a break.
func renderWellness(ctx *Context) string {
	last := transcript.LastUserMessageTime(ctx.Input.TranscriptPath)
	if last.IsZero() {
		return ""
	}
	elapsed := time.Since(last)
	if elapsed < 0 {
		return ""
	}
	mins := int(elapsed.Minutes())
	if mins > 24*60 {
		// Stale session — don't claim someone's been coding for days.
		return ""
	}

	t := ctx.Theme
	icon := "☕"
	color := t.WellnessOk
	if mins >= 45 {
		icon = "🚶"
		color = t.WellnessNudge
	}
	return fmt.Sprintf("%s %sCoding %dm/45m%s", icon, color, mins, t.Reset)
}
