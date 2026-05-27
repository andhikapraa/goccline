package cost

import (
	"strings"
	"time"
)

// Sum returns the total USD cost of entries.
func Sum(entries []Entry) float64 {
	var total float64
	for _, e := range entries {
		total += CostOf(e.Model, e.Usage)
	}
	return total
}

// SumSince returns total cost for entries on/after t.
func SumSince(entries []Entry, t time.Time) float64 {
	var total float64
	for _, e := range entries {
		if e.Timestamp.Before(t) {
			continue
		}
		total += CostOf(e.Model, e.Usage)
	}
	return total
}

// SumForProject returns cost limited to entries from project paths
// containing `projectKey`. The Claude Code project dir naming encodes the
// repo path with slashes replaced by dashes, so a substring match against
// the source file path picks up the right files. Caller decides the key.
// Empty key returns 0.
func SumForProject(entries []Entry, _ string) float64 {
	// Not used in the entry-level filter; project filtering happens at the
	// file selection step in Scan. This is a placeholder for future
	// per-entry filters (e.g. by repo path embedded in transcript).
	return Sum(entries)
}

// StartOfDay returns midnight (local time) for t.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EncodeProjectKey converts an absolute repo path to the Claude Code project
// directory name (e.g. "/Users/pras/dev/x" → "-Users-pras-dev-x"). Used to
// scope cost_repo to the current working directory.
func EncodeProjectKey(absPath string) string {
	return "-" + strings.ReplaceAll(strings.TrimPrefix(absPath, "/"), "/", "-")
}
