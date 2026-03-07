// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/providers/protocoltypes"
)

// MockLLMProvider for testing
type MockLLMProvider struct {
	shouldFail   bool
	delay        time.Duration
	responseText string
	callCount    int
	lastMessages []protocoltypes.Message
	lastModel    string
	lastOptions  map[string]any
}

func (m *MockLLMProvider) Chat(
	ctx context.Context,
	messages []protocoltypes.Message,
	tools []protocoltypes.ToolDefinition,
	model string,
	options map[string]any,
) (*protocoltypes.LLMResponse, error) {
	m.callCount++
	m.lastMessages = messages
	m.lastModel = model
	m.lastOptions = options

	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if m.shouldFail {
		return nil, errors.New("mock LLM error")
	}

	response := m.responseText
	if response == "" {
		response = "## Project Overview\nTest project\n\n## Key Decisions\n- Decision 1\n- Decision 2"
	}

	return &protocoltypes.LLMResponse{
		Content:      response,
		FinishReason: "stop",
	}, nil
}

func (m *MockLLMProvider) GetDefaultModel() string {
	return "gpt-4o-mini"
}

func TestNewLLMSummarizer(t *testing.T) {
	config := CompactionConfig{
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
		SummaryMaxLength: 2000,
	}
	provider := &MockLLMProvider{}

	summarizer := NewLLMSummarizer(config, provider)

	if summarizer == nil {
		t.Fatal("Expected summarizer to be created")
	}
	if summarizer.config.LLMModel != "gpt-4o-mini" {
		t.Errorf("Expected model 'gpt-4o-mini', got '%s'", summarizer.config.LLMModel)
	}
}

func TestBuildPrompt_WithoutExistingSummary(t *testing.T) {
	config := CompactionConfig{
		SummaryMaxLength: 2000,
	}
	summarizer := NewLLMSummarizer(config, &MockLLMProvider{})

	messages := []Message{
		{Role: "user", Content: "Let's build a REST API", Timestamp: time.Now()},
		{Role: "architect", Content: "I'll design it", Timestamp: time.Now()},
	}

	req := &CompactionRequest{
		SessionID:       "test123",
		Messages:        messages,
		ExistingSummary: "",
		Config:          config,
	}

	prompt := summarizer.buildPrompt(req)

	// Should contain task description
	if !strings.Contains(prompt, "context summarizer") {
		t.Error("Expected task description in prompt")
	}

	// Should contain messages
	if !strings.Contains(prompt, "REST API") {
		t.Error("Expected message content in prompt")
	}

	// Should NOT contain previous summary section
	if strings.Contains(prompt, "=== Previous Summary ===") {
		t.Error("Should not contain previous summary section")
	}

	// Should contain format instructions
	if !strings.Contains(prompt, "## Project Overview") {
		t.Error("Expected format instructions")
	}

	// Should contain max length
	if !strings.Contains(prompt, "2000 characters") {
		t.Error("Expected max length in prompt")
	}
}

func TestBuildPrompt_WithExistingSummary(t *testing.T) {
	config := CompactionConfig{
		SummaryMaxLength: 2000,
	}
	summarizer := NewLLMSummarizer(config, &MockLLMProvider{})

	messages := []Message{
		{Role: "user", Content: "Add authentication", Timestamp: time.Now()},
	}

	req := &CompactionRequest{
		SessionID:       "test123",
		Messages:        messages,
		ExistingSummary: "Previous: Building REST API",
		Config:          config,
	}

	prompt := summarizer.buildPrompt(req)

	// Should contain previous summary
	if !strings.Contains(prompt, "=== Previous Summary ===") {
		t.Error("Expected previous summary section")
	}
	if !strings.Contains(prompt, "Previous: Building REST API") {
		t.Error("Expected previous summary content")
	}

	// Should contain new messages
	if !strings.Contains(prompt, "Add authentication") {
		t.Error("Expected new message content")
	}
}

