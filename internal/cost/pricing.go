// Package cost computes Claude Code spend from session JSONL files.
package cost

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed pricing_data.json
var pricingJSON []byte

// Pricing rates are USD per 1M tokens, except WebSearchCost which is per request.
type Pricing struct {
	Input        float64 `json:"input"`
	Output       float64 `json:"output"`
	CacheRead    float64 `json:"cache_read"`
	CacheWrite5m float64 `json:"cache_write_5m"`
	CacheWrite1h float64 `json:"cache_write_1h"`
}

type pricingFile struct {
	Models  map[string]Pricing `json:"models"`
	Aliases []struct {
		Prefix string `json:"prefix"`
		Model  string `json:"model"`
	} `json:"aliases"`
	DefaultModel        string  `json:"default_model"`
	LongContextThresh   int     `json:"long_context_threshold"`
	WebSearchCost       float64 `json:"web_search_cost"`
}

var pricing pricingFile

func init() {
	if err := json.Unmarshal(pricingJSON, &pricing); err != nil {
		// embed/parse failure is a build-time issue; fall back to empty map
		// so the binary still runs (cost components will return 0).
		pricing.Models = map[string]Pricing{}
	}
}

// LookupPricing resolves a model string (e.g. "claude-opus-4-7-20251001")
// against the pricing table. Strategy:
//  1. Exact match on the table key.
//  2. Try progressive prefix matches by stripping date/variant suffixes.
//  3. Match via aliases (prefix mapping).
//  4. Fall through to the default model.
func LookupPricing(model string) (Pricing, bool) {
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" {
		return Pricing{}, false
	}
	if p, ok := pricing.Models[m]; ok {
		return p, true
	}
	// Strip trailing -YYYYMMDD or fast variants.
	for _, suffix := range []string{":fast"} {
		m = strings.TrimSuffix(m, suffix)
	}
	if i := strings.LastIndex(m, "-"); i > 0 {
		// If the trailing segment is all digits, drop it.
		if isDigits(m[i+1:]) {
			if p, ok := pricing.Models[m[:i]]; ok {
				return p, true
			}
		}
	}
	for _, alias := range pricing.Aliases {
		if strings.HasPrefix(m, alias.Prefix) {
			if p, ok := pricing.Models[alias.Model]; ok {
				return p, true
			}
		}
	}
	if pricing.DefaultModel != "" {
		if p, ok := pricing.Models[pricing.DefaultModel]; ok {
			return p, true
		}
	}
	return Pricing{}, false
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// CostOf returns the dollar cost of a single usage entry.
func CostOf(model string, u Usage) float64 {
	p, ok := LookupPricing(model)
	if !ok {
		return 0
	}
	c := float64(u.InputTokens)*p.Input +
		float64(u.OutputTokens)*p.Output +
		float64(u.CacheReadTokens)*p.CacheRead +
		float64(u.CacheWrite5m)*p.CacheWrite5m +
		float64(u.CacheWrite1h)*p.CacheWrite1h
	return c / 1e6
}
