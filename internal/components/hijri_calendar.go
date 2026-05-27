package components

import (
	"strings"
	"time"

	"github.com/andhikapraa/goccline/internal/hijri"
)

func init() {
	Register("hijri_calendar", renderHijriCalendar)
	Register("location_display", renderLocationDisplay)
}

// renderHijriCalendar prints "🕌 10 Dhu al-Hijjah 1447 🌔 │ May 27 2026" —
// the Hijri date, lunar phase emoji, and Gregorian date.
func renderHijriCalendar(ctx *Context) string {
	d, err := hijri.Today()
	if err != nil || d.MonthName == "" {
		return ""
	}
	t := ctx.Theme
	parts := []string{
		"🕌 " + t.Hijri + d.Format() + t.Reset,
		t.HijriMoon + hijri.MoonPhase(d.Day) + t.Reset,
		time.Now().Format("Jan 2 2006"),
	}
	if label := ctx.Config.Prayer.LocationLabel; label != "" {
		parts = append(parts, t.Location+"📍 "+label+t.Reset)
	}
	return strings.Join(parts, " │ ")
}

// renderLocationDisplay prints just the user's location label (if set).
// Useful when you want location on a different line from hijri_calendar.
func renderLocationDisplay(ctx *Context) string {
	label := ctx.Config.Prayer.LocationLabel
	if label == "" {
		return ""
	}
	t := ctx.Theme
	return t.Location + "📍 " + label + t.Reset
}
