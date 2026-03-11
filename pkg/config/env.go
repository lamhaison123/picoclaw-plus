// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// LoadEnvFile loads environment variables from a .env file
// v0.2.1: Support for .env file loading
func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// .env file is optional
			return nil
		}
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logger.WarnCF("config", "Invalid .env line, skipping",
				map[string]any{
					"file": path,
					"line": lineNum,
					"text": line,
				})
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Reject empty keys
		if key == "" {
			logger.WarnCF("config", "Empty key in .env file, skipping",
				map[string]any{
					"file": path,
					"line": lineNum,
				})
			continue
		}

		// Remove quotes from value if present
		value = strings.Trim(value, `"'`)

		// Only set if not already set (existing env vars take precedence)
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				logger.WarnCF("config", "Failed to set environment variable",
					map[string]any{
						"key":   key,
						"error": err.Error(),
					})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	logger.InfoCF("config", "Loaded .env file",
		map[string]any{"path": path})

	return nil
}

// LoadEnvFiles loads multiple .env files in order
// Later files override earlier ones
func LoadEnvFiles(paths ...string) error {
	for _, path := range paths {
		if err := LoadEnvFile(path); err != nil {
			return err
		}
	}
	return nil
}
