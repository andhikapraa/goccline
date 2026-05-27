package components

import "strings"

func init() {
	Register("model_info", renderModelInfo)
}

// renderModelInfo emits a model emoji + name. The bash statusline picks the
// emoji from the model ID (opus/haiku/sonnet); we mirror that mapping.
func renderModelInfo(ctx *Context) string {
	name := ctx.Input.Model.DisplayName
	if name == "" {
		name = ctx.Input.Model.ID
	}
	if name == "" {
		return ""
	}

	t := ctx.Theme
	emoji := "🤖"
	color := t.Model
	id := strings.ToLower(ctx.Input.Model.ID)
	switch {
	case strings.Contains(id, "opus"):
		emoji = "🧠"
		color = t.Opus
	case strings.Contains(id, "haiku"):
		emoji = "⚡"
		color = t.Haiku
	case strings.Contains(id, "sonnet"):
		emoji = "🎵"
		color = t.Sonnet
	}

	return emoji + " " + color + name + t.Reset
}
