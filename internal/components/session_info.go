package components

func init() {
	Register("session_info", renderSessionInfo)
}

func renderSessionInfo(ctx *Context) string {
	id := ctx.Input.SessionID
	if id == "" {
		return ""
	}
	if len(id) > 8 {
		id = id[:8]
	}
	return "🔗 " + id
}
