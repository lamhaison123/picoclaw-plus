// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Mem0Client Tests

package memory

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMem0Client(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})

	tests := []struct {
		name    string
		cfg     Mem0Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Mem0Config{
				URL:       "http://localhost:8001",
				APIKey:    "test-key",
				TimeoutMS: 1200,
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			cfg: Mem0Config{
				APIKey:    "test-key",
				TimeoutMS: 1200,
			},
			wantErr: true,
		},
		{
			name: "default timeout",
			cfg: Mem0Config{
				URL:    "http://localhost:8001",
				APIKey: "test-key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewMem0Client(tt.cfg, breaker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMem0Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewMem0Client() returned nil client")
			}
			if client != nil {
				_ = client.Close()
			}
		})
	}
}

func TestMem0Client_Store(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/" {
			t.Errorf("Expected path /v1/memories/, got %s", r.URL.Path)
		}

		// Check authorization header (Token prefix)
		auth := r.Header.Get("Authorization")
		if auth != "Token test-key" {
			t.Errorf("Expected Authorization: Token test-key, got %s", auth)
		}

		// Parse request
		var req Mem0AddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Send response
		resp := Mem0AddResponse{
			Results: []Mem0MemoryResult{
				{
					ID:        "mem-123",
					Memory:    req.Messages[0].Content,
					UserID:    req.UserID,
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	metadata := map[string]interface{}{
		"user_id": "user-123",
		"type":    "test",
	}

	id, err := client.Store(ctx, "Test memory content", metadata)
	if err != nil {
		t.Errorf("Store() error = %v", err)
	}
	if id != "mem-123" {
		t.Errorf("Store() id = %v, want mem-123", id)
	}
}

func TestMem0Client_Recall(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/search/" {
			t.Errorf("Expected path /v1/memories/search/, got %s", r.URL.Path)
		}

		// Parse request
		var req Mem0SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Send response
		resp := Mem0SearchResponse{
			Results: []Mem0SearchResult{
				{
					ID:        "mem-1",
					Memory:    "Memory 1",
					Score:     0.95,
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "mem-2",
					Memory:    "Memory 2",
					Score:     0.85,
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	memories, err := client.Recall(ctx, "test query", 10)
	if err != nil {
		t.Errorf("Recall() error = %v", err)
	}
	if len(memories) != 2 {
		t.Errorf("Recall() returned %d memories, want 2", len(memories))
	}
	if memories[0].ID != "mem-1" {
		t.Errorf("Recall() first memory ID = %v, want mem-1", memories[0].ID)
	}
	if memories[0].Score != 0.95 {
		t.Errorf("Recall() first memory score = %v, want 0.95", memories[0].Score)
	}
}

func TestMem0Client_Update(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/mem-123/" {
			t.Errorf("Expected path /v1/memories/mem-123/, got %s", r.URL.Path)
		}

		// Send success response
		resp := map[string]interface{}{
			"success": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.Update(ctx, "mem-123", "Updated content", nil)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}
}

func TestMem0Client_Delete(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/mem-123/" {
			t.Errorf("Expected path /v1/memories/mem-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.Delete(ctx, "mem-123")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}

func TestMem0Client_Health(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/search/" {
			t.Errorf("Expected path /v1/memories/search/, got %s", r.URL.Path)
		}

		resp := Mem0SearchResponse{
			Results: []Mem0SearchResult{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Health(ctx); err != nil {
		t.Errorf("Health() error = %v", err)
	}
}

func TestMem0Client_ErrorHandling(t *testing.T) {
	// Create mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMem0Client(Mem0Config{
		URL:       server.URL,
		APIKey:    "test-key",
		TimeoutMS: 5000,
	}, breaker)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Test Store error
	_, err = client.Store(ctx, "test", nil)
	if err == nil {
		t.Error("Store() expected error, got nil")
	}

	// Test Recall error
	_, err = client.Recall(ctx, "test", 10)
	if err == nil {
		t.Error("Recall() expected error, got nil")
	}

	// Test Health error
	err = client.Health(ctx)
	if err == nil {
		t.Error("Health() expected error, got nil")
	}
}
