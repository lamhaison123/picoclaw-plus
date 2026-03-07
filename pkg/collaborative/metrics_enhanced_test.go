// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"testing"
	"time"
)

func TestEnhancedMetrics_Idempotency(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Record some hits and misses
	metrics.RecordIdempotencyHit()
	metrics.RecordIdempotencyHit()
	metrics.RecordIdempotencyMiss()

	if metrics.GetIdempotencyHitRatio() != 2.0/3.0 {
		t.Errorf("Expected hit ratio 0.666, got %f", metrics.GetIdempotencyHitRatio())
	}

	// Record duplicate created (should be 0 in production)
	metrics.RecordDuplicateCreated()
	if metrics.GetDuplicateCreatedCount() != 1 {
		t.Errorf("Expected duplicate count 1, got %d", metrics.GetDuplicateCreatedCount())
	}

	// Record conflict
	metrics.RecordConflict()
	if metrics.GetConflictCount() != 1 {
		t.Errorf("Expected conflict count 1, got %d", metrics.GetConflictCount())
	}
}

func TestEnhancedMetrics_CircuitBreaker(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Initial state should be closed
	if metrics.GetCircuitBreakerState() != "closed" {
		t.Errorf("Expected initial state 'closed', got %s", metrics.GetCircuitBreakerState())
	}

	// Record state changes
	metrics.SetCircuitBreakerState("open")
	metrics.RecordCircuitBreakerOpen()

	if metrics.GetCircuitBreakerState() != "open" {
		t.Errorf("Expected state 'open', got %s", metrics.GetCircuitBreakerState())
	}

	if metrics.GetCircuitBreakerOpenTotal() != 1 {
		t.Errorf("Expected open count 1, got %d", metrics.GetCircuitBreakerOpenTotal())
	}

	// Record rejects
	metrics.RecordCircuitBreakerReject()
	metrics.RecordCircuitBreakerReject()

	if metrics.GetCircuitBreakerRejectTotal() != 2 {
		t.Errorf("Expected reject count 2, got %d", metrics.GetCircuitBreakerRejectTotal())
	}

	// Record recovery
	recoveryDuration := 5 * time.Second
	metrics.RecordCircuitBreakerRecovery(recoveryDuration)

	if metrics.GetCircuitBreakerRecoveryDuration() != recoveryDuration {
		t.Errorf("Expected recovery duration %v, got %v", recoveryDuration, metrics.GetCircuitBreakerRecoveryDuration())
	}
}

func TestEnhancedMetrics_Depth(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Record various depths
	metrics.RecordDepth(0)
	metrics.RecordDepth(1)
	metrics.RecordDepth(1)
	metrics.RecordDepth(2)
	metrics.RecordDepth(20)

	dist := metrics.GetDepthDistribution()
	if dist[0] != 1 {
		t.Errorf("Expected depth 0 count 1, got %d", dist[0])
	}
	if dist[1] != 2 {
		t.Errorf("Expected depth 1 count 2, got %d", dist[1])
	}
	if dist[2] != 1 {
		t.Errorf("Expected depth 2 count 1, got %d", dist[2])
	}
	if dist[20] != 1 {
		t.Errorf("Expected depth 20 count 1, got %d", dist[20])
	}

	// Record depth rejection
	metrics.RecordDepthRejected()
	if metrics.GetDepthRejectedTotal() != 1 {
		t.Errorf("Expected depth rejected count 1, got %d", metrics.GetDepthRejectedTotal())
	}
}

func TestEnhancedMetrics_Memory(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Set memory usage
	metrics.SetMemoryUsage(1024 * 1024) // 1MB
	if metrics.GetMemoryUsageBytes() != 1024*1024 {
		t.Errorf("Expected memory usage 1MB, got %d", metrics.GetMemoryUsageBytes())
	}

	// Set higher memory (should update peak)
	metrics.SetMemoryUsage(2 * 1024 * 1024) // 2MB
	if metrics.GetMemoryPeakBytes() != 2*1024*1024 {
		t.Errorf("Expected peak memory 2MB, got %d", metrics.GetMemoryPeakBytes())
	}

	// Set lower memory (peak should remain)
	metrics.SetMemoryUsage(512 * 1024) // 512KB
	if metrics.GetMemoryPeakBytes() != 2*1024*1024 {
		t.Errorf("Expected peak memory to remain 2MB, got %d", metrics.GetMemoryPeakBytes())
	}

	// Record GC pauses
	metrics.RecordGCPause(10)
	metrics.RecordGCPause(15)

	totalPause, cycles := metrics.GetGCStats()
	if totalPause != 25 {
		t.Errorf("Expected total GC pause 25ms, got %d", totalPause)
	}
	if cycles != 2 {
		t.Errorf("Expected 2 GC cycles, got %d", cycles)
	}
}

