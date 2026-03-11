// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package routing

import (
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
)

// Router selects models based on message complexity.
// v0.2.1: Model routing for cost optimization
type Router struct {
	config config.RoutingConfig
}

// NewRouter creates a new model router.
func NewRouter(cfg config.RoutingConfig) *Router {
	return &Router{
		config: cfg,
	}
}

// SelectModel chooses the best model for a message based on complexity.
// Returns the selected model name, or empty string if routing is disabled.
func (r *Router) SelectModel(msg string, history []providers.Message, primaryModel string) string {
	// If routing disabled, use primary model
	if !r.config.Enabled || len(r.config.Tiers) == 0 {
		return primaryModel
	}

	// Extract features and score complexity
	features := ExtractFeatures(msg, history)
	complexity := ScoreComplexity(features)

	// Log routing decision
	logger.DebugCF("routing", "Complexity analysis",
		map[string]any{
			"complexity":      complexity.String(),
			"token_estimate":  features.TokenEstimate,
			"code_blocks":     features.CodeBlockCount,
			"recent_tools":    features.RecentToolCalls,
			"conv_depth":      features.ConversationDepth,
			"has_attachments": features.HasAttachments,
		})

	// Select tier based on complexity
	var selectedTier *config.ModelRoutingTier
	switch complexity {
	case ComplexityLow:
		// Use cheapest tier
		if len(r.config.Tiers) > 0 {
			selectedTier = &r.config.Tiers[0]
		}
	case ComplexityMedium:
		// Use middle tier if available, otherwise cheapest
		if len(r.config.Tiers) > 1 {
			selectedTier = &r.config.Tiers[1]
		} else if len(r.config.Tiers) > 0 {
			selectedTier = &r.config.Tiers[0]
		}
	case ComplexityHigh:
		// Use most expensive tier
		if len(r.config.Tiers) > 0 {
			selectedTier = &r.config.Tiers[len(r.config.Tiers)-1]
		}
	}

	// If no tier selected or tier has no models, use primary
	if selectedTier == nil || len(selectedTier.Models) == 0 {
		logger.DebugCF("routing", "No tier available, using primary model",
			map[string]any{"model": primaryModel})
		return primaryModel
	}

	// Select first model from tier
	selectedModel := selectedTier.Models[0]

	logger.InfoCF("routing", "Model selected",
		map[string]any{
			"complexity": complexity.String(),
			"tier":       selectedTier.Name,
			"model":      selectedModel,
			"primary":    primaryModel,
		})

	return selectedModel
}
