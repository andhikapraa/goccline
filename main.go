// Command goccline is a fast Claude Code statusline.
//
// Usage:
//
//	echo '{"workspace":{"current_dir":"/tmp"},...}' | goccline
//
// Config lives at $GOCCLINE_CONFIG or ~/.claude/goccline/Config.toml.
package main

import (
	"fmt"
	"os"

	_ "github.com/andhikapraa/goccline/internal/components" // register built-ins

	"github.com/andhikapraa/goccline/internal/config"
	"github.com/andhikapraa/goccline/internal/input"
	"github.com/andhikapraa/goccline/internal/render"
)

func main() {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "goccline: config:", err)
		os.Exit(2)
	}

	in, err := input.Parse(os.Stdin)
	if err != nil {
		// Don't abort — render with whatever defaults we have.
		fmt.Fprintln(os.Stderr, "goccline: input:", err)
	}

	fmt.Println(render.Render(cfg, in))
}
