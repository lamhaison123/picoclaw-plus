// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
//go:build cgo
// +build cgo

package memory

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLanceDBStore_Lifecycle(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lancedb-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})

	cfg := LanceDBConfig{
		Mode: "cgo",
		Path: tmpDir,
	}

	vectorCfg := VectorConfig{
		Dimension: 128,
		TimeoutMs: 5000,
	}

	store, err := NewLanceDBStore(cfg, vectorCfg, breaker)
	if err != nil {
		t.Fatalf("Failed to create LanceDB store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	// 1. Health check
	err = store.Health(ctx)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	// 2. Upsert vectors
	vectors := []Vector{
		{
			ID:        "vec1",
			Embedding: make([]float32, 128), // Dummy zero vector
			Metadata:  map[string]interface{}{"key": "value1"},
		},
		{
			ID:        "vec2",
			Embedding: make([]float32, 128),
			Metadata:  map[string]interface{}{"key": "value2"},
		},
	}
	// Add some distinct values for search
	vectors[0].Embedding[0] = 1.0

	err = store.Upsert(ctx, vectors)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// 3. Search
	query := Vector{
		Embedding: make([]float32, 128),
	}
	query.Embedding[0] = 1.0

	results, err := store.Search(ctx, query, 5)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("Expected results, got none")
	}

	// 4. Delete
	err = store.Delete(ctx, []string{"vec1"})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
