// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIService implements embedding generation using OpenAI-compatible API
type OpenAIService struct {
	apiKey    string
	baseURL   string
	model     string
	dimension int
	timeout   time.Duration
	client    *http.Client
}

// NewOpenAIService creates a new OpenAI embedding service
func NewOpenAIService(cfg Config) (*OpenAIService, error) {
	apiKey := cfg.APIKey
	if apiKey == "" {
		// Provide a dummy key for local endpoints (vllm, ollama, custom api) that don't require authentication
		apiKey = "sk-dummy"
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := cfg.Model
	if model == "" {
		model = "text-embedding-3-small"
	}

	dimension := cfg.Dimension
	if dimension <= 0 {
		dimension = 384 // Default for text-embedding-3-small with dimension reduction
	}

	timeout := time.Duration(cfg.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &OpenAIService{
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		dimension: dimension,
		timeout:   timeout,
		client:    &http.Client{Timeout: timeout},
	}, nil
}

// Generate creates an embedding for a single text
func (s *OpenAIService) Generate(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := s.GenerateBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return embeddings[0], nil
}

// GenerateBatch creates embeddings for multiple texts
func (s *OpenAIService) GenerateBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	// Prepare request
	reqBody := map[string]interface{}{
		"input": texts,
		"model": s.model,
	}

	// Add dimension parameter for models that support it
	if s.model == "text-embedding-3-small" || s.model == "text-embedding-3-large" {
		reqBody["dimensions"] = s.dimension
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"data"`
		Model string `json:"model"`
		Usage struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract embeddings
	embeddings := make([][]float32, len(result.Data))
	for _, item := range result.Data {
		if item.Index >= len(embeddings) {
			return nil, fmt.Errorf("invalid embedding index: %d", item.Index)
		}
		embeddings[item.Index] = item.Embedding
	}

	return embeddings, nil
}

// Dimension returns the dimension of embeddings
func (s *OpenAIService) Dimension() int {
	return s.dimension
}

// Close closes the service (no-op for HTTP client)
func (s *OpenAIService) Close() error {
	return nil
}
