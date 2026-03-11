// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Configuration (locked specs from Sprint 1)
// Sprint 1 Implementation - Phase 1

package memory

import (
	"fmt"
	"strings"
)

// MemoryConfig is the root configuration for memory subsystem
type MemoryConfig struct {
	Vector VectorConfig         `json:"vector"`
	Memory MemoryProviderConfig `json:"memory"`
	Global GlobalMemoryConfig   `json:"global"`
}

// VectorConfig holds vector store configuration
type VectorConfig struct {
	Enabled   bool   `json:"enabled"`
	Provider  string `json:"provider"`   // "qdrant" | "lancedb"
	TimeoutMs int    `json:"timeout_ms"` // default: 800
	Dimension int    `json:"dimension"`  // default: 384

	Qdrant  QdrantConfig  `json:"qdrant,omitempty"`
	LanceDB LanceDBConfig `json:"lancedb,omitempty"`
}

// QdrantConfig holds Qdrant-specific configuration
type QdrantConfig struct {
	URL        string `json:"url"`
	Collection string `json:"collection"`
	APIKey     string `json:"api_key,omitempty"`
}

// LanceDBConfig holds LanceDB-specific configuration
type LanceDBConfig struct {
	Mode string `json:"mode"` // "cgo" | "api"
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}

// MemoryProviderConfig holds memory provider configuration
type MemoryProviderConfig struct {
	Enabled   bool   `json:"enabled"`
	Provider  string `json:"provider"`   // "mem0" | "mindgraph" | "both" | "sidecar"
	TimeoutMs int    `json:"timeout_ms"` // default: 1200

	Sidecar        SidecarConfig        `json:"sidecar"`
	CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`
	Retry          RetryConfig          `json:"retry"`
}

// SidecarConfig holds sidecar endpoint configuration
type SidecarConfig struct {
	Endpoint string `json:"endpoint"`
}

// RetryConfig holds retry policy configuration
type RetryConfig struct {
	MaxAttempts int `json:"max_attempts"` // default: 2
	BackoffMs   int `json:"backoff_ms"`   // default: 100
}

// GlobalMemoryConfig holds global memory settings
type GlobalMemoryConfig struct {
	WatchdogTimeoutMs int    `json:"watchdog_timeout_ms"` // default: 5000
	FallbackMode      string `json:"fallback_mode"`       // "graceful" | "strict"
}

// Default values (locked from Sprint 1 specs)
const (
	DefaultVectorTimeoutMs   = 800
	DefaultVectorDimension   = 384
	DefaultMemoryTimeoutMs   = 1200
	DefaultMaxFailures       = 5
	DefaultResetTimeoutS     = 30
	DefaultMaxAttempts       = 2
	DefaultBackoffMs         = 100
	DefaultWatchdogTimeoutMs = 5000
	DefaultFallbackMode      = "graceful"
)

// Validate validates the memory configuration
func (c *MemoryConfig) Validate() error {
	// Apply defaults
	c.applyDefaults()

	// Vector validation
	if c.Vector.Enabled {
		if err := c.validateVector(); err != nil {
			return err
		}
	}

	// Memory validation
	if c.Memory.Enabled {
		if err := c.validateMemory(); err != nil {
			return err
		}
	}

	// Global validation
	if err := c.validateGlobal(); err != nil {
		return err
	}

	return nil
}

// applyDefaults applies default values to configuration
func (c *MemoryConfig) applyDefaults() {
	if c.Vector.TimeoutMs <= 0 {
		c.Vector.TimeoutMs = DefaultVectorTimeoutMs
	}
	if c.Vector.Dimension <= 0 {
		c.Vector.Dimension = DefaultVectorDimension
	}
	if c.Memory.TimeoutMs <= 0 {
		c.Memory.TimeoutMs = DefaultMemoryTimeoutMs
	}
	if c.Memory.CircuitBreaker.MaxFailures <= 0 {
		c.Memory.CircuitBreaker.MaxFailures = DefaultMaxFailures
	}
	if c.Memory.CircuitBreaker.ResetTimeoutS <= 0 {
		c.Memory.CircuitBreaker.ResetTimeoutS = DefaultResetTimeoutS
	}
	if c.Memory.Retry.MaxAttempts <= 0 {
		c.Memory.Retry.MaxAttempts = DefaultMaxAttempts
	}
	if c.Memory.Retry.BackoffMs <= 0 {
		c.Memory.Retry.BackoffMs = DefaultBackoffMs
	}
	if c.Global.WatchdogTimeoutMs <= 0 {
		c.Global.WatchdogTimeoutMs = DefaultWatchdogTimeoutMs
	}
	if c.Global.FallbackMode == "" {
		c.Global.FallbackMode = DefaultFallbackMode
	}
}

// validateVector validates vector store configuration
func (c *MemoryConfig) validateVector() error {
	if c.Vector.Provider == "" {
		return fmt.Errorf("%s: vector.provider required when vector.enabled=true", ErrConfigInvalid)
	}

	if c.Vector.Provider != "qdrant" && c.Vector.Provider != "lancedb" {
		return fmt.Errorf("%s: vector.provider must be 'qdrant' or 'lancedb', got: %s",
			ErrConfigInvalid, c.Vector.Provider)
	}

	if c.Vector.TimeoutMs <= 0 {
		return fmt.Errorf("%s: vector.timeout_ms must be positive, got: %d",
			ErrConfigInvalid, c.Vector.TimeoutMs)
	}

	if c.Vector.Dimension <= 0 {
		return fmt.Errorf("%s: vector.dimension must be positive, got: %d",
			ErrConfigInvalid, c.Vector.Dimension)
	}

	// Provider-specific validation
	if c.Vector.Provider == "qdrant" {
		if c.Vector.Qdrant.URL == "" {
			return fmt.Errorf("%s: vector.qdrant.url required when provider=qdrant", ErrConfigInvalid)
		}
		if c.Vector.Qdrant.Collection == "" {
			return fmt.Errorf("%s: vector.qdrant.collection required when provider=qdrant", ErrConfigInvalid)
		}
	}

	if c.Vector.Provider == "lancedb" {
		if err := c.validateLanceDB(); err != nil {
			return err
		}
	}

	return nil
}

// validateLanceDB validates LanceDB configuration with CGO fallback policy
func (c *MemoryConfig) validateLanceDB() error {
	if c.Vector.LanceDB.Mode == "" {
		return fmt.Errorf("%s: vector.lancedb.mode required when provider=lancedb", ErrConfigInvalid)
	}

	// LanceDB CGO Policy (locked from Sprint 1)
	if c.Vector.LanceDB.Mode == "cgo" && !cgoAvailable() {
		// Log warning
		fmt.Printf("WARNING: LanceDB CGO mode unavailable\n")

		// Try API mode fallback
		if c.Vector.LanceDB.URL != "" {
			fmt.Printf("INFO: Falling back to LanceDB API mode\n")
			c.Vector.LanceDB.Mode = "api"
		} else {
			// No fallback available
			fmt.Printf("WARNING: No LanceDB fallback available, disabling vector store\n")
			c.Vector.Enabled = false

			if c.Global.FallbackMode == "strict" {
				return fmt.Errorf("%s: LanceDB CGO unavailable and no API fallback configured", ErrConfigInvalid)
			}
			// graceful mode: continue with warning
		}
	}

	// Validate mode-specific requirements
	if c.Vector.LanceDB.Mode == "cgo" && c.Vector.LanceDB.Path == "" {
		return fmt.Errorf("%s: vector.lancedb.path required when mode=cgo", ErrConfigInvalid)
	}

	if c.Vector.LanceDB.Mode == "api" && c.Vector.LanceDB.URL == "" {
		return fmt.Errorf("%s: vector.lancedb.url required when mode=api", ErrConfigInvalid)
	}

	return nil
}

// validateMemory validates memory provider configuration
func (c *MemoryConfig) validateMemory() error {
	if c.Memory.Provider == "" {
		return fmt.Errorf("%s: memory.provider required when memory.enabled=true", ErrConfigInvalid)
	}

	validProviders := []string{"mem0", "mindgraph", "both", "sidecar"}
	if !containsProvider(validProviders, c.Memory.Provider) {
		return fmt.Errorf("%s: memory.provider must be one of %v, got: %s",
			ErrConfigInvalid, validProviders, c.Memory.Provider)
	}

	if c.Memory.TimeoutMs <= 0 {
		return fmt.Errorf("%s: memory.timeout_ms must be positive, got: %d",
			ErrConfigInvalid, c.Memory.TimeoutMs)
	}

	// Only require sidecar endpoint when using sidecar provider
	if c.Memory.Provider == "sidecar" && c.Memory.Sidecar.Endpoint == "" {
		return fmt.Errorf("%s: memory.sidecar.endpoint required when provider=sidecar", ErrConfigInvalid)
	}

	// Circuit breaker validation
	if c.Memory.CircuitBreaker.MaxFailures <= 0 {
		return fmt.Errorf("%s: circuit_breaker.max_failures must be positive", ErrConfigInvalid)
	}

	if c.Memory.CircuitBreaker.ResetTimeoutS <= 0 {
		return fmt.Errorf("%s: circuit_breaker.reset_timeout_s must be positive", ErrConfigInvalid)
	}

	return nil
}

// validateGlobal validates global configuration
func (c *MemoryConfig) validateGlobal() error {
	if c.Global.FallbackMode != "graceful" && c.Global.FallbackMode != "strict" {
		return fmt.Errorf("%s: global.fallback_mode must be 'graceful' or 'strict', got: %s",
			ErrConfigInvalid, c.Global.FallbackMode)
	}

	return nil
}

// cgoAvailable checks if CGO is available in the build
func cgoAvailable() bool {
	return isCGOEnabled()
}

// Helper function to check if slice contains string
func containsProvider(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
