// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
)

// legacySession represents the old JSON session format.
type legacySession struct {
	Key      string              `json:"key"`
	Messages []providers.Message `json:"messages"`
	Summary  string              `json:"summary,omitempty"`
	Created  time.Time           `json:"created"`
	Updated  time.Time           `json:"updated"`
}

// MigrateFromJSON migrates legacy JSON session files to JSONL format.
// v0.2.1: Support migration from legacy JSON sessions
//
// This function:
// 1. Scans the directory for .json files
// 2. Loads each JSON session
// 3. Writes it to JSONL format using SetHistory (idempotent)
// 4. Optionally backs up or removes the old JSON file
//
// Migration is idempotent: if a JSONL file already exists, it's skipped.
func MigrateFromJSON(ctx context.Context, jsonDir string, store Store, backup bool) error {
	files, err := os.ReadDir(jsonDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No directory to migrate
		}
		return fmt.Errorf("migration: read directory: %w", err)
	}

	migrated := 0
	skipped := 0
	failed := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		sessionPath := filepath.Join(jsonDir, file.Name())

		// Check if already migrated (JSONL file exists)
		jsonlPath := strings.TrimSuffix(sessionPath, ".json") + ".jsonl"
		if _, err := os.Stat(jsonlPath); err == nil {
			logger.DebugCF("migration", "Skipping already migrated session",
				map[string]any{"file": file.Name()})
			skipped++
			continue
		}

		// Load legacy JSON session
		data, err := os.ReadFile(sessionPath)
		if err != nil {
			logger.WarnCF("migration", "Failed to read legacy session",
				map[string]any{
					"file":  file.Name(),
					"error": err.Error(),
				})
			failed++
			continue
		}

		var legacy legacySession
		if err := json.Unmarshal(data, &legacy); err != nil {
			logger.WarnCF("migration", "Failed to parse legacy session",
				map[string]any{
					"file":  file.Name(),
					"error": err.Error(),
				})
			failed++
			continue
		}

		// Migrate to JSONL using SetHistory (idempotent)
		if err := store.SetHistory(ctx, legacy.Key, legacy.Messages); err != nil {
			logger.ErrorCF("migration", "Failed to migrate session",
				map[string]any{
					"key":   legacy.Key,
					"file":  file.Name(),
					"error": err.Error(),
				})
			failed++
			continue
		}

		// Migrate summary if present
		if legacy.Summary != "" {
			if err := store.SetSummary(ctx, legacy.Key, legacy.Summary); err != nil {
				logger.WarnCF("migration", "Failed to migrate summary",
					map[string]any{
						"key":   legacy.Key,
						"error": err.Error(),
					})
			}
		}

		logger.InfoCF("migration", "Migrated session",
			map[string]any{
				"key":           legacy.Key,
				"message_count": len(legacy.Messages),
			})
		migrated++

		// Backup or remove old JSON file
		if backup {
			backupPath := sessionPath + ".bak"
			if err := os.Rename(sessionPath, backupPath); err != nil {
				logger.WarnCF("migration", "Failed to backup old session",
					map[string]any{
						"file":  file.Name(),
						"error": err.Error(),
					})
			}
		} else {
			if err := os.Remove(sessionPath); err != nil {
				logger.WarnCF("migration", "Failed to remove old session",
					map[string]any{
						"file":  file.Name(),
						"error": err.Error(),
					})
			}
		}
	}

	logger.InfoCF("migration", "Migration complete",
		map[string]any{
			"migrated": migrated,
			"skipped":  skipped,
			"failed":   failed,
		})

	return nil
}

// AutoMigrate automatically migrates from JSON to JSONL if needed.
// This is called during initialization to transparently upgrade old sessions.
func AutoMigrate(ctx context.Context, sessionDir string, store Store) error {
	// Check if migration is needed (any .json files exist)
	files, err := os.ReadDir(sessionDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	hasJSON := false
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			hasJSON = true
			break
		}
	}

	if !hasJSON {
		return nil
	}

	logger.InfoCF("migration", "Detected legacy JSON sessions, starting migration",
		map[string]any{"dir": sessionDir})

	// Migrate with backup
	return MigrateFromJSON(ctx, sessionDir, store, true)
}
