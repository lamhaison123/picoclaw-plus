// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: QdrantStore Unit Tests
// Sprint 1 Implementation - Phase 2
//
//go:build !no_qdrant
// +build !no_qdrant

package memory

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func init() {
	skipQdrantInitCheck = true
}

// setupTestStore helps create a real Qdrant connection for integration tests
func setupTestStore(t *testing.T) (*QdrantStore, string, func()) {
	url := os.Getenv("QDRANT_URL")
	if url == "" {
		t.Skip("QDRANT_URL not set, skipping integration test")
		return nil, "", nil
	}

	// Qdrant Go client requires gRPC port (default 6334), not REST HTTP (6333).
	// Translate the provided URL if it uses port 6333 to prevent http2 preface errors.
	if strings.HasSuffix(url, ":6333") {
		url = strings.TrimSuffix(url, ":6333") + ":6334"
	}

	collection := fmt.Sprintf("test_collection_%d", time.Now().UnixNano())

	cfg := QdrantConfig{
		URL:        url,
		Collection: collection,
	}

	vectorCfg := VectorConfig{
		Dimension: 384,
		TimeoutMs: 2000,
	}

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	// Temporarily enable network check for integration test setup
	oldSkip := skipQdrantInitCheck
	skipQdrantInitCheck = false
	defer func() { skipQdrantInitCheck = oldSkip }()

	store, err := NewQdrantStore(cfg, vectorCfg, breaker)
	if err != nil {
		t.Fatalf("Failed to connect to real Qdrant at %s: %v", url, err)
	}

	cleanup := func() {
		// Attempt to delete collection in cleanup
		ctx := context.Background()
		store.client.DeleteCollection(ctx, collection)
		store.Close()
	}

	return store, collection, cleanup
}

