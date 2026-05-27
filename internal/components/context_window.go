package components

import (
	"fmt"
	"strings"

	"github.com/andhikapraa/goccline/internal/transcript"
)

func init() {
	Register("context_window", renderContextWindow)
	Register("context_alert", renderContextAlert)
	Register("total_tokens", renderTotalTokens)
	Register("token_usage", renderTotalTokens) // alias
}

// contextLimit infers the model's context window from the display name.
// Opus 4.7 ships in both 200K and 1M variants; the bash statusline encodes
// this by appending "(1M context)" to display_name. We honor that signal.
func contextLimit(ctx *Context) int {
	name := ctx.Input.Model.DisplayName
	if strings.Contains(name, "1M") {
		return 1_000_000
	}
	return 200_000
}

// renderContextWindow shows "🧠 45% (90K/200K)". The percentage colors
// shift toward warning/critical as the window fills up.
func renderContextWindow(ctx *Context) string {
	used := transcript.Context(ctx.Input.TranscriptPath)
	if used == 0 {
		return ""
	}
	limit := contextLimit(ctx)
	pct := used * 100 / limit
	t := ctx.Theme
	color := t.WellnessOk
	switch {
	case pct >= 90:
		color = t.WellnessNudge
	case pct >= 75:
		color = t.Time
	}
	return fmt.Sprintf("🧠 %s%d%% (%s/%s)%s",
		color, pct, humanK(used), humanK(limit), t.Reset)
}

// renderContextAlert shows "⚠️ >LIMIT" in red if the context window is
// exceeded. Hidden when under the limit.
func renderContextAlert(ctx *Context) string {
	used := transcript.Context(ctx.Input.TranscriptPath)
	limit := contextLimit(ctx)
	if used < limit {
		return ""
	}
	t := ctx.Theme
	return fmt.Sprintf("%s⚠️ >%s%s", t.WellnessNudge, humanK(limit), t.Reset)
}

// renderTotalTokens shows cumulative input/output for the session.
// Format: "15.2K in / 4.5K out".
func renderTotalTokens(ctx *Context) string {
	in, out := transcript.TotalTokens(ctx.Input.TranscriptPath)
	if in == 0 && out == 0 {
		return ""
	}
	t := ctx.Theme
	return fmt.Sprintf("%s%s in%s / %s%s out%s",
		t.WellnessOk, humanK(in), t.Reset,
		t.Time, humanK(out), t.Reset)
}

// humanK renders 12345 as "12K", 1234567 as "1.2M".
func humanK(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1000:
		return fmt.Sprintf("%dK", n/1000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
