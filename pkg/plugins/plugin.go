// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package plugins

import (
	"context"
)

// Plugin represents a WASM plugin interface
type Plugin interface {
	// Name returns the plugin name
	Name() string
	
	// Version returns the plugin version
	Version() string
	
	// Initialize initializes the plugin
	Initialize(ctx context.Context) error
	
	// Execute executes the plugin with input data
	Execute(ctx context.Context, input []byte) ([]byte, error)
	
	// Cleanup cleans up plugin resources
	Cleanup(ctx context.Context) error
}

// PluginMetadata contains plugin metadata
type PluginMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Capabilities []string         `json:"capabilities"`
	Config      map[string]string `json:"config,omitempty"`
}

// PluginManager manages multiple plugins
type PluginManager struct {
	runtime *WASMRuntime
	plugins map[string]*PluginMetadata
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(ctx context.Context) (*PluginManager, error) {
	runtime, err := NewWASMRuntime(ctx)
	if err != nil {
		return nil, err
	}

	return &PluginManager{
		runtime: runtime,
		plugins: make(map[string]*PluginMetadata),
	}, nil
}

// LoadPlugin loads a plugin from WASM bytes
func (pm *PluginManager) LoadPlugin(ctx context.Context, metadata *PluginMetadata, wasmBytes []byte) error {
	if err := pm.runtime.LoadPlugin(ctx, metadata.Name, wasmBytes); err != nil {
		return err
	}

	pm.plugins[metadata.Name] = metadata
	return nil
}

// UnloadPlugin unloads a plugin
func (pm *PluginManager) UnloadPlugin(ctx context.Context, name string) error {
	if err := pm.runtime.UnloadPlugin(ctx, name); err != nil {
		return err
	}

	delete(pm.plugins, name)
	return nil
}

// GetPlugin returns plugin metadata
func (pm *PluginManager) GetPlugin(name string) (*PluginMetadata, bool) {
	meta, exists := pm.plugins[name]
	return meta, exists
}

// ListPlugins returns all loaded plugins
func (pm *PluginManager) ListPlugins() []*PluginMetadata {
	list := make([]*PluginMetadata, 0, len(pm.plugins))
	for _, meta := range pm.plugins {
		list = append(list, meta)
	}
	return list
}

// Close closes the plugin manager
func (pm *PluginManager) Close(ctx context.Context) error {
	return pm.runtime.Close(ctx)
}
