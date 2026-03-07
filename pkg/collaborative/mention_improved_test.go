// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"strings"
	"testing"
)

func TestExtractMentionsImproved(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single mention",
			input:    "@developer please fix this",
			expected: []string{"developer"},
		},
		{
			name:     "Multiple mentions",
			input:    "@developer and @tester please review",
			expected: []string{"developer", "tester"},
		},
		{
			name:     "Mention with Vietnamese text",
			input:    "Cảm ơn @developer đã hoàn thành",
			expected: []string{"developer"},
		},
		{
			name:     "Mention in multi-line text",
			input:    "Line 1\n@architect please review\nLine 3",
			expected: []string{"architect"},
		},
		{
			name:     "Mention with code block",
			input:    "```\ncode here\n```\n@developer check this",
			expected: []string{"developer"},
		},
		{
			name:     "Email should be ignored",
			input:    "Contact user@gmail.com or @developer",
			expected: []string{"developer"},
		},
		{
			name:     "Duplicate mentions",
			input:    "@developer @developer @developer",
			expected: []string{"developer"},
		},
		{
			name:     "Mention at start of line",
			input:    "@manager\nPlease coordinate",
			expected: []string{"manager"},
		},
		{
			name:     "Mention after punctuation",
			input:    "Done! @tester please verify.",
			expected: []string{"tester"},
		},
		{
			name:     "Multiple mentions with Vietnamese",
			input:    "@manager ơi, nhờ bạn kêu @developer và @tester",
			expected: []string{"manager", "developer", "tester"},
		},
		{
			name:     "Mention in log block",
			input:    "Mar 06 14:14:37 picoclaw[45300]: @developer check logs",
			expected: []string{"developer"},
		},
		{
			name:     "No mentions",
			input:    "This is a regular message",
			expected: []string{},
		},
		{
			name:     "Mention with numbers",
			input:    "@developer123 please help",
			expected: []string{"developer123"},
		},
		{
			name:     "Invalid mention (starts with number)",
			input:    "@123developer should be ignored",
			expected: []string{},
		},
		{
			name:     "Mention with underscore",
			input:    "@dev_ops please deploy",
			expected: []string{"dev_ops"},
		},
		{
			name:     "Mixed case mentions",
			input:    "@Developer @TESTER @ArChItEcT",
			expected: []string{"developer", "tester", "architect"},
		},
		{
			name:     "Mention with excessive whitespace",
			input:    "   @developer    please    help   ",
			expected: []string{"developer"},
		},
		{
			name:     "Mention with CRLF line endings",
			input:    "@developer\r\n@tester\r\n@manager",
			expected: []string{"developer", "tester", "manager"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractMentionsImproved(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d mentions, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected mention[%d] = %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestValidateMentionText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid ASCII text",
			input:    "@developer please help",
			expected: true,
		},
		{
			name:     "Valid UTF-8 Vietnamese",
			input:    "Cảm ơn @developer",
			expected: true,
		},
		{
			name:     "Valid UTF-8 emoji",
			input:    "🎉 @developer great work!",
			expected: true,
		},
		{
			name:     "Too long text",
			input:    strings.Repeat("a", 100001),
			expected: false,
		},
		{
			name:     "Empty text",
			input:    "",
			expected: true,
		},
		{
			name:     "Valid multi-line",
			input:    "Line 1\nLine 2\n@developer",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMentionText(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExtractMentionsWithContext(t *testing.T) {
	input := "Hello @developer please fix the bug in @tester code"
	contexts := ExtractMentionsWithContext(input, 10)
	
	if len(contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(contexts))
		return
	}
	
	// Check first mention
	if contexts[0].Mention != "developer" {
		t.Errorf("Expected mention 'developer', got '%s'", contexts[0].Mention)
	}
	
	if !strings.Contains(contexts[0].BeforeText, "Hello") {
		t.Errorf("Expected 'Hello' in before text, got '%s'", contexts[0].BeforeText)
	}
	
	// Check second mention
	if contexts[1].Mention != "tester" {
		t.Errorf("Expected mention 'tester', got '%s'", contexts[1].Mention)
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CRLF to LF",
			input:    "Line1\r\nLine2\r\nLine3",
			expected: "Line1\nLine2\nLine3",
		},
		{
			name:     "CR to LF",
			input:    "Line1\rLine2\rLine3",
			expected: "Line1\nLine2\nLine3",
		},
		{
			name:     "Trim whitespace",
			input:    "  Line1  \n  Line2  ",
			expected: "Line1\nLine2",
		},
		{
			name:     "Mixed line endings",
			input:    "Line1\r\nLine2\nLine3\r",
			expected: "Line1\nLine2\nLine3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeText(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}

func BenchmarkExtractMentionsImproved(b *testing.B) {
	text := "Hello @developer and @tester, please review @architect's work. Contact @manager for approval."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractMentionsImproved(text)
	}
}

func BenchmarkExtractMentionsLongText(b *testing.B) {
	// Simulate long message with mentions
	text := strings.Repeat("Some text here. ", 100) + "@developer " + 
	        strings.Repeat("More text. ", 100) + "@tester " +
	        strings.Repeat("Even more. ", 100) + "@manager"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractMentionsImproved(text)
	}
}
