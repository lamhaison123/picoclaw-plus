// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: QdrantStore - Vector database integration
// Sprint 1 Implementation - Phase 1
//
// NOTE: This file requires github.com/qdrant/go-client/qdrant dependency
// To build without Qdrant support, use build tag: -tags=no_qdrant
//
//go:build !no_qdrant
// +build !no_qdrant

package memory

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

// skipQdrantInitCheck is an internal test hook to bypass network I/O during initialization
var skipQdrantInitCheck = false

// QdrantStore implements VectorStore interface for Qdrant vector database
type QdrantStore struct {
	client     *qdrant.Client
	collection string
	dimension  int
	timeout    time.Duration
	breaker    *CircuitBreaker
	maxRetries int
	retryDelay time.Duration
}

// NewQdrantStore creates a new Qdrant vector store instance
// It validates the configuration and ensures the collection exists with correct dimension
func NewQdrantStore(cfg QdrantConfig, vectorCfg VectorConfig, breaker *CircuitBreaker) (*QdrantStore, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: qdrant.url is required")
	}
	if cfg.Collection == "" {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: qdrant.collection is required")
	}
	if vectorCfg.Dimension <= 0 {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: dimension must be positive, got: %d", vectorCfg.Dimension)
	}
	if breaker == nil {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: circuit breaker is required")
	}

	timeoutMs := vectorCfg.TimeoutMs
	if timeoutMs <= 0 {
		timeoutMs = DefaultVectorTimeoutMs
	}

	maxRetries := DefaultMaxAttempts
	retryDelay := DefaultBackoffMs

	// Parse URL for gRPC Host and Port
	rawURL := cfg.URL
	if !strings.Contains(rawURL, "://") {
		rawURL = "grpc://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	var host string
	var port int
	var useTLS bool

	if err == nil && parsedURL.Host != "" {
		host = parsedURL.Hostname()
		portStr := parsedURL.Port()
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		} else {
			if parsedURL.Scheme == "https" {
				port = 443
			} else {
				port = 6334 // Default Qdrant gRPC port
			}
		}
		useTLS = parsedURL.Scheme == "https"
	} else {
		// Fallback for raw host string
		host = cfg.URL
		port = 6334
		useTLS = false
	}

	// Create Qdrant client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   port,
		UseTLS: useTLS,
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to create Qdrant client: %w", err)
	}

	// Verify collection exists and has correct dimension
	if !skipQdrantInitCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collInfo, err := client.GetCollectionInfo(ctx, cfg.Collection)
		if err != nil {
			// Collection doesn't exist - try to create it
			createErr := client.CreateCollection(ctx, &qdrant.CreateCollection{
				CollectionName: cfg.Collection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     uint64(vectorCfg.Dimension),
					Distance: qdrant.Distance_Cosine,
				}),
			})
			if createErr != nil {
				return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: collection not found and failed to create: %w", createErr)
			}
		} else {
			// Collection exists - verify dimension matches
			if collInfo.Config != nil && collInfo.Config.Params != nil {
				var collDimension uint64

				// Handle both single vector and named vectors config
				if collInfo.Config.Params.VectorsConfig != nil {
					if vectorsConfig := collInfo.Config.Params.VectorsConfig.GetParams(); vectorsConfig != nil {
						collDimension = vectorsConfig.Size
					} else if paramsMap := collInfo.Config.Params.VectorsConfig.GetParamsMap(); paramsMap != nil && len(paramsMap.GetMap()) > 0 {
						// Use first named vector dimension if exist
						for _, v := range paramsMap.GetMap() {
							collDimension = v.Size
							break
						}
					}
				}

				if collDimension > 0 && collDimension != uint64(vectorCfg.Dimension) {
					return nil, fmt.Errorf("ERR_DIMENSION_MISMATCH: collection has dimension %d, config specifies %d",
						collDimension, vectorCfg.Dimension)
				}
			}
		}
	}

	return &QdrantStore{
		client:     client,
		collection: cfg.Collection,
		dimension:  vectorCfg.Dimension,
		timeout:    time.Duration(timeoutMs) * time.Millisecond,
		breaker:    breaker,
		maxRetries: maxRetries,
		retryDelay: time.Duration(retryDelay) * time.Millisecond,
	}, nil
}

// withTimeout applies timeout to context respecting parent deadline
func (q *QdrantStore) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
		// Check if parent deadline is sooner than our timeout
		remaining := time.Until(deadline)
		if remaining < q.timeout {
			// Parent deadline is sooner, use it but still return proper cancel
			return context.WithCancel(ctx)
		}
	}
	// Create new timeout context
	return context.WithTimeout(ctx, q.timeout)
}

