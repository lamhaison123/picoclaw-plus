// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package embedding

import (
	"context"
)

// Service defines the interface for text embedding generation
type Service interface {
	// Generate creates an embedding vector for the given text
	Generate(ctx context.Context, text string) ([]float32, error)

	// GenerateBatch creates embeddings for multiple texts
	GenerateBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Dimension returns the dimension of embeddings produced by this service
	Dimension() int

	// Close closes the embedding service connection
	Close() error
}

// Config holds embedding service configuration
type Config struct {
	Provider  string `json:"provider"`   // "openai", "local", "none"
	Model     string `json:"model"`      // e.g. "text-embedding-3-small"
	Dimension int    `json:"dimension"`  // e.g. 384, 1536
	APIKey    string `json:"api_key"`    // For API-based providers
	BaseURL   string `json:"base_url"`   // Optional custom endpoint
	TimeoutMs int    `json:"timeout_ms"` // Request timeout
}
