// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"reflect"
	"testing"
)

func TestExtractMentions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single mention",
			input:    "@developer can you help?",
			expected: []string{"developer"},
		},
		{
			name:     "Multiple mentions",
			input:    "@architect @developer @tester review this",
			expected: []string{"architect", "developer", "tester"},
		},
		{
			name:     "Duplicate mentions",
			input:    "@developer @developer please help",
			expected: []string{"developer"},
		},
		{
			name:     "Email address should be ignored",
			input:    "Contact me at user@gmail.com",
			expected: []string{},
		},
		{
			name:     "Email with mention",
			input:    "@developer please email spawnacc2@gmail.com",
			expected: []string{"developer"},
		},
		{
			name:     "Multiple emails",
			input:    "Email user@yahoo.com or admin@hotmail.com",
			expected: []string{},
		},
		{
			name:     "Mention at start",
			input:    "@manager what's the status?",
			expected: []string{"manager"},
		},
		{
			name:     "Mention in middle",
			input:    "Hey @tester can you check?",
			expected: []string{"tester"},
		},
		{
			name:     "No mentions",
			input:    "This is a regular message",
			expected: []string{},
		},
		{
			name:     "Case insensitive",
			input:    "@Developer @TESTER @ArChItEcT",
			expected: []string{"developer", "tester", "architect"},
		},
		{
			name:     "Common email domains filtered",
			input:    "user@outlook.com and admin@icloud.com",
			expected: []string{},
		},
		{
			name:     "Mixed mentions and emails",
			input:    "@architect design this, then email me@company.com and @developer implement it",
			expected: []string{"architect", "developer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractMentions(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractMentions(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsEmailDomain(t *testing.T) {
	tests := []struct {
		domain   string
		expected bool
	}{
		{"gmail", true},
		{"yahoo", true},
		{"hotmail", true},
		{"outlook", true},
		{"developer", false},
		{"architect", false},
		{"tester", false},
		{"manager", false},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			result := isEmailDomain(tt.domain)
			if result != tt.expected {
				t.Errorf("isEmailDomain(%q) = %v, want %v", tt.domain, result, tt.expected)
			}
		})
	}
}