func TestCallLLM_Success(t *testing.T) {
	config := CompactionConfig{
		LLMModel:      "gpt-4o-mini",
		LLMTimeout:    30 * time.Second,
		LLMMaxRetries: 3,
	}
	provider := &MockLLMProvider{
		responseText: "Test summary content",
	}
	summarizer := NewLLMSummarizer(config, provider)

	ctx := context.Background()
	summary, err := summarizer.callLLM(ctx, "test prompt")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if summary != "Test summary content" {
		t.Errorf("Expected 'Test summary content', got '%s'", summary)
	}

	// Verify provider was called correctly
	if provider.callCount != 1 {
		t.Errorf("Expected 1 call, got %d", provider.callCount)
	}
	if len(provider.lastMessages) != 2 {
		t.Errorf("Expected 2 messages (system + user), got %d", len(provider.lastMessages))
	}
	if provider.lastMessages[0].Role != "system" {
		t.Errorf("Expected first message role 'system', got '%s'", provider.lastMessages[0].Role)
	}
	if provider.lastMessages[1].Role != "user" {
		t.Errorf("Expected second message role 'user', got '%s'", provider.lastMessages[1].Role)
	}
	if provider.lastModel != "gpt-4o-mini" {
		t.Errorf("Expected model 'gpt-4o-mini', got '%s'", provider.lastModel)
	}

	// Check options
	if temp, ok := provider.lastOptions["temperature"].(float64); !ok || temp != 0.3 {
		t.Errorf("Expected temperature 0.3, got %v", provider.lastOptions["temperature"])
	}
}

func TestCallLLM_Timeout(t *testing.T) {
	config := CompactionConfig{
		LLMModel:      "gpt-4o-mini",
		LLMTimeout:    50 * time.Millisecond, // Short timeout
		LLMMaxRetries: 0,
	}
	provider := &MockLLMProvider{
		delay: 200 * time.Millisecond, // Longer than timeout
	}
	summarizer := NewLLMSummarizer(config, provider)

	ctx := context.Background()
	_, err := summarizer.callLLM(ctx, "test prompt")

	if err == nil {
		t.Error("Expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got %v", err)
	}
}

func TestCallLLMWithRetry_Success(t *testing.T) {
	config := CompactionConfig{
		LLMMaxRetries: 3,
	}
	provider := &MockLLMProvider{
		responseText: "Success",
	}
	summarizer := NewLLMSummarizer(config, provider)

	ctx := context.Background()
	summary, err := summarizer.callLLMWithRetry(ctx, "test prompt")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if summary != "Success" {
		t.Errorf("Expected 'Success', got '%s'", summary)
	}
	if provider.callCount != 1 {
		t.Errorf("Expected 1 call, got %d", provider.callCount)
	}
}

