// Package hijri converts Gregorian dates to Hijri (Islamic) calendar dates.
package hijri

import (
	"fmt"
	"time"

	gohijri "github.com/hablullah/go-hijri"
)

// Date holds the Hijri calendar fields we display.
type Date struct {
	Day       int
	Month     int
	MonthName string
	Year      int64
}

var monthNames = [...]string{
	"Muharram", "Safar", "Rabi al-Awwal", "Rabi al-Thani",
	"Jumada al-Awwal", "Jumada al-Thani", "Rajab", "Sha'ban",
	"Ramadan", "Shawwal", "Dhu al-Qi'dah", "Dhu al-Hijjah",
}

// Today returns the Hijri date for now.
func Today() (Date, error) {
	h, err := gohijri.CreateHijriDate(time.Now(), 0)
	if err != nil {
		return Date{}, err
	}
	d := Date{
		Day:   int(h.Day),
		Month: int(h.Month),
		Year:  h.Year,
	}
	if d.Month >= 1 && d.Month <= 12 {
		d.MonthName = monthNames[d.Month-1]
	}
	return d, nil
}

// Format renders the date as "10 Dhu al-Hijjah 1447".
func (d Date) Format() string {
	if d.MonthName == "" {
		return ""
	}
	return fmt.Sprintf("%d %s %d", d.Day, d.MonthName, d.Year)
}

// MoonPhase returns the emoji approximating the current lunar phase.
// Algorithm: a synodic month is ~29.53 days; map the Hijri day-of-month
// (which roughly tracks lunar age) to one of 8 phase emojis.
func MoonPhase(hijriDay int) string {
	phases := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	// Hijri days 1..30; new moon at 1, full near 14-15.
	idx := (hijriDay - 1) * 8 / 30
	if idx < 0 {
		idx = 0
	}
	if idx >= len(phases) {
		idx = len(phases) - 1
	}
	return phases[idx]
}
