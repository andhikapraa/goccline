package components

func init() {
	Register("version_display", renderVersionDisplay)
	Register("version_info", renderVersionInfo)
}

// renderVersionDisplay shows the Claude Code version from the JSON payload.
func renderVersionDisplay(ctx *Context) string {
	if ctx.Input.Version == "" {
		return ""
	}
	t := ctx.Theme
	return t.Version + "CC:" + ctx.Input.Version + t.Reset
}

// renderVersionInfo shows both CC version and our goccline version.
func renderVersionInfo(ctx *Context) string {
	t := ctx.Theme
	out := ""
	if ctx.Input.Version != "" {
		out = t.Version + "CC:" + ctx.Input.Version + t.Reset + " │ "
	}
	return out + t.Model + "GL:" + Version + t.Reset
}

// Version is set at build time via -ldflags "-X .../components.Version=vX.Y.Z".
var Version = "dev"