// executeWithRetry wraps a function call with circuit breaker and retry logic
// Uses exponential backoff with jitter to avoid thundering herd
func (q *QdrantStore) executeWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	baseDelay := q.retryDelay

	for attempt := 0; attempt < q.maxRetries; attempt++ {
		err := q.breaker.Call(ctx, fn)
		if err == nil {
			return nil
		}

		lastErr = err

		// If it's a context error or circuit open, don't retry
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) ||
			strings.Contains(err.Error(), "ERR_CIRCUIT_OPEN") {
			return err
		}

		// Check if error is non-retriable
		if strings.Contains(err.Error(), "ERR_AUTH_FAILED") ||
			strings.Contains(err.Error(), "ERR_CONFIG_INVALID") ||
			strings.Contains(err.Error(), "ERR_DIMENSION_MISMATCH") {
			return err
		}

		// If this is the last attempt, don't wait
		if attempt == q.maxRetries-1 {
			break
		}

		// Calculate exponential backoff with overflow protection
		var backoff time.Duration
		if attempt < 30 { // 2^30 = 1 billion, safe for time.Duration
			backoff = baseDelay * time.Duration(1<<uint(attempt))
		} else {
			backoff = 5 * time.Second // Max out for large attempts
		}

		// Cap at 5 seconds and handle overflow
		if backoff > 5*time.Second || backoff < 0 {
			backoff = 5 * time.Second
		}

		// Add jitter (±25%)
		jitter := time.Duration(float64(backoff) * 0.25 * (2*rand.Float64() - 1))
		delay := backoff + jitter
		if delay < 0 {
			delay = baseDelay
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
	}
	return lastErr
}

// Upsert inserts or updates vectors in the store
func (q *QdrantStore) Upsert(ctx context.Context, vectors []Vector) error {
	if len(vectors) == 0 {
		return nil
	}

	ctx, cancel := q.withTimeout(ctx)
	defer cancel()

	// Process in batches to avoid memory issues and timeouts
	const maxBatchSize = 100
	successCount := 0
	for i := 0; i < len(vectors); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(vectors) {
			end = len(vectors)
		}

		batch := vectors[i:end]
		if err := q.upsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to upsert batch %d-%d (successfully upserted %d vectors): %w",
				i, end, successCount, err)
		}
		successCount += len(batch)
	}

	return nil
}

// parseQdrantID converts a string ID to a Qdrant PointId ensuring consistency across Upsert and Delete.
func parseQdrantID(id string) *qdrant.PointId {
	if id == "" {
		return qdrant.NewIDUUID(uuid.New().String())
	}
	if numID, err := strconv.ParseUint(id, 10, 64); err == nil {
		return qdrant.NewIDNum(numID) // valid integer ID
	}
	if _, err := uuid.Parse(id); err == nil {
		return qdrant.NewIDUUID(id) // already a UUID
	}
	// Deterministically hash regular strings to UUIDs
	newUUID := uuid.NewMD5(uuid.NameSpaceOID, []byte(id))
	return qdrant.NewIDUUID(newUUID.String())
}

// upsertBatch handles a single batch of vectors
func (q *QdrantStore) upsertBatch(ctx context.Context, vectors []Vector) error {
	if len(vectors) == 0 {
		return nil
	}

	return q.executeWithRetry(ctx, func() error {
		// Validate dimensions
		for i, v := range vectors {
			if len(v.Embedding) != q.dimension {
				return fmt.Errorf("ERR_DIMENSION_MISMATCH: vector[%d] has dimension %d, expected %d",
					i, len(v.Embedding), q.dimension)
			}
		}

		// Convert to Qdrant points
		points := make([]*qdrant.PointStruct, len(vectors))
		for i, v := range vectors {
			// Parse ID properly using unified converter
			pointID := parseQdrantID(v.ID)

			// Convert metadata to Qdrant Value format
			payload := make(map[string]*qdrant.Value)
			for k, v := range v.Metadata {
				val, err := qdrant.NewValue(v)
				if err != nil {
					// Skip invalid values
					continue
				}
				payload[k] = val
			}

			points[i] = &qdrant.PointStruct{
				Id:      pointID,
				Vectors: qdrant.NewVectors(v.Embedding...),
				Payload: payload,
			}
		}

		// Upsert to Qdrant
		_, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: q.collection,
			Points:         points,
		})
		if err != nil {
			return q.mapError(ctx, err)
		}

		return nil
	})
}

