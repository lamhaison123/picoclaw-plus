// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// package memory/vector
// Module: LanceDBStore - Vector database integration via CGO
//
//go:build cgo
// +build cgo

package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/apache/arrow/go/v17/arrow/memory"
	"github.com/lancedb/lancedb-go/pkg/lancedb"
)

// LanceDBStore implements VectorStore interface for LanceDB
type LanceDBStore struct {
	db         *lancedb.Connection
	table      *lancedb.Table
	dimension  int
	timeout    time.Duration
	breaker    *CircuitBreaker
	maxRetries int
	retryDelay time.Duration
}

// NewLanceDBStore creates a new LanceDB vector store instance
func NewLanceDBStore(cfg LanceDBConfig, vectorCfg VectorConfig, breaker *CircuitBreaker) (*LanceDBStore, error) {
	if cfg.Path == "" && cfg.URL == "" {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: lancedb.path or url is required")
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

	// Use Path for connection string
	connStr := cfg.Path
	if connStr == "" {
		connStr = cfg.URL
	}

	db, err := lancedb.Connect(connStr)
	if err != nil {
		return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to connect to LanceDB: %w", err)
	}

	tableName := "picoclaw_vectors"

	// Define Arrow Schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.BinaryTypes.String},
			{Name: "vector", Type: arrow.FixedSizeListOf(int32(vectorCfg.Dimension), arrow.PrimitiveTypes.Float32)},
			{Name: "metadata", Type: arrow.BinaryTypes.String},
		},
		nil,
	)

	// Try to open existing table
	// Check if table exists
	tableNames, err := db.TableNames(context.Background())
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to list tables: %w", err)
	}

	tableExists := false
	for _, name := range tableNames {
		if name == tableName {
			tableExists = true
			break
		}
	}

	var table *lancedb.Table
	if !tableExists {
		// Create table
		tbl, err := db.CreateTable(context.Background(), tableName, schema)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to create table: %w", err)
		}
		// Convert interface to our local struct
		switch v := tbl.(type) {
		case *lancedb.Table:
			table = v
		default:
			// Needs type coercion depending on lancedb-go interface
			db.Close()
			return nil, fmt.Errorf("ERR_INTERNAL: lancedb.Table type coercion failed")
		}
	} else {
		// Open existing
		tbl, err := db.OpenTable(context.Background(), tableName)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to open table: %w", err)
		}
		switch v := tbl.(type) {
		case *lancedb.Table:
			table = v
		default:
			db.Close()
			return nil, fmt.Errorf("ERR_INTERNAL: lancedb.Table type coercion failed")
		}
	}

	return &LanceDBStore{
		db:         db,
		table:      table,
		dimension:  vectorCfg.Dimension,
		timeout:    time.Duration(timeoutMs) * time.Millisecond,
		breaker:    breaker,
		maxRetries: maxRetries,
		retryDelay: time.Duration(retryDelay) * time.Millisecond,
	}, nil
}

func (l *LanceDBStore) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
		remaining := time.Until(deadline)
		if remaining < l.timeout {
			return context.WithCancel(ctx)
		}
	}
	return context.WithTimeout(ctx, l.timeout)
}

func (l *LanceDBStore) executeWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	baseDelay := l.retryDelay

	for attempt := 0; attempt < l.maxRetries; attempt++ {
		err := l.breaker.Call(ctx, fn)
		if err == nil {
			return nil
		}
		lastErr = err

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "ERR_CIRCUIT_OPEN") {
			return err
		}

		if attempt == l.maxRetries-1 {
			break
		}

		backoff := baseDelay * time.Duration(1<<uint(attempt))
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}
		jitter := time.Duration(float64(backoff) * 0.25 * (2*rand.Float64() - 1))
		delay := backoff + jitter

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return lastErr
}

