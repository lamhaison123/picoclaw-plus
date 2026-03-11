// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build !no_qdrant
// +build !no_qdrant

package agent

import (
	"github.com/sipeed/picoclaw/pkg/config"
	memory "github.com/sipeed/picoclaw/pkg/memory/vector"
)

// createQdrantStore creates a Qdrant vector store when Qdrant support is available
func createQdrantStore(cfg config.VectorStoreConfig, dimension int) (memory.VectorStore, error) {
	breaker := memory.NewCircuitBreaker(memory.CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	vectorCfg := memory.VectorConfig{
		Dimension: dimension,
		TimeoutMs: cfg.Qdrant.TimeoutMS,
	}

	qdrantCfg := memory.QdrantConfig{
		URL:        cfg.Qdrant.URL,
		Collection: cfg.Qdrant.Collection,
		APIKey:     cfg.Qdrant.APIKey,
	}

	return memory.NewQdrantStore(qdrantCfg, vectorCfg, breaker)
}
