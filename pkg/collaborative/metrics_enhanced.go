// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"sync"
	"sync/atomic"
	"time"
)

// EnhancedMetrics provides comprehensive metrics for testing and monitoring
type EnhancedMetrics struct {
	// Idempotency metrics
	idempotencyHitCount   int64
	idempotencyMissCount  int64
	duplicateCreatedCount int64
	conflictCount         int64

	// Circuit breaker metrics
	circuitBreakerState         string
	circuitBreakerOpenTotal     int64
	circuitBreakerRejectTotal   int64
	circuitBreakerRecoveryCount int64
	lastRecoveryDuration        time.Duration

	// Depth metrics
	depthRejectedTotal int64
	depthDistribution  map[int]int64
	depthMu            sync.RWMutex

	// Memory metrics
	memoryUsageBytes     int64
	memoryPeakBytes      int64
	gcPauseMs            int64
	gcCycles             int64

	// Request metrics
	requestInflight int64

	// TTL cleanup metrics
	idempotencyKeyTTLExpiredTotal int64

	mu sync.RWMutex
}

// NewEnhancedMetrics creates a new enhanced metrics collector
func NewEnhancedMetrics() *EnhancedMetrics {
	return &EnhancedMetrics{
		depthDistribution:   make(map[int]int64),
		circuitBreakerState: "closed",
	}
}

// Idempotency Metrics

func (em *EnhancedMetrics) RecordIdempotencyHit() {
	atomic.AddInt64(&em.idempotencyHitCount, 1)
}

func (em *EnhancedMetrics) RecordIdempotencyMiss() {
	atomic.AddInt64(&em.idempotencyMissCount, 1)
}

func (em *EnhancedMetrics) RecordDuplicateCreated() {
	atomic.AddInt64(&em.duplicateCreatedCount, 1)
}

func (em *EnhancedMetrics) RecordConflict() {
	atomic.AddInt64(&em.conflictCount, 1)
}

func (em *EnhancedMetrics) GetIdempotencyHitRatio() float64 {
	hits := atomic.LoadInt64(&em.idempotencyHitCount)
	misses := atomic.LoadInt64(&em.idempotencyMissCount)
	total := hits + misses
	if total == 0 {
		return 0.0
	}
	return float64(hits) / float64(total)
}

func (em *EnhancedMetrics) GetDuplicateCreatedCount() int64 {
	return atomic.LoadInt64(&em.duplicateCreatedCount)
}

func (em *EnhancedMetrics) GetConflictCount() int64 {
	return atomic.LoadInt64(&em.conflictCount)
}

// Circuit Breaker Metrics

func (em *EnhancedMetrics) SetCircuitBreakerState(state string) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.circuitBreakerState = state
}

func (em *EnhancedMetrics) GetCircuitBreakerState() string {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.circuitBreakerState
}

func (em *EnhancedMetrics) RecordCircuitBreakerOpen() {
	atomic.AddInt64(&em.circuitBreakerOpenTotal, 1)
}

func (em *EnhancedMetrics) RecordCircuitBreakerReject() {
	atomic.AddInt64(&em.circuitBreakerRejectTotal, 1)
}

func (em *EnhancedMetrics) RecordCircuitBreakerRecovery(duration time.Duration) {
	atomic.AddInt64(&em.circuitBreakerRecoveryCount, 1)
	em.mu.Lock()
	em.lastRecoveryDuration = duration
	em.mu.Unlock()
}

func (em *EnhancedMetrics) GetCircuitBreakerOpenTotal() int64 {
	return atomic.LoadInt64(&em.circuitBreakerOpenTotal)
}

func (em *EnhancedMetrics) GetCircuitBreakerRejectTotal() int64 {
	return atomic.LoadInt64(&em.circuitBreakerRejectTotal)
}

func (em *EnhancedMetrics) GetCircuitBreakerRecoveryDuration() time.Duration {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.lastRecoveryDuration
}

// Depth Metrics

func (em *EnhancedMetrics) RecordDepthRejected() {
	atomic.AddInt64(&em.depthRejectedTotal, 1)
}

func (em *EnhancedMetrics) RecordDepth(depth int) {
	em.depthMu.Lock()
	defer em.depthMu.Unlock()
	em.depthDistribution[depth]++
}

func (em *EnhancedMetrics) GetDepthRejectedTotal() int64 {
	return atomic.LoadInt64(&em.depthRejectedTotal)
}

func (em *EnhancedMetrics) GetDepthDistribution() map[int]int64 {
	em.depthMu.RLock()
	defer em.depthMu.RUnlock()
	
	result := make(map[int]int64)
	for k, v := range em.depthDistribution {
		result[k] = v
	}
	return result
}

// Memory Metrics

func (em *EnhancedMetrics) SetMemoryUsage(bytes int64) {
	atomic.StoreInt64(&em.memoryUsageBytes, bytes)
	
	// Update peak if necessary
	for {
		peak := atomic.LoadInt64(&em.memoryPeakBytes)
		if bytes <= peak {
			break
		}
		if atomic.CompareAndSwapInt64(&em.memoryPeakBytes, peak, bytes) {
			break
		}
	}
}

func (em *EnhancedMetrics) GetMemoryUsageBytes() int64 {
	return atomic.LoadInt64(&em.memoryUsageBytes)
}

func (em *EnhancedMetrics) GetMemoryPeakBytes() int64 {
	return atomic.LoadInt64(&em.memoryPeakBytes)
}

func (em *EnhancedMetrics) RecordGCPause(pauseMs int64) {
	atomic.AddInt64(&em.gcPauseMs, pauseMs)
	atomic.AddInt64(&em.gcCycles, 1)
}

