// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package plugins

import (
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// WASMRuntime manages WASM plugin execution
type WASMRuntime struct {
	runtime wazero.Runtime
	mu      sync.RWMutex
	modules map[string]api.Module
}

// NewWASMRuntime creates a new WASM runtime
func NewWASMRuntime(ctx context.Context) (*WASMRuntime, error) {
	// Create runtime with default config
	r := wazero.NewRuntime(ctx)

	// Instantiate WASI for basic I/O support
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		r.Close(ctx)
		return nil, fmt.Errorf("failed to instantiate WASI: %w", err)
	}

	return &WASMRuntime{
		runtime: r,
		modules: make(map[string]api.Module),
	}, nil
}

// LoadPlugin loads a WASM plugin from bytes
func (w *WASMRuntime) LoadPlugin(ctx context.Context, name string, wasmBytes []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if already loaded
	if _, exists := w.modules[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Compile and instantiate module
	mod, err := w.runtime.InstantiateWithConfig(ctx, wasmBytes,
		wazero.NewModuleConfig().WithName(name))
	if err != nil {
		return fmt.Errorf("failed to instantiate plugin %s: %w", name, err)
	}

	w.modules[name] = mod
	return nil
}

// CallFunction calls an exported function from a loaded plugin
func (w *WASMRuntime) CallFunction(ctx context.Context, pluginName, funcName string, args ...uint64) ([]uint64, error) {
	w.mu.RLock()
	mod, exists := w.modules[pluginName]
	w.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not loaded", pluginName)
	}

	fn := mod.ExportedFunction(funcName)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in plugin %s", funcName, pluginName)
	}

	return fn.Call(ctx, args...)
}

// UnloadPlugin unloads a plugin
func (w *WASMRuntime) UnloadPlugin(ctx context.Context, name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	mod, exists := w.modules[name]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	if err := mod.Close(ctx); err != nil {
		return fmt.Errorf("failed to close plugin %s: %w", name, err)
	}

	delete(w.modules, name)
	return nil
}

// Close closes the runtime and all loaded modules
func (w *WASMRuntime) Close(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	for name, mod := range w.modules {
		if err := mod.Close(ctx); err != nil {
			return fmt.Errorf("failed to close plugin %s: %w", name, err)
		}
	}

	return w.runtime.Close(ctx)
}

// ListPlugins returns names of all loaded plugins
func (w *WASMRuntime) ListPlugins() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	names := make([]string, 0, len(w.modules))
	for name := range w.modules {
		names = append(names, name)
	}
	return names
}
