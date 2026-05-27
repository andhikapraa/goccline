package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/andhikapraa/goccline/internal/cost"
)

func init() {
	Register("burn_rate", renderBurnRate)
	Register("block_projection", renderBlockProjection)
}

// blockInfo computes the current 5-hour billing block boundaries and the
// spend incurred so far within that block.
//
// Boundary source priority:
//  1. rate_limits.five_hour.resets_at from the Claude Code JSON payload
//     (the canonical answer when Claude Code v2.1.80+ provides it).
//  2. UTC-aligned 5h grid (00, 05, 10, 15, 20) as a fallback.
type blockInfo struct {
	Start      time.Time
	End        time.Time
	Now        time.Time
	SpentUSD   float64
	TokensUsed int
}

func currentBlock(ctx *Context) (blockInfo, bool) {
	bi := blockInfo{Now: time.Now()}

	// Source 1: Anthropic-reported reset time.
	if rawEnd := fmt.Sprint(ctx.Input.RateLimits.FiveHour.ResetsAt); rawEnd != "" && rawEnd != "<nil>" {
		raw := strings.Trim(rawEnd, `"`)
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
			bi.End = time.Unix(n, 0)
		} else if t, err := time.Parse(time.RFC3339, raw); err == nil {
			bi.End = t
		}
	}
	if bi.End.IsZero() {
		// Source 2: UTC 5h grid.
		utc := bi.Now.UTC()
		hourBucket := (utc.Hour() / 5) * 5
		startUTC := time.Date(utc.Year(), utc.Month(), utc.Day(), hourBucket, 0, 0, 0, time.UTC)
		bi.Start = startUTC.Local()
		bi.End = bi.Start.Add(5 * time.Hour)
	} else {
		bi.Start = bi.End.Add(-5 * time.Hour)
	}

	// Sum cost from session JSONL files for entries in [Start, Now].
	root := projectsRoot()
	if root == "" {
		return bi, false
	}
	entries := ctx.CostMemo.Scan(root, bi.Start)
	for _, e := range entries {
		if e.Timestamp.Before(bi.Start) || e.Timestamp.After(bi.Now) {
			continue
		}
		bi.SpentUSD += cost.CostOf(e.Model, e.Usage)
		bi.TokensUsed += e.Usage.InputTokens + e.Usage.OutputTokens +
			e.Usage.CacheReadTokens + e.Usage.CacheWrite5m + e.Usage.CacheWrite1h
	}
	return bi, true
}

// projectsRoot is defined in cost.go; we reference it here too. Keeping the
// reference for clarity even though Go shares the package scope.
var _ = filepath.Join
var _ = os.UserHomeDir

// renderBurnRate prints "🔥 $X.XX/hr" — the current 5h block's spend rate.
func renderBurnRate(ctx *Context) string {
	bi, ok := currentBlock(ctx)
	if !ok || bi.SpentUSD == 0 {
		return ""
	}
	elapsed := bi.Now.Sub(bi.Start)
	if elapsed <= 0 {
		return ""
	}
	perHour := bi.SpentUSD / elapsed.Hours()
	t := ctx.Theme
	color := t.WellnessOk
	switch {
	case perHour >= 10:
		color = t.WellnessNudge
	case perHour >= 3:
		color = t.Time
	}
	return fmt.Sprintf("🔥 %s$%.2f/hr%s", color, perHour, t.Reset)
}

// renderBlockProjection prints "Est: $X.XX (YK)" — projected end-of-block
// cost and token count, extrapolated linearly from the current burn rate.
func renderBlockProjection(ctx *Context) string {
	bi, ok := currentBlock(ctx)
	if !ok || bi.SpentUSD == 0 {
		return ""
	}
	elapsed := bi.Now.Sub(bi.Start)
	total := bi.End.Sub(bi.Start)
	if elapsed <= 0 || total <= 0 {
		return ""
	}
	scale := total.Seconds() / elapsed.Seconds()
	projCost := bi.SpentUSD * scale
	projTokens := int(float64(bi.TokensUsed) * scale)
	t := ctx.Theme
	return fmt.Sprintf("Est: %s$%.2f%s (%s)",
		t.CostLive, projCost, t.Reset, humanK(projTokens))
}
