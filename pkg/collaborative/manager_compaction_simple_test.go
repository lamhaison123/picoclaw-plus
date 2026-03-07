package collaborative

import (
	"testing"
)

func TestCompactionConfigFields(t *testing.T) {
	provider := &MockLLMProvider{}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
	}

	t.Logf("CompactionEnabled: %v", config.CompactionEnabled)
	t.Logf("LLMProvider: %v", config.LLMProvider)
	t.Logf("LLMProvider != nil: %v", config.LLMProvider != nil)

	if !config.CompactionEnabled {
		t.Error("CompactionEnabled should be true")
	}

	if config.LLMProvider == nil {
		t.Error("LLMProvider should not be nil")
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	t.Logf("Manager created")
	t.Logf("GetCompactionManager: %v", manager.GetCompactionManager())

	if manager.GetCompactionManager() == nil {
		t.Error("Compaction manager should be initialized")
	}
}
