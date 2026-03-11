// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Config Unit Tests
// Sprint 1 Implementation - Phase 2

package memory

import (
	"errors"
	"testing"
	"time"
)

func init() {
	skipQdrantInitCheck = true
}

// TestConfig_Validation tests configuration validation logic
func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid qdrant config",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test_collection",
			},
			wantErr: false,
		},
		{
			name: "valid circuit breaker config",
			cfg: CircuitBreakerConfig{
				MaxFailures:   5,
				ResetTimeoutS: 30,
				HalfOpenMax:   3,
			},
			wantErr: false,
		},
		{
			name: "qdrant missing url",
			cfg: QdrantConfig{
				Collection: "test",
			},
			wantErr: true,
			errMsg:  "url is required",
		},
		{
			name: "qdrant missing collection",
			cfg: QdrantConfig{
				URL: "http://localhost:6333",
			},
			wantErr: true,
			errMsg:  "collection is required",
		},
		{
			name: "qdrant invalid dimension",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test",
			},
			wantErr: true,
			errMsg:  "dimension must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			switch cfg := tt.cfg.(type) {
			case QdrantConfig:
				breaker := NewCircuitBreaker(CircuitBreakerConfig{
					MaxFailures:   5,
					ResetTimeoutS: 30,
					HalfOpenMax:   3,
				})
				vectorCfg := VectorConfig{
					Dimension: 384,
					TimeoutMs: 800,
				}
				if tt.name == "qdrant invalid dimension" {
					vectorCfg.Dimension = 0
				}
				_, err = NewQdrantStore(cfg, vectorCfg, breaker)
			case CircuitBreakerConfig:
				breaker := NewCircuitBreaker(cfg)
				if breaker == nil {
					err = errors.New("failed to create circuit breaker")
				}
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errMsg)
				} else if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Error = %v, want error containing '%s'", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error = %v", err)
				}
			}
		})
	}
}

// TestConfig_Defaults tests default value application
func TestConfig_Defaults(t *testing.T) {
	t.Run("qdrant timeout default", func(t *testing.T) {
		cfg := QdrantConfig{
			URL:        "http://localhost:6333",
			Collection: "test",
		}

		vectorCfg := VectorConfig{
			Dimension: 384,
			TimeoutMs: 0, // Should default to 800
		}

		breaker := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures:   5,
			ResetTimeoutS: 30,
			HalfOpenMax:   3,
		})

		store, err := NewQdrantStore(cfg, vectorCfg, breaker)
		if err != nil {
			t.Fatalf("NewQdrantStore() failed: %v", err)
		}

		if store.timeout != 800*time.Millisecond {
			t.Errorf("timeout = %v, want 800ms", store.timeout)
		}
	})

	t.Run("circuit breaker defaults", func(t *testing.T) {
		cfg := CircuitBreakerConfig{} // All zeros

		breaker := NewCircuitBreaker(cfg)

		if breaker.maxFailures != 5 {
			t.Errorf("maxFailures = %d, want 5", breaker.maxFailures)
		}
		if breaker.resetTimeout != 30*time.Second {
			t.Errorf("resetTimeout = %v, want 30s", breaker.resetTimeout)
		}
		if breaker.halfOpenMax != 3 {
			t.Errorf("halfOpenMax = %d, want 3", breaker.halfOpenMax)
		}
	})
}

// TestConfig_SpecCompliance tests compliance with locked specs
func TestConfig_SpecCompliance(t *testing.T) {
	t.Run("timeout specs", func(t *testing.T) {
		// Vector Search: 800ms
		cfg := QdrantConfig{
			URL:        "http://localhost:6333",
			Collection: "test",
		}
		vectorCfg := VectorConfig{
			Dimension: 384,
			TimeoutMs: 800,
		}

		breaker := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures:   5,
			ResetTimeoutS: 30,
			HalfOpenMax:   3,
		})

		store, err := NewQdrantStore(cfg, vectorCfg, breaker)
		if err != nil {
			t.Fatalf("NewQdrantStore() failed: %v", err)
		}

		if store.timeout != 800*time.Millisecond {
			t.Errorf("timeout = %v, want 800ms per specs", store.timeout)
		}
	})

	t.Run("circuit breaker specs", func(t *testing.T) {
		// Failures: 5, Reset: 30s, Half-open: 3
		cfg := CircuitBreakerConfig{
			MaxFailures:   5,
			ResetTimeoutS: 30,
			HalfOpenMax:   3,
		}

		breaker := NewCircuitBreaker(cfg)

		if breaker.maxFailures != 5 {
			t.Errorf("maxFailures = %d, want 5 per specs", breaker.maxFailures)
		}
		if breaker.resetTimeout != 30*time.Second {
			t.Errorf("resetTimeout = %v, want 30s per specs", breaker.resetTimeout)
		}
		if breaker.halfOpenMax != 3 {
			t.Errorf("halfOpenMax = %d, want 3 per specs", breaker.halfOpenMax)
		}
	})

	t.Run("dimension specs", func(t *testing.T) {
		// Dimension: 384
		cfg := QdrantConfig{
			URL:        "http://localhost:6333",
			Collection: "test",
		}
		vectorCfg := VectorConfig{
			Dimension: 384,
			TimeoutMs: 800,
		}

		breaker := NewCircuitBreaker(CircuitBreakerConfig{
			MaxFailures:   5,
			ResetTimeoutS: 30,
			HalfOpenMax:   3,
		})

		store, err := NewQdrantStore(cfg, vectorCfg, breaker)
		if err != nil {
			t.Fatalf("NewQdrantStore() failed: %v", err)
		}

		if store.dimension != 384 {
			t.Errorf("dimension = %d, want 384 per specs", store.dimension)
		}
	})
}

// TestConfig_EdgeCases tests edge case handling
func TestConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		cfg     QdrantConfig
		wantErr bool
	}{
		{
			name: "negative timeout",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test",
			},
			wantErr: false, // Should apply default
		},
		{
			name: "very large timeout",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test",
			},
			wantErr: false, // Should accept
		},
		{
			name: "very large dimension",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "test",
			},
			wantErr: false, // Should accept
		},
		{
			name: "empty collection name",
			cfg: QdrantConfig{
				URL:        "http://localhost:6333",
				Collection: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			breaker := NewCircuitBreaker(CircuitBreakerConfig{
				MaxFailures:   5,
				ResetTimeoutS: 30,
				HalfOpenMax:   3,
			})

			vectorCfg := VectorConfig{
				Dimension: 384,
				TimeoutMs: 800,
			}
			switch tt.name {
			case "very large dimension":
				vectorCfg.Dimension = 10000
			case "negative timeout":
				vectorCfg.TimeoutMs = -100
			case "very large timeout":
				vectorCfg.TimeoutMs = 999999
			}
			_, err := NewQdrantStore(tt.cfg, vectorCfg, breaker)

			if tt.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error = %v", err)
			}
		})
	}
}
