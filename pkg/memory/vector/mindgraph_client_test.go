// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: MindGraphClient Tests

package memory

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMindGraphClient(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})

	tests := []struct {
		name    string
		cfg     MindGraphConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: MindGraphConfig{
				URL:       "http://localhost:8002",
				APIKey:    "test-key",
				TimeoutMS: 1200,
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			cfg: MindGraphConfig{
				APIKey:    "test-key",
				TimeoutMS: 1200,
			},
			wantErr: true,
		},
		{
			name: "default timeout",
			cfg: MindGraphConfig{
				URL:    "http://localhost:8002",
				APIKey: "test-key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewMindGraphClient(tt.cfg, breaker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMindGraphClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewMindGraphClient() returned nil client")
			}
			if client != nil {
				_ = client.Close()
			}
		})
	}
}

func TestMindGraphClient_Store(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/ingest/chunk" {
			t.Errorf("Expected path /ingest/chunk, got %s", r.URL.Path)
		}

		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Expected Authorization: Bearer test-key, got %s", auth)
		}

		// Parse request
		var req MindGraphMemoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Send response
		resp := MindGraphMemoryResponse{
			ChunkUID:     "mem-123",
			NodesCreated: 1,
			EdgesCreated: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMindGraphClient(MindGraphConfig{
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
		"type": "test",
		"user": "test-user",
	}

	id, err := client.Store(ctx, "Test memory content", metadata)
	if err != nil {
		t.Errorf("Store() error = %v", err)
	}
	if id != "mem-123" {
		t.Errorf("Store() id = %v, want mem-123", id)
	}
}

func TestMindGraphClient_Recall(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/retrieve/context" {
			t.Errorf("Expected path /retrieve/context, got %s", r.URL.Path)
		}

		// Parse request
		var req MindGraphRecallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Send response
		resp := MindGraphRecallResponse{
			Chunks: []MindGraphChunkResponse{
				{
					ChunkUID: "mem-1",
					Content:  "Memory 1",
					Score:    0.95,
				},
				{
					ChunkUID: "mem-2",
					Content:  "Memory 2",
					Score:    0.85,
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
	client, err := NewMindGraphClient(MindGraphConfig{
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
}

func TestMindGraphClient_Health(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/health" {
			t.Errorf("Expected path /health, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
	})
	client, err := NewMindGraphClient(MindGraphConfig{
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

func TestMindGraphClient_ErrorHandling(t *testing.T) {
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
	client, err := NewMindGraphClient(MindGraphConfig{
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
