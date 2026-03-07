// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"regexp"
	"strings"
)

var mentionRegex = regexp.MustCompile(`(?:^|[^\w.])@(\w+)`)

// ExtractMentions extracts @mentions from text
// Ignores email addresses (e.g., user@gmail.com)
func ExtractMentions(text string) []string {
	matches := mentionRegex.FindAllStringSubmatch(text, -1)
	mentions := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			mention := strings.ToLower(match[1])
			// Skip common email domains to avoid false positives
			if isEmailDomain(mention) {
				continue
			}
			if !seen[mention] {
				mentions = append(mentions, mention)
				seen[mention] = true
			}
		}
	}

	return mentions
}

// isEmailDomain checks if a mention is likely an email domain
func isEmailDomain(mention string) bool {
	emailDomains := map[string]bool{
		"gmail":    true,
		"yahoo":    true,
		"hotmail":  true,
		"outlook":  true,
		"icloud":   true,
		"proton":   true,
		"mail":     true,
		"email":    true,
		"live":     true,
		"msn":      true,
		"aol":      true,
		"yandex":   true,
		"zoho":     true,
		"gmx":      true,
		"fastmail": true,
	}
	return emailDomains[mention]
}
