// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Core Interfaces
// Sprint 1 Implementation - Phase 1

package memory

import (
	"context"
)

// Vector represents an embedding vector with metadata
type Vector struct {
	ID        string
	Embedding []float32
	Metadata  map[string]interface{}
}

// SearchResult represents a vector search result
type SearchResult struct {
	ID       string
	Score    float32
	Metadata map[string]interface{}
}

// VectorStore defines the interface for vector database operations
// Implementations: QdrantStore, LanceDBStore
type VectorStore interface {
	// Upsert inserts or updates vectors in the store
	Upsert(ctx context.Context, vectors []Vector) error

	// Search performs vector similarity search
	Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error)

	// Delete removes vectors from the store by IDs
	Delete(ctx context.Context, ids []string) error

	// Health checks the health of the vector store connection
	Health(ctx context.Context) error

	// Close closes the vector store connection
	Close() error
}

// MemoryProvider defines the interface for memory provider operations
// Implementations: Mem0Client, MindGraphClient, SidecarClient
type MemoryProvider interface {
	// Store stores a memory entry
	Store(ctx context.Context, content string, metadata map[string]interface{}) (string, error)

	// Recall retrieves relevant memories based on query
	Recall(ctx context.Context, query string, limit int) ([]Memory, error)

	// Update updates an existing memory entry
	Update(ctx context.Context, id string, content string, metadata map[string]interface{}) error

	// Delete removes a memory entry by ID
	Delete(ctx context.Context, id string) error

	// Health checks the health of the memory provider connection
	Health(ctx context.Context) error

	// Close closes the memory provider connection
	Close() error
}

// Memory represents a memory entry from the provider
type Memory struct {
	ID        string
	Content   string
	Metadata  map[string]interface{}
	Score     float32
	Timestamp int64
}

// Canonical Error Codes (aligned with @architect specs)
const (
	ErrConfigInvalid       = "ERR_CONFIG_INVALID"
	ErrProviderUnavailable = "ERR_PROVIDER_UNAVAILABLE"
	ErrTimeout             = "ERR_TIMEOUT"
	ErrCircuitOpen         = "ERR_CIRCUIT_OPEN"
	ErrSchemaMismatch      = "ERR_SCHEMA_MISMATCH"
	ErrDimensionMismatch   = "ERR_DIMENSION_MISMATCH"
	ErrAuthFailed          = "ERR_AUTH_FAILED"
	ErrRateLimited         = "ERR_RATE_LIMITED"
	ErrInternal            = "ERR_INTERNAL"
)
