// Package render assembles the configured lines into the final ANSI output.
package render

import (
	"strings"

	"github.com/andhikapraa/goccline/internal/components"
	"github.com/andhikapraa/goccline/internal/config"
	"github.com/andhikapraa/goccline/internal/cost"
	"github.com/andhikapraa/goccline/internal/input"
	"github.com/andhikapraa/goccline/internal/theme"
)

// Render walks display.line1..display.lines, invokes each component, and
// joins the results. Returns the full multi-line statusline.
func Render(cfg *config.Config, in input.Payload) string {
	ctx := &components.Context{
		Input:    in,
		Config:   cfg,
		CostMemo: cost.NewMemo(),
		Theme:    theme.Resolve(cfg.Theme.Name),
	}

	var b strings.Builder
	for i, line := range cfg.Display.Line {
		if i >= cfg.Display.Lines {
			break
		}
		parts := make([]string, 0, len(line.Components))
		for _, name := range line.Components {
			fn := components.Lookup(name)
			if fn == nil {
				continue
			}
			if s := fn(ctx); s != "" {
				parts = append(parts, s)
			}
		}
		if len(parts) == 0 && !line.ShowWhenEmpty {
			continue
		}
		sep := line.Separator
		if sep == "" {
			sep = " │ "
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(strings.Join(parts, sep))
	}
	return b.String()
}
