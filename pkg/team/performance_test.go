// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"testing"

	"github.com/sipeed/picoclaw/pkg/agent"
)

// TestAgentPool tests agent instance pooling
func TestAgentPool(t *testing.T) {
	pool := NewAgentPool()

	// Create a mock agent instance
	createFn := func() *agent.AgentInstance {
		return &agent.AgentInstance{}
	}

	// Get first instance
	instance1 := pool.GetOrCreateInstance("developer", createFn)
	if instance1 == nil {
		t.Fatal("Expected instance, got nil")
	}

	// Get same instance again
	instance2 := pool.GetOrCreateInstance("developer", createFn)
	if instance1 != instance2 {
		t.Error("Expected same instance to be reused")
	}

	// Check usage count
	if count := pool.GetUsageCount("developer"); count != 2 {
		t.Errorf("Expected usage count 2, got %d", count)
	}

	// Release instance
	pool.ReleaseInstance("developer")
	if count := pool.GetUsageCount("developer"); count != 1 {
		t.Errorf("Expected usage count 1 after release, got %d", count)
	}

	// Release again should remove from pool
	pool.ReleaseInstance("developer")
	if count := pool.GetUsageCount("developer"); count != 0 {
		t.Errorf("Expected usage count 0 after final release, got %d", count)
	}
}

// TestAgentPoolClear tests clearing the agent pool
func TestAgentPoolClear(t *testing.T) {
	pool := NewAgentPool()

	createFn := func() *agent.AgentInstance {
		return &agent.AgentInstance{}
	}

	// Add multiple instances
	pool.GetOrCreateInstance("developer", createFn)
	pool.GetOrCreateInstance("tester", createFn)

	// Clear pool
	pool.Clear()

	// Verify pool is empty
	if count := pool.GetUsageCount("developer"); count != 0 {
		t.Errorf("Expected usage count 0 after clear, got %d", count)
	}
	if count := pool.GetUsageCount("tester"); count != 0 {
		t.Errorf("Expected usage count 0 after clear, got %d", count)
	}
}

// TestRoleCache tests role-to-capability caching
func TestRoleCache(t *testing.T) {
	cache := NewRoleCache()

	// Set capabilities
	caps := []string{"code", "test", "review"}
	cache.Set("developer", caps)

	// Get capabilities
	retrieved, exists := cache.Get("developer")
	if !exists {
		t.Fatal("Expected capabilities to exist in cache")
	}

	if len(retrieved) != len(caps) {
		t.Errorf("Expected %d capabilities, got %d", len(caps), len(retrieved))
	}

	for i, cap := range caps {
		if retrieved[i] != cap {
			t.Errorf("Expected capability %s, got %s", cap, retrieved[i])
		}
	}
}

// TestRoleCacheInvalidate tests cache invalidation
func TestRoleCacheInvalidate(t *testing.T) {
	cache := NewRoleCache()

	// Set capabilities
	caps := []string{"code", "test"}
	cache.Set("developer", caps)

	// Invalidate
	cache.Invalidate("developer")

	// Verify removed
	_, exists := cache.Get("developer")
	if exists {
		t.Error("Expected capabilities to be removed from cache")
	}
}

// TestRoleCacheClear tests clearing the cache
func TestRoleCacheClear(t *testing.T) {
	cache := NewRoleCache()

	// Set multiple roles
	cache.Set("developer", []string{"code"})
	cache.Set("tester", []string{"test"})

	// Clear cache
	cache.Clear()

	// Verify all removed
	_, exists1 := cache.Get("developer")
	_, exists2 := cache.Get("tester")

	if exists1 || exists2 {
		t.Error("Expected all entries to be removed from cache")
	}
}

// TestGetCapabilitiesForRole tests cached capability retrieval
func TestGetCapabilitiesForRole(t *testing.T) {
	tm := &TeamManager{
		roleCapabilities: map[string][]string{
			"developer": {"code", "test"},
		},
		roleCache: NewRoleCache(),
	}

	// First call should populate cache
	caps1, exists1 := tm.GetCapabilitiesForRole("developer")
	if !exists1 {
		t.Fatal("Expected capabilities to exist")
	}

	// Second call should use cache
	caps2, exists2 := tm.GetCapabilitiesForRole("developer")
	if !exists2 {
		t.Fatal("Expected capabilities to exist in cache")
	}

	if len(caps1) != len(caps2) {
		t.Error("Expected same capabilities from cache")
	}
}

// TestInvalidateRoleCache tests cache invalidation through TeamManager
func TestInvalidateRoleCache(t *testing.T) {
	tm := &TeamManager{
		roleCapabilities: map[string][]string{
			"developer": {"code", "test"},
		},
		roleCache: NewRoleCache(),
	}

	// Populate cache
	tm.GetCapabilitiesForRole("developer")

	// Invalidate
	tm.InvalidateRoleCache("developer")

	// Verify cache is empty but roleCapabilities still has data
	_, cached := tm.roleCache.Get("developer")
	if cached {
		t.Error("Expected cache to be invalidated")
	}

	// Should still be able to retrieve from roleCapabilities
	caps, exists := tm.GetCapabilitiesForRole("developer")
	if !exists || len(caps) != 2 {
		t.Error("Expected to retrieve capabilities from roleCapabilities map")
	}
}

// TestAgentPoolConcurrency tests concurrent access to agent pool
func TestAgentPoolConcurrency(t *testing.T) {
	pool := NewAgentPool()
	createFn := func() *agent.AgentInstance {
		return &agent.AgentInstance{}
	}

	// Concurrent gets
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			instance := pool.GetOrCreateInstance("developer", createFn)
			if instance == nil {
				t.Error("Expected instance, got nil")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have usage count of 10
	if count := pool.GetUsageCount("developer"); count != 10 {
		t.Errorf("Expected usage count 10, got %d", count)
	}
}

// TestRoleCacheConcurrency tests concurrent access to role cache
func TestRoleCacheConcurrency(t *testing.T) {
	cache := NewRoleCache()

	// Concurrent sets and gets
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			role := "developer"
			caps := []string{"code", "test"}
			cache.Set(role, caps)
			_, exists := cache.Get(role)
			if !exists {
				t.Error("Expected capabilities to exist")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