// Search performs vector similarity search
func (q *QdrantStore) Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error) {
	// Validate topK parameter
	if topK <= 0 {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: topK must be positive, got: %d", topK)
	}
	if topK > 1000 {
		topK = 1000 // Cap at reasonable limit to prevent performance issues
	}

	ctx, cancel := q.withTimeout(ctx)
	defer cancel()

	var results []SearchResult

	err := q.executeWithRetry(ctx, func() error {
		// Validate dimension
		if len(query.Embedding) != q.dimension {
			return fmt.Errorf("ERR_DIMENSION_MISMATCH: query has dimension %d, expected %d",
				len(query.Embedding), q.dimension)
		}

		// Search in Qdrant using Query API
		searchResults, err := q.client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: q.collection,
			Query:          qdrant.NewQuery(query.Embedding...),
			Limit:          qdrant.PtrOf(uint64(topK)),
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err != nil {
			return q.mapError(ctx, err)
		}

		// Convert results preserving ID type information
		results = make([]SearchResult, len(searchResults))
		for i, r := range searchResults {
			var idStr string
			switch id := r.Id.GetPointIdOptions().(type) {
			case *qdrant.PointId_Num:
				idStr = strconv.FormatUint(id.Num, 10)
			case *qdrant.PointId_Uuid:
				idStr = id.Uuid
			default:
				idStr = fmt.Sprintf("%v", r.Id)
			}

			// Convert Qdrant Value payload back to interface{}
			metadata := make(map[string]interface{})
			for k, v := range r.Payload {
				switch kind := v.GetKind().(type) {
				case *qdrant.Value_StringValue:
					metadata[k] = kind.StringValue
				case *qdrant.Value_IntegerValue:
					metadata[k] = kind.IntegerValue
				case *qdrant.Value_DoubleValue:
					metadata[k] = kind.DoubleValue
				case *qdrant.Value_BoolValue:
					metadata[k] = kind.BoolValue
				default:
					metadata[k] = v.GetKind()
				}
			}

			results[i] = SearchResult{
				ID:       idStr,
				Score:    r.Score,
				Metadata: metadata,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// Delete removes vectors from the store
func (q *QdrantStore) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	ctx, cancel := q.withTimeout(ctx)
	defer cancel()

	return q.executeWithRetry(ctx, func() error {
		// Convert IDs to Qdrant format uniformly across Upsert and Delete
		pointIDs := make([]*qdrant.PointId, len(ids))
		for i, id := range ids {
			pointIDs[i] = parseQdrantID(id)
		}

		// Delete from Qdrant
		_, err := q.client.Delete(ctx, &qdrant.DeletePoints{
			CollectionName: q.collection,
			Points:         qdrant.NewPointsSelector(pointIDs...),
		})
		if err != nil {
			return q.mapError(ctx, err)
		}

		return nil
	})
}

// Health checks the health of the Qdrant connection
func (q *QdrantStore) Health(ctx context.Context) error {
	ctx, cancel := q.withTimeout(ctx)
	defer cancel()

	// Check collection exists
	_, err := q.client.GetCollectionInfo(ctx, q.collection)
	if err != nil {
		return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: health check failed: %w", err)
	}

	return nil
}

// Close closes the Qdrant client connection
func (q *QdrantStore) Close() error {
	// Qdrant Go client doesn't require explicit close
	// Connection pooling is handled internally
	return nil
}

// mapError maps Qdrant errors to canonical error codes
func (q *QdrantStore) mapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Check context errors first
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("ERR_TIMEOUT: operation timed out: %w", err)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("ERR_CANCELED: operation canceled: %w", err)
	}

	errStr := strings.ToLower(err.Error())

	// Map Qdrant-specific errors to canonical codes
	switch {
	case strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "authentication"):
		return fmt.Errorf("ERR_AUTH_FAILED: %w", err)
	case strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "too many requests"):
		return fmt.Errorf("ERR_RATE_LIMITED: %w", err)
	case strings.Contains(errStr, "collection") && strings.Contains(errStr, "not found"):
		return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: collection not found: %w", err)
	case strings.Contains(errStr, "dimension"):
		return fmt.Errorf("ERR_DIMENSION_MISMATCH: %w", err)
	case strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dial"):
		return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: cannot connect to Qdrant: %w", err)
	case strings.Contains(errStr, "timeout"):
		return fmt.Errorf("ERR_TIMEOUT: %w", err)
	default:
		return fmt.Errorf("ERR_INTERNAL: %w", err)
	}
}
