// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
//go:build cgo
// +build cgo

package agent

import (
	"github.com/sipeed/picoclaw/pkg/config"
	memory "github.com/sipeed/picoclaw/pkg/memory/vector"
)

// createLanceDBStore creates a LanceDB vector store when CGO support is available
func createLanceDBStore(cfg config.VectorStoreConfig, dimension int) (memory.VectorStore, error) {
	breaker := memory.NewCircuitBreaker(memory.CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})

	return memory.NewLanceDBStore(memory.LanceDBConfig{
		Mode: cfg.LanceDB.Mode,
		Path: cfg.LanceDB.Path,
		URL:  cfg.LanceDB.URL,
	}, memory.VectorConfig{
		Dimension: dimension,
		TimeoutMs: cfg.TimeoutMs,
	}, breaker)
}
