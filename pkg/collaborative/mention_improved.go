// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ImprovedMentionParser handles UTF-8 and multi-line mentions correctly
var (
	// Match @mention with proper Unicode word boundaries
	// Supports: @developer, @architect, @开发者 (Unicode letters)
	// Ignores: user@gmail.com, @123numbers
	mentionRegexImproved = regexp.MustCompile(`(?:^|[\s\p{P}])@([\p{L}][\p{L}\p{N}_]*)`)
)

// ExtractMentionsImproved extracts @mentions with better UTF-8 handling
// Handles multi-line text, Unicode characters, and edge cases
func ExtractMentionsImproved(text string) []string {
	// Normalize text: handle different line endings and excessive whitespace
	text = normalizeText(text)

	matches := mentionRegexImproved.FindAllStringSubmatch(text, -1)
	mentions := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			mention := strings.ToLower(match[1])

			// Skip if empty or too short
			if len(mention) < 2 {
				continue
			}

			// Skip common email domains to avoid false positives
			if isEmailDomain(mention) {
				continue
			}

			// Skip if already seen (deduplication)
			if !seen[mention] {
				mentions = append(mentions, mention)
				seen[mention] = true
			}
		}
	}

	return mentions
}

// normalizeText normalizes text for consistent parsing
func normalizeText(text string) string {
	// Replace different line endings with standard \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Normalize excessive whitespace but preserve single spaces
	lines := strings.Split(text, "\n")
	trimmedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines at the end
		if trimmed != "" || len(trimmedLines) > 0 {
			trimmedLines = append(trimmedLines, trimmed)
		}
	}

	// Remove trailing empty lines
	for len(trimmedLines) > 0 && trimmedLines[len(trimmedLines)-1] == "" {
		trimmedLines = trimmedLines[:len(trimmedLines)-1]
	}

	return strings.Join(trimmedLines, "\n")
}

// ValidateMentionText checks if text is valid UTF-8 and safe to process
func ValidateMentionText(text string) bool {
	// Check if valid UTF-8
	if !isValidUTF8(text) {
		return false
	}

	// Check reasonable length (prevent DoS)
	if len(text) > 100000 { // 100KB limit
		return false
	}

	return true
}

// isValidUTF8 checks if string is valid UTF-8
func isValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

// ExtractMentionsWithContext extracts mentions with surrounding context
// Useful for debugging and logging
func ExtractMentionsWithContext(text string, contextChars int) []MentionContext {
	if contextChars <= 0 {
		contextChars = 20
	}

	matches := mentionRegexImproved.FindAllStringSubmatchIndex(text, -1)
	contexts := make([]MentionContext, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		// Extract mention text
		mentionStart := match[2]
		mentionEnd := match[3]
		mention := strings.ToLower(text[mentionStart:mentionEnd])

		// Skip duplicates and invalid mentions
		if seen[mention] || isEmailDomain(mention) || len(mention) < 2 {
			continue
		}
		seen[mention] = true

		// Extract context
		contextStart := match[0] - contextChars
		if contextStart < 0 {
			contextStart = 0
		}

		contextEnd := match[1] + contextChars
		if contextEnd > len(text) {
			contextEnd = len(text)
		}

		contexts = append(contexts, MentionContext{
			Mention:     mention,
			Position:    mentionStart,
			BeforeText:  text[contextStart:match[0]],
			MentionText: text[match[0]:match[1]],
			AfterText:   text[match[1]:contextEnd],
		})
	}

	return contexts
}

// MentionContext provides context around a mention for debugging
type MentionContext struct {
	Mention     string
	Position    int
	BeforeText  string
	MentionText string
	AfterText   string
}

// String formats the mention context for logging
func (mc MentionContext) String() string {
	return fmt.Sprintf("...%s[%s]%s... (pos: %d)",
		mc.BeforeText, mc.MentionText, mc.AfterText, mc.Position)
}

// ValidateMentionSafety checks for malicious patterns in mentions
func ValidateMentionSafety(text string) error {
	// Check for excessive repetition (e.g., @@@@@@@@)
	if strings.Count(text, "@") > 50 {
		return fmt.Errorf("excessive @ symbols detected (%d)", strings.Count(text, "@"))
	}

	// Check for extremely long mention names
	mentions := ExtractMentionsImproved(text)
	for _, mention := range mentions {
		if len(mention) > 100 {
			return fmt.Errorf("mention name too long: %s (max 100 chars)", mention)
		}

		// Check for control characters
		for _, r := range mention {
			if unicode.IsControl(r) {
				return fmt.Errorf("mention contains control character: %s", mention)
			}
		}
	}

	return nil
}

// SanitizeMentionText removes potentially dangerous content from mention text
func SanitizeMentionText(text string) string {
	// Remove control characters except newlines and tabs
	var sb strings.Builder
	for _, r := range text {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue // Skip control characters
		}
		sb.WriteRune(r)
	}

	return sb.String()
}
