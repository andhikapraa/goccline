// Package components defines the contract every statusline widget implements,
// plus the registry that the render loop dispatches against.
package components

import (
	"sync"

	"github.com/andhikapraa/goccline/internal/config"
	"github.com/andhikapraa/goccline/internal/input"
)

// Context is the per-render state passed to every component. Components
// MUST NOT mutate it.
type Context struct {
	Input  input.Payload
	Config *config.Config
}

// RenderFn returns the string this component contributes to the line.
// An empty string means "nothing to show" — the render loop will omit it.
type RenderFn func(ctx *Context) string

var (
	registryMu sync.RWMutex
	registry   = map[string]RenderFn{}
)

// Register adds a component to the global registry. Call from init() in each
// component file.
func Register(name string, fn RenderFn) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = fn
}

// Lookup returns the renderer for name, or nil if unregistered.
func Lookup(name string) RenderFn {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry[name]
}

// Names returns the registered component names (for diagnostics).
func Names() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}