func TestCallLLMWithRetry_FailThenSuccess(t *testing.T) {
	config := CompactionConfig{
		LLMMaxRetries: 3,
		LLMTimeout:    30 * time.Second,
	}
	provider := &MockLLMProvider{
		shouldFail: true,
	}
	summarizer := NewLLMSummarizer(config, provider)

	ctx := context.Background()

	// First call fails
	_, err := summarizer.callLLMWithRetry(ctx, "test prompt")
	if err == nil {
		t.Error("Expected error on first call")
	}

	// Should have tried 4 times (initial + 3 retries)
	if provider.callCount != 4 {
		t.Errorf("Expected 4 calls, got %d", provider.callCount)
	}

	// Now make it succeed
	provider.shouldFail = false
	provider.callCount = 0

	summary, err := summarizer.callLLMWithRetry(ctx, "test prompt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
}

func TestValidateAndTruncate_Normal(t *testing.T) {
	config := CompactionConfig{
		SummaryMaxLength: 2000,
	}
	summarizer := NewLLMSummarizer(config, &MockLLMProvider{})

	summary := "  Test summary with spaces  "
	result := summarizer.validateAndTruncate(summary)

	if result != "Test summary with spaces" {
		t.Errorf("Expected trimmed summary, got '%s'", result)
	}
}

func TestValidateAndTruncate_TooLong(t *testing.T) {
	config := CompactionConfig{
		SummaryMaxLength: 50,
	}
	summarizer := NewLLMSummarizer(config, &MockLLMProvider{})

	summary := "This is a very long summary that exceeds the maximum length. It should be truncated."
	result := summarizer.validateAndTruncate(summary)

	if len(result) > 50 {
		t.Errorf("Expected length <= 50, got %d", len(result))
	}
}

func TestValidateAndTruncate_SentenceBoundary(t *testing.T) {
	config := CompactionConfig{
		SummaryMaxLength: 60,
	}
	summarizer := NewLLMSummarizer(config, &MockLLMProvider{})

	summary := "First sentence. Second sentence. Third sentence that is very long."
	result := summarizer.validateAndTruncate(summary)

	// Should cut at sentence boundary
	if !strings.HasSuffix(result, ".") {
		t.Error("Expected to cut at sentence boundary")
	}
	if len(result) > 60 {
		t.Errorf("Expected length <= 60, got %d", len(result))
	}
}

func TestCalculateSize(t *testing.T) {
	summarizer := NewLLMSummarizer(CompactionConfig{}, &MockLLMProvider{})

	messages := []Message{
		{Role: "user", Content: "Hello"},      // 5 + 4 + 50 = 59
		{Role: "architect", Content: "World"}, // 5 + 9 + 50 = 64
	}

	size := summarizer.calculateSize(messages)

	expected := 59 + 64
	if size != expected {
		t.Errorf("Expected size %d, got %d", expected, size)
	}
}

func TestSummarize_Success(t *testing.T) {
	config := CompactionConfig{
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
		SummaryMaxLength: 2000,
	}
	provider := &MockLLMProvider{
		responseText: "## Project Overview\nREST API project\n\n## Key Decisions\n- Use JWT auth",
	}
	summarizer := NewLLMSummarizer(config, provider)

	messages := []Message{
		{Role: "user", Content: "Build REST API", Timestamp: time.Now()},
		{Role: "architect", Content: "Use JWT", Timestamp: time.Now()},
	}

	req := &CompactionRequest{
		SessionID:       "test123",
		Messages:        messages,
		ExistingSummary: "",
		Config:          config,
		Timestamp:       time.Now(),
	}

	ctx := context.Background()
	result, err := summarizer.Summarize(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if !result.Success {
		t.Error("Expected success to be true")
	}
	if result.Summary == "" {
		t.Error("Expected non-empty summary")
	}
	if !strings.Contains(result.Summary, "REST API") {
		t.Error("Expected summary to contain 'REST API'")
	}
	if result.MessagesCount != 2 {
		t.Errorf("Expected 2 messages, got %d", result.MessagesCount)
	}
	if result.OriginalSize <= 0 {
		t.Error("Expected positive original size")
	}
	if result.CompressedSize <= 0 {
		t.Error("Expected positive compressed size")
	}
}

func TestSummarize_Failure(t *testing.T) {
	config := CompactionConfig{
		LLMModel:      "gpt-4o-mini",
		LLMTimeout:    30 * time.Second,
		LLMMaxRetries: 2,
	}
	provider := &MockLLMProvider{
		shouldFail: true,
	}
	summarizer := NewLLMSummarizer(config, provider)

	messages := []Message{
		{Role: "user", Content: "Test", Timestamp: time.Now()},
	}

	req := &CompactionRequest{
		SessionID: "test123",
		Messages:  messages,
		Config:    config,
		Timestamp: time.Now(),
	}

	ctx := context.Background()
	result, err := summarizer.Summarize(ctx, req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Error("Expected nil result on error")
	}
	if !strings.Contains(err.Error(), "LLM call failed") {
		t.Errorf("Expected 'LLM call failed' error, got %v", err)
	}
}

func TestSummarize_WithExistingSummary(t *testing.T) {
	config := CompactionConfig{
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
		SummaryMaxLength: 2000,
	}
	provider := &MockLLMProvider{
		responseText: "Updated summary with new info",
	}
	summarizer := NewLLMSummarizer(config, provider)

	messages := []Message{
		{Role: "user", Content: "Add feature X", Timestamp: time.Now()},
	}

	req := &CompactionRequest{
		SessionID:       "test123",
		Messages:        messages,
		ExistingSummary: "Previous summary about REST API",
		Config:          config,
		Timestamp:       time.Now(),
	}

	ctx := context.Background()
	result, err := summarizer.Summarize(ctx, req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify prompt included existing summary
	if provider.callCount != 1 {
		t.Errorf("Expected 1 call, got %d", provider.callCount)
	}

	// Check that user message contains existing summary
	userMsg := provider.lastMessages[1].Content
	if !strings.Contains(userMsg, "Previous summary about REST API") {
		t.Error("Expected prompt to include existing summary")
	}

	if result.Summary != "Updated summary with new info" {
		t.Errorf("Expected updated summary, got '%s'", result.Summary)
	}
}
