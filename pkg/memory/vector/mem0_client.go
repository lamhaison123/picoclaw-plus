// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Mem0Client - Personalized memory provider
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

// Mem0Config is imported from config package but redefined here for package independence
type Mem0Config struct {
	Enabled   bool
	URL       string
	APIKey    string
	TimeoutMS int
}

// Mem0Client implements MemoryProvider interface for Mem0 personalized memory
type Mem0Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
	breaker    *CircuitBreaker
}

// Mem0AddRequest represents a request to add memory
type Mem0AddRequest struct {
	Messages []Mem0Message          `json:"messages"`
	UserID   string                 `json:"user_id,omitempty"`
	AgentID  string                 `json:"agent_id,omitempty"`
	RunID    string                 `json:"run_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Mem0Message represents a message in the conversation
type Mem0Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Mem0AddResponse represents response from add operation
type Mem0AddResponse struct {
	Results []Mem0MemoryResult `json:"results"`
}

// Mem0MemoryResult represents a memory result
type Mem0MemoryResult struct {
	ID        string                 `json:"id,omitempty"`
	Memory    string                 `json:"memory,omitempty"`
	EventID   string                 `json:"event_id,omitempty"` // For async processing
	Status    string                 `json:"status,omitempty"`   // e.g. "PENDING"
	Message   string                 `json:"message,omitempty"`  // Async message
	UserID    string                 `json:"user_id,omitempty"`
	Hash      string                 `json:"hash,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
}

// Mem0SearchRequest represents a search request
type Mem0SearchRequest struct {
	Query   string                 `json:"query"`
	UserID  string                 `json:"user_id,omitempty"`
	AgentID string                 `json:"agent_id,omitempty"`
	RunID   string                 `json:"run_id,omitempty"`
	Limit   int                    `json:"limit,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	Version string                 `json:"version,omitempty"` // "v2"
}

// Mem0SearchResponse represents search results
type Mem0SearchResponse struct {
	Results []Mem0SearchResult `json:"results"`
}

// Mem0SearchResult represents a search result
type Mem0SearchResult struct {
	ID         string                 `json:"id"`
	Memory     string                 `json:"memory"`
	UserID     string                 `json:"user_id,omitempty"`
	Hash       string                 `json:"hash,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Score      float32                `json:"score,omitempty"`
	CreatedAt  string                 `json:"created_at,omitempty"`
	UpdatedAt  string                 `json:"updated_at,omitempty"`
	Categories []string               `json:"categories,omitempty"`
}

// NewMem0Client creates a new Mem0 client instance
func NewMem0Client(cfg Mem0Config, breaker *CircuitBreaker) (*Mem0Client, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("%s: mem0.url is required", ErrConfigInvalid)
	}
	if breaker == nil {
		return nil, fmt.Errorf("%s: circuit breaker is required", ErrConfigInvalid)
	}

	timeoutMs := cfg.TimeoutMS
	if timeoutMs <= 0 {
		timeoutMs = DefaultMemoryTimeoutMs
	}

	return &Mem0Client{
		baseURL: cfg.URL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		},
		timeout: time.Duration(timeoutMs) * time.Millisecond,
		breaker: breaker,
	}, nil
}

