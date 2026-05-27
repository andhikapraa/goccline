package components

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/andhikapraa/goccline/internal/git"
)

func init() {
	Register("repo_info", renderRepoInfo)
}

// renderRepoInfo prints the working directory (~-shortened) and, if the
// directory is a git repo, the current branch.
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
	dir := filepath.Clean(cwd)

	if branch := git.Branch(ctx.Input.CWD); branch != "" {
		dir += " ⎇ " + branch
		if git.HasWorktree(ctx.Input.CWD) {
			dir += " [WT]"
		}
	}
	return dir
}
