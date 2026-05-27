package components

import "time"

func init() {
	Register("time_display", renderTimeDisplay)
}

func renderTimeDisplay(_ *Context) string {
	return "🕐 " + time.Now().Format("15:04")
}
