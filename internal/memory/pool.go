// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package memory

import (
	"strings"
	"sync"
)

var (
	// Buffer pool for general purpose byte slices (4KB default)
	bufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 4096)
			return &b
		},
	}

	// String builder pool for efficient string concatenation
	builderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	// Small buffer pool for short strings (512B)
	smallBufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 512)
			return &b
		},
	}
)

// GetBuffer returns a pooled byte buffer (4KB capacity)
func GetBuffer() *[]byte {
	return bufferPool.Get().(*[]byte)
}

// PutBuffer returns a buffer to the pool after clearing it
func PutBuffer(b *[]byte) {
	*b = (*b)[:0] // Clear but keep capacity
	bufferPool.Put(b)
}

// GetSmallBuffer returns a pooled small byte buffer (512B capacity)
func GetSmallBuffer() *[]byte {
	return smallBufferPool.Get().(*[]byte)
}

// PutSmallBuffer returns a small buffer to the pool after clearing it
func PutSmallBuffer(b *[]byte) {
	*b = (*b)[:0]
	smallBufferPool.Put(b)
}

// GetBuilder returns a pooled strings.Builder
func GetBuilder() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

// PutBuilder returns a builder to the pool after resetting it
func PutBuilder(sb *strings.Builder) {
	sb.Reset()
	builderPool.Put(sb)
}