func TestEnhancedMetrics_RequestInflight(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Increment inflight
	metrics.IncrementInflight()
	metrics.IncrementInflight()

	if metrics.GetRequestInflight() != 2 {
		t.Errorf("Expected 2 inflight requests, got %d", metrics.GetRequestInflight())
	}

	// Decrement
	metrics.DecrementInflight()

	if metrics.GetRequestInflight() != 1 {
		t.Errorf("Expected 1 inflight request, got %d", metrics.GetRequestInflight())
	}
}

func TestEnhancedMetrics_TTLExpired(t *testing.T) {
	metrics := NewEnhancedMetrics()

	metrics.RecordIdempotencyKeyExpired()
	metrics.RecordIdempotencyKeyExpired()

	if metrics.GetIdempotencyKeyTTLExpiredTotal() != 2 {
		t.Errorf("Expected 2 expired keys, got %d", metrics.GetIdempotencyKeyTTLExpiredTotal())
	}
}

func TestEnhancedMetrics_Snapshot(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Record various metrics
	metrics.RecordIdempotencyHit()
	metrics.RecordIdempotencyMiss()
	metrics.SetCircuitBreakerState("open")
	metrics.RecordDepth(5)
	metrics.SetMemoryUsage(1024)
	metrics.IncrementInflight()

	// Get snapshot
	snapshot := metrics.Snapshot()

	// Verify snapshot
	if snapshot.IdempotencyHitCount != 1 {
		t.Errorf("Expected snapshot hit count 1, got %d", snapshot.IdempotencyHitCount)
	}
	if snapshot.IdempotencyMissCount != 1 {
		t.Errorf("Expected snapshot miss count 1, got %d", snapshot.IdempotencyMissCount)
	}
	if snapshot.CircuitBreakerState != "open" {
		t.Errorf("Expected snapshot state 'open', got %s", snapshot.CircuitBreakerState)
	}
	if snapshot.DepthDistribution[5] != 1 {
		t.Errorf("Expected depth 5 count 1, got %d", snapshot.DepthDistribution[5])
	}
	if snapshot.MemoryUsageBytes != 1024 {
		t.Errorf("Expected memory 1024, got %d", snapshot.MemoryUsageBytes)
	}
	if snapshot.RequestInflight != 1 {
		t.Errorf("Expected 1 inflight, got %d", snapshot.RequestInflight)
	}
}

func TestEnhancedMetrics_Reset(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Record various metrics
	metrics.RecordIdempotencyHit()
	metrics.RecordDuplicateCreated()
	metrics.SetCircuitBreakerState("open")
	metrics.RecordDepth(10)
	metrics.SetMemoryUsage(2048)

	// Reset
	metrics.Reset()

	// Verify all metrics are cleared
	snapshot := metrics.Snapshot()
	if snapshot.IdempotencyHitCount != 0 {
		t.Errorf("Expected hit count 0 after reset, got %d", snapshot.IdempotencyHitCount)
	}
	if snapshot.DuplicateCreatedCount != 0 {
		t.Errorf("Expected duplicate count 0 after reset, got %d", snapshot.DuplicateCreatedCount)
	}
	if snapshot.CircuitBreakerState != "closed" {
		t.Errorf("Expected state 'closed' after reset, got %s", snapshot.CircuitBreakerState)
	}
	if len(snapshot.DepthDistribution) != 0 {
		t.Errorf("Expected empty depth distribution after reset, got %d entries", len(snapshot.DepthDistribution))
	}
	if snapshot.MemoryUsageBytes != 0 {
		t.Errorf("Expected memory 0 after reset, got %d", snapshot.MemoryUsageBytes)
	}
}

func TestEnhancedMetrics_Concurrency(t *testing.T) {
	metrics := NewEnhancedMetrics()

	// Test concurrent access
	done := make(chan bool)

	// Goroutine 1: Record hits
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordIdempotencyHit()
		}
		done <- true
	}()

	// Goroutine 2: Record misses
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordIdempotencyMiss()
		}
		done <- true
	}()

	// Goroutine 3: Record depths
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordDepth(i % 10)
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	// Verify counts
	if metrics.GetIdempotencyHitRatio() != 0.5 {
		t.Errorf("Expected hit ratio 0.5, got %f", metrics.GetIdempotencyHitRatio())
	}

	dist := metrics.GetDepthDistribution()
	totalDepth := int64(0)
	for _, count := range dist {
		totalDepth += count
	}
	if totalDepth != 100 {
		t.Errorf("Expected 100 depth records, got %d", totalDepth)
	}
}
