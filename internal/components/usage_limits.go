package components

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andhikapraa/goccline/internal/theme"
)

func init() {
	Register("usage_limits", renderUsageLimits)
	Register("reset_timer", renderResetTimer)
}

// renderUsageLimits prints "⏱ Used 5h:22% • 7d:54%" using the rate_limits
// block Claude Code v2.1.80+ ships in the statusLine JSON. Percentages are
// USED (not remaining) — the rest is your headroom.
func renderUsageLimits(ctx *Context) string {
	rl := ctx.Input.RateLimits
	five := rl.FiveHour.UsedPercentage
	seven := rl.SevenDay.UsedPercentage
	if five == 0 && seven == 0 {
		return ""
	}
	parts := []string{}
	if five > 0 {
		parts = append(parts, "5h:"+colorPct(five, ctx.Theme))
	}
	if seven > 0 {
		parts = append(parts, "7d:"+colorPct(seven, ctx.Theme))
	}
	if len(parts) == 0 {
		return ""
	}
	return "⏱ Used " + strings.Join(parts, " • ")
}

// renderResetTimer prints "⏱ 5h:2h5m • 7d:Sun 8AM" — the time-to-next-reset
// for each rate window.
func renderResetTimer(ctx *Context) string {
	rl := ctx.Input.RateLimits
	parts := []string{}
	if d := timeUntilAny(rl.FiveHour.ResetsAt); d != "" {
		parts = append(parts, "5h:"+d)
	}
	if d := timeUntilAny(rl.SevenDay.ResetsAt); d != "" {
		parts = append(parts, "7d:"+d)
	}
	if len(parts) == 0 {
		return ""
	}
	return "⏱ " + ctx.Theme.Dim + strings.Join(parts, " • ") + ctx.Theme.Reset
}

// timeUntilAny coerces the JSON value into a time and returns the formatted
// remaining duration. Handles three shapes Claude Code has emitted across
// versions: float64 (e.g. 1.78e9 — current v2.1.80+), int64, and ISO string.
func timeUntilAny(raw any) string {
	switch v := raw.(type) {
	case nil:
		return ""
	case float64:
		return timeUntil(strconv.FormatInt(int64(v), 10))
	case int64:
		return timeUntil(strconv.FormatInt(v, 10))
	case int:
		return timeUntil(strconv.Itoa(v))
	case string:
		return timeUntil(v)
	default:
		return timeUntil(fmt.Sprint(v))
	}
}

// timeUntil parses an ISO 8601 timestamp or Unix epoch seconds and returns
// a human-friendly remaining duration.
//
//	< 1h   -> "Xm"
//	< 24h  -> "XhYm"
//	>= 24h -> "Mon 3PM"
func timeUntil(raw string) string {
	raw = strings.TrimSpace(strings.Trim(raw, `"`))
	if raw == "" || raw == "null" {
		return ""
	}
	var resetAt time.Time
	if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
		resetAt = time.Unix(n, 0)
	} else {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return ""
		}
		resetAt = t
	}
	d := time.Until(resetAt)
	if d <= 0 {
		return "now"
	}
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		h := int(d / time.Hour)
		m := int(d/time.Minute) % 60
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		return resetAt.Format("Mon 3PM")
	}
}

func colorPct(pct float64, t theme.Theme) string {
	color := t.WellnessOk
	switch {
	case pct >= 80:
		color = t.WellnessNudge
	case pct >= 50:
		color = t.Time
	}
	return fmt.Sprintf("%s%.0f%%%s", color, pct, t.Reset)
}
