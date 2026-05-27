// Package config loads the user's TOML config for goccline.
//
// The format mirrors the bash claude-code-statusline so users can migrate
// existing Config.toml files. Unknown keys are tolerated.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Display struct {
	Lines int    `toml:"lines"`
	Line  []Line `toml:"-"` // populated from line1/line2/... flat keys
}

type Line struct {
	Components    []string `toml:"components"`
	Separator     string   `toml:"separator"`
	ShowWhenEmpty bool     `toml:"show_when_empty"`
}

type Theme struct {
	Name string `toml:"name"`
}

type Features struct {
	ShowCommits     bool `toml:"show_commits"`
	ShowVersion     bool `toml:"show_version"`
	ShowSubmodules  bool `toml:"show_submodules"`
	ShowCost        bool `toml:"show_cost_tracking"`
	ShowMCP         bool `toml:"show_mcp_status"`
	ShowPrayerTimes bool `toml:"show_prayer_times"`
	ShowHijriDate   bool `toml:"show_hijri_date"`
	ShowContext     bool `toml:"show_context_window"`
	ShowUsageLimits bool `toml:"show_usage_limits"`
}

type SessionInfo struct {
	// ShowFirstMessage reads the transcript file to show the first user
	// prompt next to the session ID. Off by default — transcript I/O can
	// add ~130ms per render on large files.
	ShowFirstMessage bool `toml:"show_first_message"`
	// IDLength controls how many characters of the session UUID to show.
	IDLength int `toml:"id_length"`
	// TitleLength bounds the first-message preview when enabled.
	TitleLength int `toml:"title_length"`
}

// Prayer configures Islamic prayer time fetching and display.
type Prayer struct {
	// LocationMode controls coordinate source. "gps" calls
	// CoreLocationCLI (macOS only) and caches for 1h. "manual" (or empty)
	// uses the lat/lon below.
	LocationMode string  `toml:"location_mode"`
	Latitude     float64 `toml:"latitude"`
	Longitude    float64 `toml:"longitude"`
	// Method is the AlAdhan calculation method (2 = Islamic Society of
	// North America). See https://aladhan.com/calculation-methods.
	Method int `toml:"method"`
	// School: 1 = Hanafi (later Asr), 0 = Shafi/Maliki/Hanbali.
	School int `toml:"school"`
	// LocationLabel is displayed when present (e.g. "Jakarta, ID").
	LocationLabel string `toml:"location_label"`
}

type Config struct {
	Display     Display     `toml:"display"`
	Theme       Theme       `toml:"theme"`
	Features    Features    `toml:"features"`
	SessionInfo SessionInfo `toml:"session_info"`
	Prayer      Prayer      `toml:"prayer"`
}

// DefaultPath returns the user's preferred config location.
func DefaultPath() string {
	if v := os.Getenv("GOCCLINE_CONFIG"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "goccline", "Config.toml")
}

// Load reads and parses the config file. Returns a default config if the
// file is missing.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Display: Display{Lines: 2},
		Theme:   Theme{Name: "catppuccin"},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Display.Line = []Line{
				{Components: []string{"repo_info"}, Separator: " │ ", ShowWhenEmpty: true},
				{Components: []string{"model_info", "time_display"}, Separator: " │ ", ShowWhenEmpty: true},
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	// The bash format uses flat dot-keys like `display.line1.components`.
	// BurntSushi/toml understands both, but we need to gather lineN keys
	// manually because the line count is dynamic.
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse toml: %w", err)
	}

	if d, ok := raw["display"].(map[string]any); ok {
		if v, ok := d["lines"].(int64); ok {
			cfg.Display.Lines = int(v)
		}
		for i := 1; i <= cfg.Display.Lines; i++ {
			key := fmt.Sprintf("line%d", i)
			lineMap, ok := d[key].(map[string]any)
			if !ok {
				cfg.Display.Line = append(cfg.Display.Line, Line{})
				continue
			}
			ln := Line{Separator: " │ ", ShowWhenEmpty: true}
			if comps, ok := lineMap["components"].([]any); ok {
				for _, c := range comps {
					if s, ok := c.(string); ok {
						ln.Components = append(ln.Components, s)
					}
				}
			}
			if s, ok := lineMap["separator"].(string); ok {
				ln.Separator = s
			}
			if b, ok := lineMap["show_when_empty"].(bool); ok {
				ln.ShowWhenEmpty = b
			}
			cfg.Display.Line = append(cfg.Display.Line, ln)
		}
	}

	if t, ok := raw["theme"].(map[string]any); ok {
		if s, ok := t["name"].(string); ok {
			cfg.Theme.Name = s
		}
	}

	if f, ok := raw["features"].(map[string]any); ok {
		assignBool(f, "show_commits", &cfg.Features.ShowCommits)
		assignBool(f, "show_version", &cfg.Features.ShowVersion)
		assignBool(f, "show_submodules", &cfg.Features.ShowSubmodules)
		assignBool(f, "show_cost_tracking", &cfg.Features.ShowCost)
		assignBool(f, "show_mcp_status", &cfg.Features.ShowMCP)
		assignBool(f, "show_prayer_times", &cfg.Features.ShowPrayerTimes)
		assignBool(f, "show_hijri_date", &cfg.Features.ShowHijriDate)
		assignBool(f, "show_context_window", &cfg.Features.ShowContext)
		assignBool(f, "show_usage_limits", &cfg.Features.ShowUsageLimits)
	}

	if s, ok := raw["session_info"].(map[string]any); ok {
		assignBool(s, "show_first_message", &cfg.SessionInfo.ShowFirstMessage)
		// Bash statusline compat alias.
		assignBool(s, "show_project", &cfg.SessionInfo.ShowFirstMessage)
		assignInt(s, "id_length", &cfg.SessionInfo.IDLength)
		assignInt(s, "title_length", &cfg.SessionInfo.TitleLength)
	}

	if p, ok := raw["prayer"].(map[string]any); ok {
		assignString(p, "location_mode", &cfg.Prayer.LocationMode)
		assignFloat(p, "latitude", &cfg.Prayer.Latitude)
		assignFloat(p, "longitude", &cfg.Prayer.Longitude)
		assignInt(p, "method", &cfg.Prayer.Method)
		assignInt(p, "school", &cfg.Prayer.School)
		assignString(p, "location_label", &cfg.Prayer.LocationLabel)
	}

	return cfg, nil
}

func assignFloat(m map[string]any, key string, dest *float64) {
	switch v := m[key].(type) {
	case float64:
		*dest = v
	case int64:
		*dest = float64(v)
	}
}

func assignString(m map[string]any, key string, dest *string) {
	if v, ok := m[key].(string); ok {
		*dest = v
	}
}

func assignInt(m map[string]any, key string, dest *int) {
	if v, ok := m[key].(int64); ok {
		*dest = int(v)
	}
}

func assignBool(m map[string]any, key string, dest *bool) {
	if v, ok := m[key].(bool); ok {
		*dest = v
	}
}