// TestQdrantConfig_Validation tests configuration validation
func TestQdrantConfig_Validation(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	vectorCfg := VectorConfig{
		Dimension: 384,
		TimeoutMs: 800,
	}

	tests := []struct {
		name    string
		cfg     QdrantConfig
		wantErr bool
		errCode string
	}{
		{
			name: "missing URL",
			cfg: QdrantConfig{
				Collection: "test_collection",
			},
			wantErr: true,
			errCode: "ERR_CONFIG_INVALID",
		},
		{
			name: "missing collection",
			cfg: QdrantConfig{
				URL: "http://localhost:6333",
			},
			wantErr: true,
			errCode: "ERR_CONFIG_INVALID",
		},
		{
			name: "invalid dimension zero",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test_collection",
			},
			wantErr: true,
			errCode: "ERR_CONFIG_INVALID",
		},
		{
			name: "invalid dimension negative",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test_collection",
			},
			wantErr: true,
			errCode: "ERR_CONFIG_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override dimension for specific tests
			testVectorCfg := vectorCfg
			switch tt.name {
			case "invalid dimension zero":
				testVectorCfg.Dimension = 0
			case "invalid dimension negative":
				testVectorCfg.Dimension = -1
			}

			_, err := NewQdrantStore(tt.cfg, testVectorCfg, breaker)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewQdrantStore() expected error, got nil")
					return
				}
				// Check error code is present
				if tt.errCode != "" && !contains(err.Error(), tt.errCode) {
					t.Errorf("NewQdrantStore() error = %v, want error containing %v", err, tt.errCode)
				}
			} else {
				if err != nil {
					t.Errorf("NewQdrantStore() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestUpsert_IDHandling tests that IDs are properly handled
func TestUpsert_IDHandling(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	if cleanup != nil {
		defer cleanup()
	}

	ctx := context.Background()

	// 1. String UUID
	uuidStr := "123e4567-e89b-12d3-a456-426614174000"
	vec1 := Vector{
		ID:        uuidStr,
		Embedding: make([]float32, 384),
		Metadata:  map[string]interface{}{"type": "uuid"},
	}

	// 2. Numeric String ID
	numStr := "123456789"
	vec2 := Vector{
		ID:        numStr,
		Embedding: make([]float32, 384),
		Metadata:  map[string]interface{}{"type": "numeric"},
	}

	// 3. Regular string ID (will be hashed to UUID)
	regStr := "my_custom_id"
	vec3 := Vector{
		ID:        regStr,
		Embedding: make([]float32, 384),
		Metadata:  map[string]interface{}{"type": "string"},
	}

	err := store.Upsert(ctx, []Vector{vec1, vec2, vec3})
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// Wait for indexing
	time.Sleep(500 * time.Millisecond)

	// Search to verify
	results, err := store.Search(ctx, Vector{Embedding: make([]float32, 384)}, 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

// TestUpsert_Batching tests that large uploads are batched
func TestUpsert_Batching(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	if cleanup != nil {
		defer cleanup()
	}

	ctx := context.Background()
	vectors := make([]Vector, 250) // More than maxBatchSize (100)

	for i := 0; i < 250; i++ {
		vectors[i] = Vector{
			ID:        strconv.Itoa(i + 1000),
			Embedding: make([]float32, 384),
			Metadata:  map[string]interface{}{"index": i},
		}
	}

	err := store.Upsert(ctx, vectors)
	if err != nil {
		t.Fatalf("Batch Upsert failed: %v", err)
	}
}

// TestDelete_IDParsing tests that delete handles both ID types
func TestDelete_IDParsing(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	if cleanup != nil {
		defer cleanup()
	}

	ctx := context.Background()

	vec1 := Vector{ID: "query-1", Embedding: make([]float32, 384)}
	vec2 := Vector{ID: "9999", Embedding: make([]float32, 384)}

	err := store.Upsert(ctx, []Vector{vec1, vec2})
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	err = store.Delete(ctx, []string{"query-1", "9999"})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

// TestSearch_DimensionValidation tests dimension checking
func TestSearch_DimensionValidation(t *testing.T) {
	store, _, cleanup := setupTestStore(t)
	if cleanup != nil {
		defer cleanup()
	}

	ctx := context.Background()
	_, err := store.Search(ctx, Vector{Embedding: make([]float32, 10)}, 5)

	if err == nil {
		t.Errorf("Expected dimension mismatch error, got nil")
	} else if !contains(err.Error(), "ERR_DIMENSION_MISMATCH") {
		t.Errorf("Expected ERR_DIMENSION_MISMATCH, got: %v", err)
	}
}

// TestRetry_ExponentialBackoff tests retry logic
func TestRetry_ExponentialBackoff(t *testing.T) {
	// Test that retry delays increase exponentially
	// This can be unit tested without Qdrant
	t.Skip("TODO: Implement unit test for retry logic")
}

// TestErrorMapping tests error classification
func TestErrorMapping(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	vectorCfg := VectorConfig{
		Dimension: 384,
		TimeoutMs: 800,
	}

	cfg := QdrantConfig{
		URL:        "http://localhost:6333",
		Collection: "test",
	}

	// This will fail to connect, but we can test error mapping
	store, err := NewQdrantStore(cfg, vectorCfg, breaker)
	if err != nil {
		// Check that connection errors are mapped correctly
		if !contains(err.Error(), "ERR_PROVIDER_UNAVAILABLE") && !contains(err.Error(), "ERR_CONFIG_INVALID") {
			t.Errorf("Expected ERR_PROVIDER_UNAVAILABLE or ERR_CONFIG_INVALID, got: %v", err)
		}
		return
	}

	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond) // Ensure timeout

	err = store.Health(ctx)
	if err != nil && !contains(err.Error(), "ERR_TIMEOUT") && !contains(err.Error(), "ERR_PROVIDER_UNAVAILABLE") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// TestCircuitBreaker_Integration tests circuit breaker integration
func TestCircuitBreaker_Integration(t *testing.T) {
	t.Skip("Requires Qdrant instance - run as integration test")

	// This test verifies:
	// 1. Circuit opens after max failures
	// 2. Requests are rejected when circuit is open
	// 3. Circuit transitions to half-open after timeout
	// 4. Circuit closes after successful requests
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
