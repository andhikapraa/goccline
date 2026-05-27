package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/andhikapraa/goccline/internal/prayer"
)

func init() {
	Register("prayer_icon", renderPrayerIcon)
	Register("prayer_times", renderPrayerTimes)
	Register("prayer_times_only", renderPrayerTimesOnly)
}

// renderPrayerIcon returns the mosque emoji unconditionally — meant to live
// next to prayer_times for a small visual anchor.
func renderPrayerIcon(_ *Context) string { return "🕌" }

// renderPrayerTimes prints all five prayers with a checkmark for ones that
// have already passed and a "Next prayer (Xm)" annotation for the upcoming one.
func renderPrayerTimes(ctx *Context) string {
	return renderPrayerInternal(ctx, true)
}

// renderPrayerTimesOnly is the same display without the leading 🕌 prefix
// (so it composes well after a separate prayer_icon component).
func renderPrayerTimesOnly(ctx *Context) string {
	return renderPrayerInternal(ctx, false)
}

func renderPrayerInternal(ctx *Context, includeIcon bool) string {
	cfg := ctx.Config.Prayer
	if cfg.Latitude == 0 && cfg.Longitude == 0 {
		// No location configured → silently skip rather than confusing display.
		return ""
	}
	method := cfg.Method
	if method == 0 {
		method = 2 // ISNA default
	}
	t, err := prayer.Fetch(cfg.Latitude, cfg.Longitude, method, cfg.School)
	if err != nil {
		return ""
	}

	type slot struct{ name, value string }
	slots := []slot{
		{"Fajr", t.Fajr},
		{"Dhuhr", t.Dhuhr},
		{"Asr", t.Asr},
		{"Maghrib", t.Maghrib},
		{"Isha", t.Isha},
	}

	now := time.Now()
	nextIdx := -1
	for i, s := range slots {
		if parseClock(s.value, now).After(now) {
			nextIdx = i
			break
		}
	}

	parts := make([]string, 0, len(slots))
	for i, s := range slots {
		switch {
		case nextIdx == -1 || i < nextIdx:
			parts = append(parts, fmt.Sprintf("%s %s ✓", s.name, s.value))
		case i == nextIdx:
			remaining := time.Until(parseClock(s.value, now))
			parts = append(parts, fmt.Sprintf("%s %s (%s)", s.name, s.value, humanDuration(remaining)))
		default:
			parts = append(parts, fmt.Sprintf("%s %s", s.name, s.value))
		}
	}

	out := strings.Join(parts, " │ ")
	if includeIcon {
		out = "🕌 " + out
	}
	return out
}

// parseClock combines a HH:MM string with today's date in the current TZ.
func parseClock(hhmm string, today time.Time) time.Time {
	t, err := time.ParseInLocation("15:04", hhmm, today.Location())
	if err != nil {
		return today
	}
	return time.Date(today.Year(), today.Month(), today.Day(),
		t.Hour(), t.Minute(), 0, 0, today.Location())
}

func humanDuration(d time.Duration) string {
	if d <= 0 {
		return "now"
	}
	h := int(d / time.Hour)
	m := int(d/time.Minute) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