// Store stores a memory entry in Mem0
func (c *Mem0Client) Store(ctx context.Context, content string, metadata map[string]interface{}) (string, error) {
	var resultID string
	var resultErr error

	err := c.breaker.Call(ctx, func() error {
		// Extract user_id from metadata if present
		userID := "picoclaw_user" // Default fallback value
		if metadata != nil {
			if uid, ok := metadata["user_id"].(string); ok {
				userID = uid
			}
		}

		req := Mem0AddRequest{
			Messages: []Mem0Message{
				{
					Role:    "user",
					Content: content,
				},
			},
			UserID:   userID,
			Metadata: metadata,
		}

		resp, err := c.doRequest(ctx, "POST", "/v1/memories/", req)
		if err != nil {
			resultErr = err
			return err
		}

		var results []Mem0MemoryResult
		if err := json.Unmarshal(resp, &results); err != nil {
			var errResp map[string]interface{}
			if jsonErr := json.Unmarshal(resp, &errResp); jsonErr == nil {
				resultErr = fmt.Errorf("%s: API returned structured error: %v", ErrInternal, errResp)
			} else {
				resultErr = fmt.Errorf("%s: failed to parse response array (body: %s): %w", ErrInternal, string(resp), err)
			}
			return resultErr
		}

		if len(results) == 0 {
			resultErr = fmt.Errorf("%s: no memory created or queued", ErrInternal)
			return resultErr
		}

		// Support synchronous IDs and asynchronous EventIDs
		if results[0].ID != "" {
			resultID = results[0].ID
		} else if results[0].EventID != "" {
			resultID = results[0].EventID
		} else {
			resultID = "async_queued"
		}
		
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
func (c *Mem0Client) Recall(ctx context.Context, query string, limit int) ([]Memory, error) {
	var resultMemories []Memory
	var resultErr error

	err := c.breaker.Call(ctx, func() error {
		req := Mem0SearchRequest{
			Query:   query,
			UserID:  "picoclaw_user", // Mem0 requires at least one ID
			Limit:   limit,
			Version: "v2",
		}

		resp, err := c.doRequest(ctx, "POST", "/v1/memories/search/", req)
		if err != nil {
			resultErr = err
			return err
		}

		var results []Mem0SearchResult
		if err := json.Unmarshal(resp, &results); err != nil {
			var errResp map[string]interface{}
			if jsonErr := json.Unmarshal(resp, &errResp); jsonErr == nil {
				resultErr = fmt.Errorf("%s: API returned structured error: %v", ErrInternal, errResp)
			} else {
				resultErr = fmt.Errorf("%s: failed to parse search response array: %w", ErrInternal, err)
			}
			return resultErr
		}

		// Convert to Memory slice
		memories := make([]Memory, len(results))
		for i, m := range results {
			// Parse timestamp
			var timestamp int64
			if m.CreatedAt != "" {
				if t, err := time.Parse(time.RFC3339, m.CreatedAt); err == nil {
					timestamp = t.Unix()
				}
			}

			memories[i] = Memory{
				ID:        m.ID,
				Content:   m.Memory,
				Metadata:  m.Metadata,
				Score:     m.Score,
				Timestamp: timestamp,
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
func (c *Mem0Client) Update(ctx context.Context, id string, content string, metadata map[string]interface{}) error {
	return c.breaker.Call(ctx, func() error {
		// Mem0 uses PUT for updates
		req := map[string]interface{}{
			"data": content,
		}
		if metadata != nil {
			req["metadata"] = metadata
		}

		resp, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/v1/memories/%s/", id), req)
		if err != nil {
			return err
		}

		// Check if response indicates success
		var result map[string]interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("%s: failed to parse response: %w", ErrInternal, err)
		}

		return nil
	})
}

// Delete removes a memory entry by ID
func (c *Mem0Client) Delete(ctx context.Context, id string) error {
	return c.breaker.Call(ctx, func() error {
		_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/v1/memories/%s/", id), nil)
		if err != nil {
			return err
		}
		return nil
	})
}

// Health checks the health of the Mem0 connection
func (c *Mem0Client) Health(ctx context.Context) error {
	return c.breaker.Call(ctx, func() error {
		// Mem0 doesn't have a dedicated health endpoint
		// We'll do a lightweight search to verify connectivity
		req := Mem0SearchRequest{
			Query:   "health check",
			UserID:  "picoclaw_user", // Mem0 requires at least one ID
			Limit:   1,
			Version: "v2",
		}

		_, err := c.doRequest(ctx, "POST", "/v1/memories/search/", req)
		if err != nil {
			return fmt.Errorf("%s: health check failed: %w", ErrProviderUnavailable, err)
		}

		return nil
	})
}

// Close closes the Mem0 client connection
func (c *Mem0Client) Close() error {
	// HTTP client doesn't need explicit closing
	// But we can close idle connections
	c.httpClient.CloseIdleConnections()
	return nil
}

// doRequest performs an HTTP request to Mem0 API
func (c *Mem0Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
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

	// Add API key with Token prefix (Mem0 specific)
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Token "+c.apiKey)
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
