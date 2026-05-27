package components

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andhikapraa/goccline/internal/cost"
)

func init() {
	Register("cost_live", renderCostLive)
	Register("cost_today", renderCostToday)
	Register("cost_daily", renderCostToday) // alias matching bash naming
	Register("cost_weekly", renderCostWeekly)
	Register("cost_monthly", renderCostMonthly)
	Register("cost_repo", renderCostRepo)
}

// renderCostLive reads cost.total_cost_usd from Claude Code's stdin JSON.
// Effectively free — no file I/O.
func renderCostLive(ctx *Context) string {
	c := ctx.Input.Cost.TotalCostUSD
	if c <= 0 {
		return ""
	}
	return fmt.Sprintf("💰 $%.2f", c)
}

// renderCostToday sums spend across all projects from midnight local.
func renderCostToday(ctx *Context) string {
	root := projectsRoot()
	if root == "" {
		return ""
	}
	since := cost.StartOfDay(time.Now())
	entries := ctx.CostMemo.Scan(root, since)
	total := cost.SumSince(entries, since)
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("📅 $%.2f today", total)
}

// renderCostWeekly sums spend across the last 7 days.
func renderCostWeekly(ctx *Context) string {
	root := projectsRoot()
	if root == "" {
		return ""
	}
	since := time.Now().AddDate(0, 0, -7)
	entries := ctx.CostMemo.Scan(root, since)
	total := cost.SumSince(entries, since)
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("📊 $%.2f 7d", total)
}

// renderCostMonthly sums spend across the last 30 days.
func renderCostMonthly(ctx *Context) string {
	root := projectsRoot()
	if root == "" {
		return ""
	}
	since := time.Now().AddDate(0, 0, -30)
	entries := ctx.CostMemo.Scan(root, since)
	total := cost.SumSince(entries, since)
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("📈 $%.2f 30d", total)
}

// renderCostRepo sums all-time spend for the current project's transcript
// directory only. Limits the scan to ~/.claude/projects/<encoded-cwd>/.
func renderCostRepo(ctx *Context) string {
	root := projectsRoot()
	if root == "" {
		return ""
	}
	cwd := ctx.Input.CWD
	if cwd == "" {
		return ""
	}
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return ""
	}
	repoDir := filepath.Join(root, cost.EncodeProjectKey(abs))
	if _, err := os.Stat(repoDir); err != nil {
		return ""
	}
	entries := cost.Scan(repoDir, time.Time{})
	total := cost.Sum(entries)
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("📁 $%.2f repo", total)
}

// projectsRoot returns the Claude Code projects directory.
func projectsRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "projects")
}
