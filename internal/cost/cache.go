package cost

import (
	"sync"
	"time"
)

// MemoCache caches the result of Scan within a single goccline invocation
// so multiple cost components don't re-parse the same files. Not concurrent
// safe across goroutines; intended for the sequential render loop.
type MemoCache struct {
	mu      sync.Mutex
	scans   map[memoKey][]Entry
}

type memoKey struct {
	root  string
	since int64 // unix seconds; 0 == no filter
}

func NewMemo() *MemoCache {
	return &MemoCache{scans: map[memoKey][]Entry{}}
}

// Scan returns cached entries for (root, since), or computes and stores.
// When a caller asks for a narrower `since` than a previously cached
// broader scan of the same root, we filter the broader result by mtime
// equivalent — but to keep this simple, the broadest scan wins: callers
// should ask for the widest window first.
func (m *MemoCache) Scan(root string, since time.Time) []Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := memoKey{root: root, since: since.Unix()}
	if v, ok := m.scans[k]; ok {
		return v
	}
	// If we already scanned the same root with an earlier `since`, reuse it.
	// (We just return the entries; per-entry timestamp filtering happens
	// at SumSince.)
	for existing, entries := range m.scans {
		if existing.root != root {
			continue
		}
		if since.IsZero() && existing.since == 0 {
			m.scans[k] = entries
			return entries
		}
		// A previous scan was broader (earlier `since`) — reuse it.
		if existing.since != 0 && existing.since <= since.Unix() {
			m.scans[k] = entries
			return entries
		}
	}
	entries := Scan(root, since)
	m.scans[k] = entries
	return entries
}
