// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: MindGraphClient - Knowledge graph memory provider
// Sprint 1 Implementation - Phase 1

package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MindGraphConfig is imported from config package but redefined here for package independence
type MindGraphConfig struct {
	Enabled   bool
	URL       string
	APIKey    string
	TimeoutMS int
}

// MindGraphClient implements MemoryProvider interface for MindGraph knowledge graph
type MindGraphClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
	breaker    *CircuitBreaker
}

// MindGraphMemoryRequest represents a request to store memory
type MindGraphMemoryRequest struct {
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Layers   []string               `json:"layers"`
}

// MindGraphMemoryResponse represents a response from MindGraph Store
type MindGraphMemoryResponse struct {
	ChunkUID     string `json:"chunk_uid"`
	NodesCreated int    `json:"nodes_created"`
	EdgesCreated int    `json:"edges_created"`
}

// MindGraphRecallRequest represents a request to recall memories
type MindGraphRecallRequest struct {
	Query string `json:"query"`
	K     int    `json:"k"`
	Depth int    `json:"depth,omitempty"`
}

// MindGraphNodeResponse represents a node in the recall or generic CRUD response
type MindGraphNodeResponse struct {
	UID       string                 `json:"uid"`
	Label     string                 `json:"label"`
	Summary   string                 `json:"summary,omitempty"`
	NodeType  string                 `json:"node_type,omitempty"`
	Props     map[string]interface{} `json:"props,omitempty"`
	CreatedAt float64                `json:"created_at,omitempty"`
}

// MindGraphChunkResponse represents a chunk from the retrieve response
type MindGraphChunkResponse struct {
	ChunkUID string  `json:"chunk_uid"`
	Content  string  `json:"content"`
	Score    float64 `json:"score"`
}

// MindGraphRecallResponse represents recall results
type MindGraphRecallResponse struct {
	Chunks []MindGraphChunkResponse `json:"chunks"`
	Graph  struct {
		Nodes []MindGraphNodeResponse `json:"nodes"`
	} `json:"graph"`
}

// NewMindGraphClient creates a new MindGraph client instance
func NewMindGraphClient(cfg MindGraphConfig, breaker *CircuitBreaker) (*MindGraphClient, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("%s: mindgraph.url is required", ErrConfigInvalid)
	}
	if breaker == nil {
		return nil, fmt.Errorf("%s: circuit breaker is required", ErrConfigInvalid)
	}

	timeoutMs := cfg.TimeoutMS
	if timeoutMs <= 0 {
		timeoutMs = DefaultMemoryTimeoutMs
	}

	return &MindGraphClient{
		baseURL: cfg.URL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		},
		timeout: time.Duration(timeoutMs) * time.Millisecond,
		breaker: breaker,
	}, nil
}

// Store stores a memory entry in MindGraph
func (c *MindGraphClient) Store(ctx context.Context, content string, metadata map[string]interface{}) (string, error) {
	var resultID string
	var resultErr error

	err := c.breaker.Call(ctx, func() error {
		req := MindGraphMemoryRequest{
			Content:  content,
			Metadata: metadata,
			Layers:   []string{"memory"},
		}

		resp, err := c.doRequest(ctx, "POST", "/ingest/chunk", req)
		if err != nil {
			resultErr = err
			return err
		}

		var result MindGraphMemoryResponse
		if err := json.Unmarshal(resp, &result); err != nil {
			resultErr = fmt.Errorf("%s: failed to parse response: %w", ErrInternal, err)
			return resultErr
		}

		if result.ChunkUID == "" {
			resultErr = fmt.Errorf("%s: no chunk_uid returned", ErrInternal)
			return resultErr
		}

		resultID = result.ChunkUID
		return nil
	})

	if err != nil {
		return "", err
	}
	if resultErr != nil {
		return "", resultErr
	}
	return resultID, nil
}

// Recall retrieves relevant memories based on query
func (c *MindGraphClient) Recall(ctx context.Context, query string, limit int) ([]Memory, error) {
	var resultMemories []Memory
	var resultErr error

	err := c.breaker.Call(ctx, func() error {
		req := MindGraphRecallRequest{
			Query: query,
			K:     limit,
			Depth: 1, // Default retrieval depth
		}

		resp, err := c.doRequest(ctx, "POST", "/retrieve/context", req)
		if err != nil {
			resultErr = err
			return err
		}

		var result MindGraphRecallResponse
		if err := json.Unmarshal(resp, &result); err != nil {
			resultErr = fmt.Errorf("%s: failed to parse response: %w, raw body: %s", ErrInternal, err, string(resp))
			return resultErr
		}

		// Convert to Memory slice
		memories := make([]Memory, len(result.Chunks))
		for i, c := range result.Chunks {
			memories[i] = Memory{
				ID:        c.ChunkUID,
				Content:   c.Content,
				Metadata:  map[string]interface{}{"score": c.Score},
				Timestamp: time.Now().Unix(),
			}
		}

		resultMemories = memories
		return nil
	})

	if err != nil {
		return nil, err
	}
	if resultErr != nil {
		return nil, resultErr
	}
	return resultMemories, nil
}

// Update updates an existing memory entry
func (c *MindGraphClient) Update(ctx context.Context, id string, content string, metadata map[string]interface{}) error {
	return c.breaker.Call(ctx, func() error {
		req := struct {
			Label string                 `json:"label"`
			Props map[string]interface{} `json:"props,omitempty"`
		}{
			Label: content,
			Props: metadata,
		}

		_, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/node/%s", id), req)
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete removes a memory entry by ID
func (c *MindGraphClient) Delete(ctx context.Context, id string) error {
	return c.breaker.Call(ctx, func() error {
		_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/node/%s", id), nil)
		if err != nil {
			return err
		}

		return nil
	})
}

// Health checks the health of the MindGraph connection
func (c *MindGraphClient) Health(ctx context.Context) error {
	return c.breaker.Call(ctx, func() error {
		_, err := c.doRequest(ctx, "GET", "/health", nil)
		if err != nil {
			return err
		}

		return nil
	})
}

// Close closes the MindGraph client connection
func (c *MindGraphClient) Close() error {
	// HTTP client doesn't need explicit closing
	// But we can close idle connections
	c.httpClient.CloseIdleConnections()
	return nil
}

// doRequest performs an HTTP request to MindGraph API
func (c *MindGraphClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to marshal request: %w", ErrInternal, err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create request: %w", ErrInternal, err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add API key if provided
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read response: %w", ErrInternal, err)
	}

	// Check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: HTTP %d: %s", ErrProviderUnavailable, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
