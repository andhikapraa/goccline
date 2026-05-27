package components

import "fmt"

func init() {
	Register("code_productivity", renderCodeProductivity)
}

// renderCodeProductivity reads cost.total_lines_added / total_lines_removed
// from the Claude Code JSON payload. Free — no I/O.
// Format: "+123 / -45" matches the bash statusline.
func renderCodeProductivity(ctx *Context) string {
	added := ctx.Input.Cost.TotalLinesAdded
	removed := ctx.Input.Cost.TotalLinesRemvd
	if added == 0 && removed == 0 {
		return ""
	}
	t := ctx.Theme
	return fmt.Sprintf("%s+%d%s / %s-%d%s",
		t.WellnessOk, added, t.Reset,
		t.WellnessNudge, removed, t.Reset)
}
