// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package embedding

import (
	"context"
	"fmt"
)

// NullService is a no-op embedding service used when embeddings are disabled
type NullService struct {
	dimension int
}

// NewNullService creates a null embedding service
func NewNullService(dimension int) *NullService {
	if dimension <= 0 {
		dimension = 384
	}
	return &NullService{dimension: dimension}
}

// Generate returns an error indicating embeddings are disabled
func (s *NullService) Generate(ctx context.Context, text string) ([]float32, error) {
	return nil, fmt.Errorf("embedding service is disabled")
}

// GenerateBatch returns an error indicating embeddings are disabled
func (s *NullService) GenerateBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, fmt.Errorf("embedding service is disabled")
}

// Dimension returns the configured dimension
func (s *NullService) Dimension() int {
	return s.dimension
}

// Close is a no-op
func (s *NullService) Close() error {
	return nil
}
