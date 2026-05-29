package cost

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
)

// fileCacheEntry is what we persist per source transcript file: the source's
// mtime + its parsed assistant entries. We re-parse only when the source's
// mtime has shifted.
type fileCacheEntry struct {
	SourceMTimeUnixNano int64
	Entries             []Entry
}

func cacheDir() (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(h, ".cache")
	}
	d := filepath.Join(base, "goccline", "cost")
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	return d, nil
}

// pathKey flattens a source path into a single safe filename. Multiple
// transcript files share a parent dir, so we keep the project segment +
// uuid filename via replacement, no hashing — debuggability over density.
func pathKey(srcPath string) string {
	s := strings.TrimPrefix(srcPath, string(os.PathSeparator))
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	return s
}

// loadCached returns cached entries if a cache exists AND its recorded
// source mtime matches the current source mtime.
func loadCached(srcPath string, srcMTime int64) ([]Entry, bool) {
	dir, err := cacheDir()
	if err != nil {
		return nil, false
	}
	f, err := os.Open(filepath.Join(dir, pathKey(srcPath)+".gob"))
	if err != nil {
		return nil, false
	}
	defer f.Close()
	var ce fileCacheEntry
	if err := gob.NewDecoder(f).Decode(&ce); err != nil {
		return nil, false
	}
	if ce.SourceMTimeUnixNano != srcMTime {
		return nil, false
	}
	return ce.Entries, true
}

// storeCached writes parsed entries atomically (tmp + rename) so concurrent
// renders don't see a half-written file.
func storeCached(srcPath string, srcMTime int64, entries []Entry) {
	dir, err := cacheDir()
	if err != nil {
		return
	}
	final := filepath.Join(dir, pathKey(srcPath)+".gob")
	tmp := final + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return
	}
	if err := gob.NewEncoder(f).Encode(fileCacheEntry{
		SourceMTimeUnixNano: srcMTime,
		Entries:             entries,
	}); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return
	}
	_ = os.Rename(tmp, final)
}

// parseFileCached is the cache-fronted equivalent of parseFile. Stats the
// source for its mtime, returns cached entries on hit, parses + writes
// otherwise. Errors are swallowed (best-effort cache) so a render still
// succeeds when the cache dir is unwritable.
func parseFileCached(path string) []Entry {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	mtime := info.ModTime().UnixNano()
	if cached, ok := loadCached(path, mtime); ok {
		return cached
	}
	entries := parseFile(path)
	storeCached(path, mtime, entries)
	return entries
}
