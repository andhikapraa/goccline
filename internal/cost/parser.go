package cost

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Usage is the per-message token breakdown we care about for cost.
type Usage struct {
	InputTokens     int
	OutputTokens    int
	CacheReadTokens int
	CacheWrite5m    int
	CacheWrite1h    int
}

// Entry is one observation: dedup key + model + usage + when.
type Entry struct {
	RequestID string
	UUID      string
	Model     string
	Timestamp time.Time
	Usage     Usage
}

// rawEntry mirrors the relevant subset of a Claude transcript JSONL line.
type rawEntry struct {
	Type      string `json:"type"`
	UUID      string `json:"uuid"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Message   struct {
		Model string `json:"model"`
		Usage struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheCreation            struct {
				Ephemeral5m int `json:"ephemeral_5m_input_tokens"`
				Ephemeral1h int `json:"ephemeral_1h_input_tokens"`
			} `json:"cache_creation"`
		} `json:"usage"`
	} `json:"message"`
}

// parseFile streams one transcript JSONL file and emits assistant entries.
// Errors on a single line are skipped; truncated files yield whatever we
// could read.
func parseFile(path string) []Entry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	var out []Entry
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var r rawEntry
		if err := json.Unmarshal(line, &r); err != nil {
			continue
		}
		if r.Type != "assistant" {
			continue
		}
		// Skip records that have no usage info at all (tool-only messages).
		u := r.Message.Usage
		if u.InputTokens == 0 && u.OutputTokens == 0 && u.CacheReadInputTokens == 0 && u.CacheCreationInputTokens == 0 {
			continue
		}
		ts, _ := time.Parse(time.RFC3339, r.Timestamp)
		out = append(out, Entry{
			RequestID: r.RequestID,
			UUID:      r.UUID,
			Model:     r.Message.Model,
			Timestamp: ts,
			Usage: Usage{
				InputTokens:     u.InputTokens,
				OutputTokens:    u.OutputTokens,
				CacheReadTokens: u.CacheReadInputTokens,
				CacheWrite5m:    u.CacheCreation.Ephemeral5m,
				CacheWrite1h:    u.CacheCreation.Ephemeral1h,
			},
		})
	}
	return out
}

// Scan walks rootDir, finds *.jsonl files modified since `since`, and
// parses them in parallel. Returns deduplicated entries (by RequestID+UUID).
// `since` zero-value means "no time filter".
func Scan(rootDir string, since time.Time) []Entry {
	var files []string
	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".jsonl" {
			return nil
		}
		if !since.IsZero() {
			info, _ := d.Info()
			if info != nil && info.ModTime().Before(since) {
				return nil
			}
		}
		files = append(files, path)
		return nil
	})
	return parseParallel(files)
}

func parseParallel(files []string) []Entry {
	if len(files) == 0 {
		return nil
	}
	workers := runtime.NumCPU()
	if workers > len(files) {
		workers = len(files)
	}
	ch := make(chan string, len(files))
	results := make(chan []Entry, len(files))
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range ch {
				results <- parseFile(path)
			}
		}()
	}
	for _, f := range files {
		ch <- f
	}
	close(ch)
	go func() { wg.Wait(); close(results) }()

	// Dedup by RequestID. A single API call can be written to the JSONL
	// multiple times with different UUIDs as Claude Code iterates / retries,
	// so UUID alone (or composite with RequestID) over-counts.
	// Entries without a RequestID fall back to UUID so we don't collapse
	// non-request rows (system/tool-only messages) into a single bucket.
	seen := make(map[string]struct{}, 1024)
	var out []Entry
	for batch := range results {
		for _, e := range batch {
			key := e.RequestID
			if key == "" {
				key = "uuid:" + e.UUID
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, e)
		}
	}
	return out
}
