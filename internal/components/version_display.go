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
	return "CC:" + ctx.Input.Version
}

// renderVersionInfo shows both CC version and our goccline version.
func renderVersionInfo(ctx *Context) string {
	cc := ""
	if ctx.Input.Version != "" {
		cc = "CC:" + ctx.Input.Version + " │ "
	}
	return cc + "GL:" + Version
}

// Version is set at build time via -ldflags "-X .../components.Version=vX.Y.Z".
var Version = "dev"
