// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package plugins

import (
	"context"
	"testing"
)

// TestWASMRuntimeCreation tests basic runtime creation
func TestWASMRuntimeCreation(t *testing.T) {
	ctx := context.Background()
	
	runtime, err := NewWASMRuntime(ctx)
	if err != nil {
		t.Fatalf("Failed to create WASM runtime: %v", err)
	}
	defer runtime.Close(ctx)

	if runtime == nil {
		t.Fatal("Runtime should not be nil")
	}

	if runtime.runtime == nil {
		t.Fatal("Internal runtime should not be nil")
	}
}

// TestPluginManagerCreation tests plugin manager creation
func TestPluginManagerCreation(t *testing.T) {
	ctx := context.Background()
	
	pm, err := NewPluginManager(ctx)
	if err != nil {
		t.Fatalf("Failed to create plugin manager: %v", err)
	}
	defer pm.Close(ctx)

	if pm == nil {
		t.Fatal("Plugin manager should not be nil")
	}

	plugins := pm.ListPlugins()
	if len(plugins) != 0 {
		t.Fatalf("Expected 0 plugins, got %d", len(plugins))
	}
}

// TestSimpleWASMPlugin tests loading a simple WASM module
func TestSimpleWASMPlugin(t *testing.T) {
	ctx := context.Background()
	
	runtime, err := NewWASMRuntime(ctx)
	if err != nil {
		t.Fatalf("Failed to create WASM runtime: %v", err)
	}
	defer runtime.Close(ctx)

	// Simple WASM module that exports an "add" function
	// (func (export "add") (param i32 i32) (result i32)
	//   local.get 0
	//   local.get 1
	//   i32.add)
	wasmBytes := []byte{
		0x00, 0x61, 0x73, 0x6d, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // WASM version
		0x01, 0x07, // Type section
		0x01, // 1 type
		0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f, // (i32, i32) -> i32
		0x03, 0x02, // Function section
		0x01, 0x00, // 1 function, type 0
		0x07, 0x07, // Export section
		0x01, // 1 export
		0x03, 0x61, 0x64, 0x64, // "add"
		0x00, 0x00, // function 0
		0x0a, 0x09, // Code section
		0x01, // 1 function body
		0x07, // body size
		0x00, // 0 locals
		0x20, 0x00, // local.get 0
		0x20, 0x01, // local.get 1
		0x6a, // i32.add
		0x0b, // end
	}

	err = runtime.LoadPlugin(ctx, "test_add", wasmBytes)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Test calling the add function
	result, err := runtime.CallFunction(ctx, "test_add", "add", 5, 7)
	if err != nil {
		t.Fatalf("Failed to call function: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(result))
	}

	if result[0] != 12 {
		t.Fatalf("Expected 12, got %d", result[0])
	}

	// Test unloading
	err = runtime.UnloadPlugin(ctx, "test_add")
	if err != nil {
		t.Fatalf("Failed to unload plugin: %v", err)
	}

	// Verify plugin is unloaded
	plugins := runtime.ListPlugins()
	if len(plugins) != 0 {
		t.Fatalf("Expected 0 plugins after unload, got %d", len(plugins))
	}
}

// TestPluginManagerWithMetadata tests plugin manager with metadata
func TestPluginManagerWithMetadata(t *testing.T) {
	ctx := context.Background()
	
	pm, err := NewPluginManager(ctx)
	if err != nil {
		t.Fatalf("Failed to create plugin manager: %v", err)
	}
	defer pm.Close(ctx)

	metadata := &PluginMetadata{
		Name:        "test_plugin",
		Version:     "1.0.0",
		Description: "Test plugin",
		Author:      "PicoClaw Team",
		Capabilities: []string{"add", "subtract"},
	}

	// Simple add function WASM
	wasmBytes := []byte{
		0x00, 0x61, 0x73, 0x6d,
		0x01, 0x00, 0x00, 0x00,
		0x01, 0x07,
		0x01,
		0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f,
		0x03, 0x02,
		0x01, 0x00,
		0x07, 0x07,
		0x01,
		0x03, 0x61, 0x64, 0x64,
		0x00, 0x00,
		0x0a, 0x09,
		0x01,
		0x07,
		0x00,
		0x20, 0x00,
		0x20, 0x01,
		0x6a,
		0x0b,
	}

	err = pm.LoadPlugin(ctx, metadata, wasmBytes)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Verify plugin is loaded
	loadedMeta, exists := pm.GetPlugin("test_plugin")
	if !exists {
		t.Fatal("Plugin should exist")
	}

	if loadedMeta.Name != metadata.Name {
		t.Fatalf("Expected name %s, got %s", metadata.Name, loadedMeta.Name)
	}

	if loadedMeta.Version != metadata.Version {
		t.Fatalf("Expected version %s, got %s", metadata.Version, loadedMeta.Version)
	}

	// List plugins
	plugins := pm.ListPlugins()
	if len(plugins) != 1 {
		t.Fatalf("Expected 1 plugin, got %d", len(plugins))
	}

	// Unload plugin
	err = pm.UnloadPlugin(ctx, "test_plugin")
	if err != nil {
		t.Fatalf("Failed to unload plugin: %v", err)
	}

	// Verify unloaded
	_, exists = pm.GetPlugin("test_plugin")
	if exists {
		t.Fatal("Plugin should not exist after unload")
	}
}
