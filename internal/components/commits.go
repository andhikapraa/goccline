package components

import (
	"fmt"

	"github.com/andhikapraa/goccline/internal/git"
)

func init() {
	Register("commits", renderCommits)
}

func renderCommits(ctx *Context) string {
	dir := ctx.Input.CWD
	if dir == "" {
		return ""
	}
	n := git.CommitsToday(dir)
	return fmt.Sprintf("Commits:%d", n)
}
