// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/providers"
)

// MockProvider implements providers.LLMProvider interface for testing
type MockProvider struct{}

func NewMockProvider() providers.LLMProvider {
	return &MockProvider{}
}

func (m *MockProvider) Chat(
	ctx context.Context,
	messages []providers.Message,
	tools []providers.ToolDefinition,
	model string,
	options map[string]any,
) (*providers.LLMResponse, error) {
	return &providers.LLMResponse{
		Content: "mock response",
	}, nil
}

func (m *MockProvider) GetDefaultModel() string {
	return "mock-model"
}
