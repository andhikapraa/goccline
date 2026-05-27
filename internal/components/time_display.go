package components

import "time"

func init() {
	Register("time_display", renderTimeDisplay)
}

func renderTimeDisplay(ctx *Context) string {
	t := ctx.Theme
	return "🕐 " + t.Time + time.Now().Format("15:04") + t.Reset
}
