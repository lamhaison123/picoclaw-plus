// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sipeed/picoclaw/pkg/plugins"
)

func main() {
	ctx := context.Background()

	// Create plugin manager
	pm, err := plugins.NewPluginManager(ctx)
	if err != nil {
		log.Fatalf("Failed to create plugin manager: %v", err)
	}
	defer pm.Close(ctx)

	// Check if WASM file provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: plugin_example <path-to-wasm-file>")
	}

	wasmPath := os.Args[1]

	// Read WASM file
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		log.Fatalf("Failed to read WASM file: %v", err)
	}

	// Plugin metadata
	metadata := &plugins.PluginMetadata{
		Name:        "math_plugin",
		Version:     "1.0.0",
		Description: "Simple math operations plugin",
		Author:      "PicoClaw Team",
		Capabilities: []string{"add", "subtract", "multiply", "divide", "factorial", "fibonacci", "power"},
	}

	// Load plugin
	fmt.Printf("Loading plugin: %s v%s\n", metadata.Name, metadata.Version)
	err = pm.LoadPlugin(ctx, metadata, wasmBytes)
	if err != nil {
		log.Fatalf("Failed to load plugin: %v", err)
	}

	fmt.Printf("✅ Plugin loaded successfully\n\n")

	// Test add function
	fmt.Println("Testing add(5, 7)...")
	result, err := pm.runtime.CallFunction(ctx, "math_plugin", "add", 5, 7)
	if err != nil {
		log.Fatalf("Failed to call add: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test subtract function
	fmt.Println("Testing subtract(10, 3)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "subtract", 10, 3)
	if err != nil {
		log.Fatalf("Failed to call subtract: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test multiply function
	fmt.Println("Testing multiply(6, 7)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "multiply", 6, 7)
	if err != nil {
		log.Fatalf("Failed to call multiply: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test divide function
	fmt.Println("Testing divide(20, 4)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "divide", 20, 4)
	if err != nil {
		log.Fatalf("Failed to call divide: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test factorial function
	fmt.Println("Testing factorial(5)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "factorial", 5)
	if err != nil {
		log.Fatalf("Failed to call factorial: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test fibonacci function
	fmt.Println("Testing fibonacci(10)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "fibonacci", 10)
	if err != nil {
		log.Fatalf("Failed to call fibonacci: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// Test power function
	fmt.Println("Testing power(2, 8)...")
	result, err = pm.runtime.CallFunction(ctx, "math_plugin", "power", 2, 8)
	if err != nil {
		log.Fatalf("Failed to call power: %v", err)
	}
	fmt.Printf("Result: %d\n\n", result[0])

	// List all loaded plugins
	fmt.Println("Loaded plugins:")
	for _, p := range pm.ListPlugins() {
		fmt.Printf("  - %s v%s: %s\n", p.Name, p.Version, p.Description)
		fmt.Printf("    Capabilities: %v\n", p.Capabilities)
	}

	// Unload plugin
	fmt.Println("\nUnloading plugin...")
	err = pm.UnloadPlugin(ctx, "math_plugin")
	if err != nil {
		log.Fatalf("Failed to unload plugin: %v", err)
	}

	fmt.Println("✅ Plugin unloaded successfully")
}
