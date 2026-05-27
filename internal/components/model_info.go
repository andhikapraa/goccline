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

	emoji := "🤖"
	id := strings.ToLower(ctx.Input.Model.ID)
	switch {
	case strings.Contains(id, "opus"):
		emoji = "🧠"
	case strings.Contains(id, "haiku"):
		emoji = "⚡"
	case strings.Contains(id, "sonnet"):
		emoji = "🎵"
	}

	return emoji + " " + name
}
