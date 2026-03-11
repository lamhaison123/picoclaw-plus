// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
//go:build !cgo
// +build !cgo

package agent

import (
	"fmt"
	"github.com/sipeed/picoclaw/pkg/config"
	memory "github.com/sipeed/picoclaw/pkg/memory/vector"
)

// createLanceDBStore returns an error when built without LanceDB/CGO support
func createLanceDBStore(cfg config.VectorStoreConfig, dimension int) (memory.VectorStore, error) {
	return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: LanceDB requires CGO, which is disabled in this build")
}
