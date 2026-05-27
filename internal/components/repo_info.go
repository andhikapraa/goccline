package components

import (
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Register("repo_info", renderRepoInfo)
}

// renderRepoInfo prints the working directory, shortened with ~ for $HOME.
// Git branch follows in a future iteration; this is the minimal v0.1.
func renderRepoInfo(ctx *Context) string {
	cwd := ctx.Input.CWD
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	if cwd == "" {
		return ""
	}
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(cwd, home) {
		cwd = "~" + strings.TrimPrefix(cwd, home)
	}
	return filepath.Clean(cwd)
}
