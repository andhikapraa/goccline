// Package theme holds truecolor palettes ported from the bash
// claude-code-statusline. Components reference colors by semantic name
// (Path, Branch, ModelOpus, etc.) and the active theme resolves them.
package theme

const (
	// reset clears all SGR styling.
	reset = "\x1b[0m"
)

// Theme is the per-render palette. Colors are full ANSI SGR escape
// sequences ready to write to stdout. The Reset field is "\x1b[0m".
type Theme struct {
	Name string

	// Semantic slots.
	Path        string
	Branch      string
	Clean       string
	Dirty       string
	Session     string
	SessionText string
	Model       string
	Sonnet      string
	Haiku       string
	Opus        string
	Style       string
	Version     string
	Time        string
	CostLive    string
	CostToday   string
	CostWeek    string
	CostMonth   string
	CostRepo    string
	Hijri       string
	HijriMoon   string
	Wellness    string
	WellnessOk  string
	WellnessNudge string
	Location    string
	Prayer      string
	PrayerDone  string
	PrayerNext  string
	Dim         string
	Reset       string
}

// truecolor returns the ANSI 24-bit SGR foreground sequence.
func truecolor(r, g, b int) string {
	return tc(r, g, b)
}

func tc(r, g, b int) string {
	// Fast-path: small static buffer. Keep imports minimal.
	return "\x1b[38;2;" + itoa(r) + ";" + itoa(g) + ";" + itoa(b) + "m"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var buf [4]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// Catppuccin is the Catppuccin Mocha palette.
var Catppuccin = Theme{
	Name: "catppuccin",

	Path:          truecolor(137, 180, 250), // blue
	Branch:        truecolor(203, 166, 247), // mauve
	Clean:         "",                       // emoji carries the meaning
	Dirty:         "",
	Session:       truecolor(137, 220, 235), // cyan
	SessionText:   truecolor(205, 214, 244), // white-ish
	Model:         truecolor(137, 220, 235), // cyan
	Sonnet:        truecolor(166, 227, 161), // green
	Haiku:         truecolor(249, 226, 175), // yellow
	Opus:          truecolor(137, 220, 235), // cyan
	Style:         truecolor(166, 227, 161), // green
	Version:       truecolor(203, 166, 247), // mauve
	Time:          truecolor(250, 179, 135), // peach
	CostLive:      truecolor(249, 226, 175), // yellow (live)
	CostToday:     truecolor(166, 227, 161), // green
	CostWeek:      truecolor(137, 180, 250), // blue
	CostMonth:     truecolor(203, 166, 247), // mauve
	CostRepo:      truecolor(245, 194, 231), // pink
	Hijri:         truecolor(148, 226, 213), // teal
	HijriMoon:     truecolor(249, 226, 175), // yellow
	Wellness:      truecolor(250, 179, 135), // peach
	WellnessOk:    truecolor(166, 227, 161), // green
	WellnessNudge: truecolor(243, 139, 168), // red
	Location:      truecolor(166, 173, 200), // light gray
	Prayer:        truecolor(205, 214, 244), // white-ish
	PrayerDone:    truecolor(166, 173, 200) + "\x1b[2m",
	PrayerNext:    truecolor(166, 227, 161), // green
	Dim:           "\x1b[2m",
	Reset:         reset,
}

// Garden is the soft pastel theme.
var Garden = Theme{
	Name: "garden",

	Path:          truecolor(173, 216, 230),
	Branch:        truecolor(221, 160, 221),
	Session:       truecolor(176, 196, 145),
	SessionText:   truecolor(255, 255, 255),
	Model:         truecolor(176, 196, 145),
	Sonnet:        truecolor(176, 196, 145),
	Haiku:         truecolor(255, 218, 185),
	Opus:          truecolor(173, 216, 230),
	Style:         truecolor(176, 196, 145),
	Version:       truecolor(221, 160, 221),
	Time:          truecolor(255, 218, 185),
	CostLive:      truecolor(255, 218, 185),
	CostToday:     truecolor(176, 196, 145),
	CostWeek:      truecolor(173, 216, 230),
	CostMonth:     truecolor(221, 160, 221),
	CostRepo:      truecolor(255, 182, 193),
	Hijri:         truecolor(176, 196, 145),
	HijriMoon:     truecolor(255, 218, 185),
	Wellness:      truecolor(255, 218, 185),
	WellnessOk:    truecolor(176, 196, 145),
	WellnessNudge: truecolor(255, 182, 193),
	Location:      truecolor(192, 192, 192),
	Prayer:        truecolor(255, 255, 255),
	PrayerDone:    truecolor(160, 160, 160) + "\x1b[2m",
	PrayerNext:    truecolor(176, 196, 145),
	Dim:           "\x1b[2m",
	Reset:         reset,
}

// Ocean is the deep-sea palette.
var Ocean = Theme{
	Name: "ocean",

	Path:          truecolor(72, 202, 228),
	Branch:        truecolor(132, 220, 198),
	Session:       truecolor(132, 220, 198),
	SessionText:   truecolor(220, 240, 255),
	Model:         truecolor(72, 202, 228),
	Sonnet:        truecolor(76, 201, 176),
	Haiku:         truecolor(255, 209, 102),
	Opus:          truecolor(72, 202, 228),
	Style:         truecolor(76, 201, 176),
	Version:       truecolor(155, 89, 182),
	Time:          truecolor(255, 209, 102),
	CostLive:      truecolor(255, 209, 102),
	CostToday:     truecolor(76, 201, 176),
	CostWeek:      truecolor(72, 202, 228),
	CostMonth:     truecolor(155, 89, 182),
	CostRepo:      truecolor(255, 107, 107),
	Hijri:         truecolor(76, 201, 176),
	HijriMoon:     truecolor(255, 209, 102),
	Wellness:      truecolor(255, 209, 102),
	WellnessOk:    truecolor(76, 201, 176),
	WellnessNudge: truecolor(255, 107, 107),
	Location:      truecolor(180, 200, 220),
	Prayer:        truecolor(220, 240, 255),
	PrayerDone:    truecolor(160, 180, 200) + "\x1b[2m",
	PrayerNext:    truecolor(76, 201, 176),
	Dim:           "\x1b[2m",
	Reset:         reset,
}

// Classic is the standard 16-color ANSI palette (no truecolor).
var Classic = Theme{
	Name: "classic",

	Path:          "\x1b[34m", // blue
	Branch:        "\x1b[35m", // magenta
	Session:       "\x1b[36m", // cyan
	SessionText:   "\x1b[37m", // white
	Model:         "\x1b[36m",
	Sonnet:        "\x1b[32m",
	Haiku:         "\x1b[33m",
	Opus:          "\x1b[36m",
	Style:         "\x1b[32m",
	Version:       "\x1b[35m",
	Time:          "\x1b[33m",
	CostLive:      "\x1b[33m",
	CostToday:     "\x1b[32m",
	CostWeek:      "\x1b[34m",
	CostMonth:     "\x1b[35m",
	CostRepo:      "\x1b[31m",
	Hijri:         "\x1b[36m",
	HijriMoon:     "\x1b[33m",
	Wellness:      "\x1b[33m",
	WellnessOk:    "\x1b[32m",
	WellnessNudge: "\x1b[31m",
	Location:      "\x1b[37m",
	Prayer:        "\x1b[37m",
	PrayerDone:    "\x1b[2;37m",
	PrayerNext:    "\x1b[32m",
	Dim:           "\x1b[2m",
	Reset:         reset,
}

// Resolve returns the theme matching name, falling back to Catppuccin.
func Resolve(name string) Theme {
	switch name {
	case "garden":
		return Garden
	case "ocean":
		return Ocean
	case "classic":
		return Classic
	default:
		return Catppuccin
	}
}

// Wrap is a small helper: color + text + reset.
func Wrap(color, text string) string {
	if color == "" {
		return text
	}
	return color + text + reset
}
