// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/providers/protocoltypes"
)

// LLMSummarizer uses LLM to generate summaries
type LLMSummarizer struct {
	config   CompactionConfig
	provider providers.LLMProvider
}

// NewLLMSummarizer creates a new LLM summarizer
func NewLLMSummarizer(config CompactionConfig, provider providers.LLMProvider) *LLMSummarizer {
	return &LLMSummarizer{
		config:   config,
		provider: provider,
	}
}

// Summarize generates a summary using LLM
func (s *LLMSummarizer) Summarize(ctx context.Context, req *CompactionRequest) (*CompactionResult, error) {
	startTime := time.Now()

	logger.DebugCF("compaction", "Starting summarization", map[string]any{
		"session_id":     req.SessionID,
		"message_count":  len(req.Messages),
		"has_existing":   req.ExistingSummary != "",
		"existing_chars": len(req.ExistingSummary),
	})

	// Build prompt
	prompt := s.buildPrompt(req)

	// Call LLM with retry
	summary, err := s.callLLMWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Validate and truncate if needed
	summary = s.validateAndTruncate(summary)

	// Calculate sizes
	originalSize := s.calculateSize(req.Messages)
	compressedSize := len(summary)

	logger.InfoCF("compaction", "Summarization completed", map[string]any{
		"session_id":        req.SessionID,
		"original_size":     originalSize,
		"compressed_size":   compressedSize,
		"compression_ratio": float64(originalSize) / float64(compressedSize),
		"duration":          time.Since(startTime).String(),
	})

	return &CompactionResult{
		Summary:        summary,
		Success:        true,
		Error:          nil,
		OriginalSize:   originalSize,
		CompressedSize: compressedSize,
		Duration:       time.Since(startTime),
		MessagesCount:  len(req.Messages),
	}, nil
}

// buildPrompt builds the summarization prompt
func (s *LLMSummarizer) buildPrompt(req *CompactionRequest) string {
	var sb strings.Builder

	sb.WriteString("You are a context summarizer for a collaborative AI team.\n\n")

	if req.ExistingSummary != "" {
		sb.WriteString("=== Previous Summary ===\n")
		sb.WriteString(req.ExistingSummary)
		sb.WriteString("\n\n")
	}

	sb.WriteString("=== New Messages to Summarize ===\n")
	for _, msg := range req.Messages {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			msg.Timestamp.Format("15:04:05"),
			strings.ToUpper(msg.Role),
			msg.Content,
		))
	}

	sb.WriteString("\n=== Task ===\n")
	sb.WriteString("Summarize the conversation above, preserving:\n")
	sb.WriteString("1. Project goals and objectives\n")
	sb.WriteString("2. Key technical decisions (architecture, tech stack, patterns)\n")
	sb.WriteString("3. Important constraints and requirements\n")
	sb.WriteString("4. Action items and assignments\n")
	sb.WriteString("5. Critical issues or blockers\n\n")

	sb.WriteString("Format as structured summary:\n")
	sb.WriteString("## Project Overview\n")
	sb.WriteString("[Brief description]\n\n")
	sb.WriteString("## Key Decisions\n")
	sb.WriteString("- Decision 1\n")
	sb.WriteString("- Decision 2\n\n")
	sb.WriteString("## Architecture\n")
	sb.WriteString("[Architecture overview]\n\n")
	sb.WriteString("## Action Items\n")
	sb.WriteString("- Item 1 (@role)\n")
	sb.WriteString("- Item 2 (@role)\n\n")

	sb.WriteString(fmt.Sprintf("Keep summary concise but comprehensive. Max %d characters.\n", req.Config.SummaryMaxLength))

	return sb.String()
}

// callLLMWithRetry calls LLM with exponential backoff retry
func (s *LLMSummarizer) callLLMWithRetry(ctx context.Context, prompt string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= s.config.LLMMaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			logger.DebugCF("compaction", "Retrying LLM call", map[string]any{
				"attempt": attempt + 1,
				"backoff": backoff.String(),
			})
			time.Sleep(backoff)
		}

		summary, err := s.callLLM(ctx, prompt)
		if err == nil {
			return summary, nil
		}

		lastErr = err
		logger.WarnCF("compaction", "LLM call failed, retrying", map[string]any{
			"attempt": attempt + 1,
			"error":   err.Error(),
		})
	}

	return "", fmt.Errorf("all retry attempts failed: %w", lastErr)
}

// callLLM makes the actual LLM API call
func (s *LLMSummarizer) callLLM(ctx context.Context, prompt string) (string, error) {
	// Create timeout context with fallback to 120s if not configured
	timeout := s.config.LLMTimeout
	if timeout == 0 {
		timeout = 120 * time.Second // Increased from 90s to 120s per architecture review
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // CRITICAL: Prevent context leak

	// Build messages
	messages := []protocoltypes.Message{
		{
			Role:    "system",
			Content: "You are a helpful context summarizer. Generate concise, structured summaries that preserve key information.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// Determine model
	model := s.config.LLMModel
	if model == "" {
		model = s.provider.GetDefaultModel()
	}

	// Call provider
	resp, err := s.provider.Chat(timeoutCtx, messages, nil, model, map[string]any{
		"temperature": 0.3, // Low temperature for consistent summaries
		"max_tokens":  1000,
	})
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// validateAndTruncate validates and truncates summary if needed
func (s *LLMSummarizer) validateAndTruncate(summary string) string {
	summary = strings.TrimSpace(summary)

	if len(summary) > s.config.SummaryMaxLength {
		summary = summary[:s.config.SummaryMaxLength]
		// Try to cut at sentence boundary
		if idx := strings.LastIndex(summary, "."); idx > s.config.SummaryMaxLength-100 {
			summary = summary[:idx+1]
		}
	}

	return summary
}

// calculateSize calculates total size of messages
func (s *LLMSummarizer) calculateSize(messages []Message) int {
	total := 0
	for _, msg := range messages {
		total += len(msg.Content) + len(msg.Role) + 50 // +50 for metadata
	}
	return total
}