func (em *EnhancedMetrics) GetGCStats() (totalPauseMs int64, cycles int64) {
	return atomic.LoadInt64(&em.gcPauseMs), atomic.LoadInt64(&em.gcCycles)
}

// Request Metrics

func (em *EnhancedMetrics) IncrementInflight() {
	atomic.AddInt64(&em.requestInflight, 1)
}

func (em *EnhancedMetrics) DecrementInflight() {
	atomic.AddInt64(&em.requestInflight, -1)
}

func (em *EnhancedMetrics) GetRequestInflight() int64 {
	return atomic.LoadInt64(&em.requestInflight)
}

// TTL Cleanup Metrics

func (em *EnhancedMetrics) RecordIdempotencyKeyExpired() {
	atomic.AddInt64(&em.idempotencyKeyTTLExpiredTotal, 1)
}

func (em *EnhancedMetrics) GetIdempotencyKeyTTLExpiredTotal() int64 {
	return atomic.LoadInt64(&em.idempotencyKeyTTLExpiredTotal)
}

// Snapshot returns a complete snapshot of all metrics
func (em *EnhancedMetrics) Snapshot() MetricsSnapshot {
	em.mu.RLock()
	em.depthMu.RLock()
	defer em.mu.RUnlock()
	defer em.depthMu.RUnlock()

	depthDist := make(map[int]int64)
	for k, v := range em.depthDistribution {
		depthDist[k] = v
	}

	return MetricsSnapshot{
		// Idempotency
		IdempotencyHitCount:   atomic.LoadInt64(&em.idempotencyHitCount),
		IdempotencyMissCount:  atomic.LoadInt64(&em.idempotencyMissCount),
		IdempotencyHitRatio:   em.GetIdempotencyHitRatio(),
		DuplicateCreatedCount: atomic.LoadInt64(&em.duplicateCreatedCount),
		ConflictCount:         atomic.LoadInt64(&em.conflictCount),
		
		// Circuit Breaker
		CircuitBreakerState:         em.circuitBreakerState,
		CircuitBreakerOpenTotal:     atomic.LoadInt64(&em.circuitBreakerOpenTotal),
		CircuitBreakerRejectTotal:   atomic.LoadInt64(&em.circuitBreakerRejectTotal),
		CircuitBreakerRecoveryCount: atomic.LoadInt64(&em.circuitBreakerRecoveryCount),
		LastRecoveryDuration:        em.lastRecoveryDuration,
		
		// Depth
		DepthRejectedTotal: atomic.LoadInt64(&em.depthRejectedTotal),
		DepthDistribution:  depthDist,
		
		// Memory
		MemoryUsageBytes: atomic.LoadInt64(&em.memoryUsageBytes),
		MemoryPeakBytes:  atomic.LoadInt64(&em.memoryPeakBytes),
		GCPauseMs:        atomic.LoadInt64(&em.gcPauseMs),
		GCCycles:         atomic.LoadInt64(&em.gcCycles),
		
		// Request
		RequestInflight: atomic.LoadInt64(&em.requestInflight),
		
		// TTL
		IdempotencyKeyTTLExpiredTotal: atomic.LoadInt64(&em.idempotencyKeyTTLExpiredTotal),
	}
}

// MetricsSnapshot represents a point-in-time snapshot of all metrics
type MetricsSnapshot struct {
	// Idempotency
	IdempotencyHitCount   int64
	IdempotencyMissCount  int64
	IdempotencyHitRatio   float64
	DuplicateCreatedCount int64
	ConflictCount         int64
	
	// Circuit Breaker
	CircuitBreakerState         string
	CircuitBreakerOpenTotal     int64
	CircuitBreakerRejectTotal   int64
	CircuitBreakerRecoveryCount int64
	LastRecoveryDuration        time.Duration
	
	// Depth
	DepthRejectedTotal int64
	DepthDistribution  map[int]int64
	
	// Memory
	MemoryUsageBytes int64
	MemoryPeakBytes  int64
	GCPauseMs        int64
	GCCycles         int64
	
	// Request
	RequestInflight int64
	
	// TTL
	IdempotencyKeyTTLExpiredTotal int64
}

// Reset clears all metrics (useful for testing)
func (em *EnhancedMetrics) Reset() {
	atomic.StoreInt64(&em.idempotencyHitCount, 0)
	atomic.StoreInt64(&em.idempotencyMissCount, 0)
	atomic.StoreInt64(&em.duplicateCreatedCount, 0)
	atomic.StoreInt64(&em.conflictCount, 0)
	
	atomic.StoreInt64(&em.circuitBreakerOpenTotal, 0)
	atomic.StoreInt64(&em.circuitBreakerRejectTotal, 0)
	atomic.StoreInt64(&em.circuitBreakerRecoveryCount, 0)
	
	atomic.StoreInt64(&em.depthRejectedTotal, 0)
	
	atomic.StoreInt64(&em.memoryUsageBytes, 0)
	atomic.StoreInt64(&em.memoryPeakBytes, 0)
	atomic.StoreInt64(&em.gcPauseMs, 0)
	atomic.StoreInt64(&em.gcCycles, 0)
	
	atomic.StoreInt64(&em.requestInflight, 0)
	atomic.StoreInt64(&em.idempotencyKeyTTLExpiredTotal, 0)
	
	em.mu.Lock()
	em.circuitBreakerState = "closed"
	em.lastRecoveryDuration = 0
	em.mu.Unlock()
	
	em.depthMu.Lock()
	em.depthDistribution = make(map[int]int64)
	em.depthMu.Unlock()
}
