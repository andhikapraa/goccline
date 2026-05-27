package components

func init() {
	Register("session_mode", renderSessionMode)
	Register("agent_display", renderAgentDisplay)
}

// renderSessionMode shows the output_style from Claude Code (e.g. "default",
// "explanatory"). The bash statusline labels it "Style:".
func renderSessionMode(ctx *Context) string {
	name := ctx.Input.OutputStyle.Name
	if name == "" {
		return ""
	}
	return "Style: " + name
}

// renderAgentDisplay is a placeholder until Claude Code surfaces the agent
// name on the JSON payload. For now it mirrors session_mode.
func renderAgentDisplay(ctx *Context) string {
	return renderSessionMode(ctx)
}
