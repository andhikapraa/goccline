// Package cache is a small on-disk JSON cache for outputs of slow operations
// (e.g. prayer-times API). Files live under $XDG_CACHE_HOME/goccline.
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// dir returns the goccline cache directory, creating it if needed.
func dir() (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".cache")
	}
	d := filepath.Join(base, "goccline")
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	return d, nil
}

// envelope wraps cached payloads with their expiry so a reader can decide
// whether the value is still good.
type envelope struct {
	ExpiresAt time.Time       `json:"expires_at"`
	Value     json.RawMessage `json:"value"`
}

// Get loads a cached value into v if the file exists and hasn't expired.
// Returns ok=true only on hit.
func Get(key string, v any) (ok bool) {
	d, err := dir()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(d, key+".json"))
	if err != nil {
		return false
	}
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return false
	}
	if time.Now().After(env.ExpiresAt) {
		return false
	}
	return json.Unmarshal(env.Value, v) == nil
}

// Set stores v at key with the given TTL.
func Set(key string, ttl time.Duration, v any) error {
	d, err := dir()
	if err != nil {
		return err
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	env := envelope{ExpiresAt: time.Now().Add(ttl), Value: raw}
	out, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, key+".json"), out, 0o644)
}
