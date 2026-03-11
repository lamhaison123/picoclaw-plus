// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// package memory/vector
// Module: LanceDBStore Stub - Used when CGO is disabled
//
//go:build !cgo
// +build !cgo

package memory

import (
	"context"
	"fmt"
)

// LanceDBStore Stub implementation
type LanceDBStore struct{}

// NewLanceDBStore creates a stub instance when LanceDB is not compiled
func NewLanceDBStore(cfg LanceDBConfig, vectorCfg VectorConfig, breaker *CircuitBreaker) (*LanceDBStore, error) {
	return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: LanceDB requires CGO, which is disabled in this build")
}

func (l *LanceDBStore) Upsert(ctx context.Context, vectors []Vector) error {
	return nil
}

func (l *LanceDBStore) Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error) {
	return nil, nil
}

func (l *LanceDBStore) Delete(ctx context.Context, ids []string) error {
	return nil
}

func (l *LanceDBStore) Health(ctx context.Context) error {
	return fmt.Errorf("LanceDB disabled")
}

func (l *LanceDBStore) Close() error {
	return nil
}
