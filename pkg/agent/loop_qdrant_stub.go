// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build no_qdrant
// +build no_qdrant

package agent

import (
	"fmt"

	"github.com/sipeed/picoclaw/pkg/config"
	memory "github.com/sipeed/picoclaw/pkg/memory/vector"
)

// createQdrantStore returns an error when built without Qdrant support
func createQdrantStore(cfg config.VectorStoreConfig, dimension int) (memory.VectorStore, error) {
	return nil, fmt.Errorf("Qdrant support not available (built with no_qdrant tag)")
}
