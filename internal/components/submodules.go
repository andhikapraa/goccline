package components

import (
	"fmt"

	"github.com/andhikapraa/goccline/internal/git"
)

func init() {
	Register("submodules", renderSubmodules)
}

func renderSubmodules(ctx *Context) string {
	dir := ctx.Input.CWD
	if dir == "" {
		return ""
	}
	n := git.SubmoduleCount(dir)
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("SUB:%d", n)
}
