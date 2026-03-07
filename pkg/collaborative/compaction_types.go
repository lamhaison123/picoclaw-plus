// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"sync"
	"time"
)

// CompactedContext holds compressed context information
type CompactedContext struct {
	Summary        string    // LLM-generated summary
	SummaryVersion int       // Version for tracking updates
	CompactedAt    time.Time // When compaction happened
	CompactedCount int       // Number of messages compacted
	OriginalSize   int       // Original size in bytes
	CompressedSize int       // Compressed size in bytes
	mu             sync.RWMutex
}

// CompactionConfig configures compaction behavior
type CompactionConfig struct {
	Enabled          bool          // Enable/disable compaction
	TriggerThreshold int           // Trigger at N messages (default: 40)
	KeepRecentCount  int           // Keep last N messages (default: 15)
	CompactBatchSize int           // Compact N messages at a time (default: 25)
	MinInterval      time.Duration // Min time between compactions (default: 5min)
	SummaryMaxLength int           // Max summary length (default: 2000 chars)

	// LLM settings
	LLMProvider   string        // "openai", "anthropic", "local"
	LLMModel      string        // Model to use (default: "gpt-4o-mini")
	LLMTimeout    time.Duration // Timeout for LLM calls (default: 30s)
	LLMMaxRetries int           // Max retry attempts (default: 3)
}

// CompactionMetrics tracks compaction performance
type CompactionMetrics struct {
	TotalCompactions   int64
	SuccessCount       int64
	FailureCount       int64
	TotalTimeSaved     time.Duration
	TotalBytesSaved    int64
	AverageCompression float64
	LastCompaction     time.Time
	mu                 sync.RWMutex
}

// CompactionRequest represents a compaction job
type CompactionRequest struct {
	SessionID       string
	Messages        []Message
	ExistingSummary string
	Config          CompactionConfig
	Timestamp       time.Time
}

// CompactionResult holds the result of compaction
type CompactionResult struct {
	Summary        string
	Success        bool
	Error          error
	OriginalSize   int
	CompressedSize int
	Duration       time.Duration
	MessagesCount  int
}