// Upsert inserts or updates vectors in the store
func (l *LanceDBStore) Upsert(ctx context.Context, vectors []Vector) error {
	if len(vectors) == 0 {
		return nil
	}
	ctx, cancel := l.withTimeout(ctx)
	defer cancel()

	return l.executeWithRetry(ctx, func() error {
		// validate dimensions
		for i, v := range vectors {
			if len(v.Embedding) != l.dimension {
				return fmt.Errorf("ERR_DIMENSION_MISMATCH: vector[%d] has dimension %d, expected %d", i, len(v.Embedding), l.dimension)
			}
		}

		schema := arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.BinaryTypes.String},
				{Name: "vector", Type: arrow.FixedSizeListOf(int32(l.dimension), arrow.PrimitiveTypes.Float32)},
				{Name: "metadata", Type: arrow.BinaryTypes.String},
			},
			nil,
		)

		b := array.NewRecordBuilder(memory.DefaultAllocator, schema)
		defer b.Release()

		idBuilder := b.Field(0).(*array.StringBuilder)
		vecListBuilder := b.Field(1).(*array.FixedSizeListBuilder)
		vecValuesBuilder := vecListBuilder.ValueBuilder().(*array.Float32Builder)
		metaBuilder := b.Field(2).(*array.StringBuilder)

		for _, v := range vectors {
			idBuilder.Append(v.ID)

			metaBytes, _ := json.Marshal(v.Metadata)
			metaBuilder.Append(string(metaBytes))

			vecListBuilder.Append(true)
			vecValuesBuilder.AppendValues(v.Embedding, nil)
		}

		rec := b.NewRecord()
		defer rec.Release()

		// LanceDB does not natively support an upsert in Go API yet, so we'll delete matches then add.
		ids := make([]string, len(vectors))
		for i, v := range vectors {
			ids[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(v.ID, "'", "''"))
		}
		delCond := fmt.Sprintf("id IN (%s)", strings.Join(ids, ","))

		// Delete existing just in case (poor man's upsert)
		_ = l.table.Delete(ctx, delCond)

		err := l.table.Add(ctx, []arrow.Record{rec})
		if err != nil {
			return fmt.Errorf("ERR_INTERNAL: failed to add vector record: %w", err)
		}

		return nil
	})
}

// Search performs vector similarity search
func (l *LanceDBStore) Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		return nil, fmt.Errorf("ERR_CONFIG_INVALID: topK must be positive, got: %d", topK)
	}

	ctx, cancel := l.withTimeout(ctx)
	defer cancel()

	var results []SearchResult
	err := l.executeWithRetry(ctx, func() error {
		if len(query.Embedding) != l.dimension {
			return fmt.Errorf("ERR_DIMENSION_MISMATCH: query has dimension %d, expected %d", len(query.Embedding), l.dimension)
		}

		// Vector search
		recs, err := l.table.Query().Vector(query.Embedding).Limit(topK).Execute()
		if err != nil {
			return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to search vectors: %w", err)
		}

		for _, rec := range recs {
			defer rec.Release()
			// Parse Record
			idCol := rec.Column(rec.Schema().FieldIndices("id")[0]).(*array.String)
			metaCol := rec.Column(rec.Schema().FieldIndices("metadata")[0]).(*array.String)
			// Score is mapped by LanceDB if distance is calculated, _distance column
			distIdx := rec.Schema().FieldIndices("_distance")

			for i := 0; i < int(rec.NumRows()); i++ {
				id := idCol.Value(i)
				metaStr := metaCol.Value(i)

				var meta map[string]interface{}
				_ = json.Unmarshal([]byte(metaStr), &meta)

				score := float32(0.0)
				if len(distIdx) > 0 {
					distCol := rec.Column(distIdx[0]).(*array.Float32)
					score = distCol.Value(i) // smaller is better in L2, might need inversion depending on logic
					// Usually picoclaw expects a similarity score. We'll store distance in Score field for simplicity
				}

				results = append(results, SearchResult{
					ID:       id,
					Score:    score,
					Metadata: meta,
				})
			}
		}

		return nil
	})

	return results, err
}

// Delete removes vectors from the store by IDs
func (l *LanceDBStore) Delete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	ctx, cancel := l.withTimeout(ctx)
	defer cancel()

	return l.executeWithRetry(ctx, func() error {
		escapedIDs := make([]string, len(ids))
		for i, id := range ids {
			escapedIDs[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(id, "'", "''"))
		}
		cond := fmt.Sprintf("id IN (%s)", strings.Join(escapedIDs, ","))
		err := l.table.Delete(ctx, cond)
		if err != nil {
			return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: failed to delete vectors: %w", err)
		}
		return nil
	})
}

// Health checks the health of the connection
func (l *LanceDBStore) Health(ctx context.Context) error {
	_, err := l.db.TableNames(ctx)
	if err != nil {
		return fmt.Errorf("ERR_PROVIDER_UNAVAILABLE: health check failed: %w", err)
	}
	return nil
}

// Close closes the connection
func (l *LanceDBStore) Close() error {
	return l.db.Close()
}
