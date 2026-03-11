// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/channels"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/constants"
	"github.com/sipeed/picoclaw/pkg/embedding"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/mcp"
	"github.com/sipeed/picoclaw/pkg/media"
	memory "github.com/sipeed/picoclaw/pkg/memory/vector"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/routing"
	"github.com/sipeed/picoclaw/pkg/skills"
	"github.com/sipeed/picoclaw/pkg/state"
	"github.com/sipeed/picoclaw/pkg/tools"
	"github.com/sipeed/picoclaw/pkg/utils"
)

type AgentLoop struct {
	bus              *bus.MessageBus
	cfg              *config.Config
	registry         *AgentRegistry
	state            *state.Manager
	running          atomic.Bool
	summarizing      sync.Map
	fallback         *providers.FallbackChain
	channelManager   *channels.Manager
	mediaStore       media.MediaStore
	reasoningSem     chan struct{}
	vectorStore      memory.VectorStore
	memoryProvider   memory.MemoryProvider
	embeddingService embedding.Service
	memoryEnabled    atomic.Bool       // BUG FIX #12: Use atomic.Bool to prevent race condition
	teamManager      tools.TeamManager // Team manager for multi-agent collaboration
	router           *routing.Router   // v0.2.1: Model router for cost optimization
}

// HandleMention implements MentionHandler interface
func (al *AgentLoop) HandleMention(ctx context.Context, mentionedID, message, channel, chatID string) (string, error) {
	// Process the mention by routing to the mentioned agent
	return al.ProcessWithAgent(ctx, mentionedID, message, fmt.Sprintf("mention:%s", mentionedID), channel, chatID)
}

// processOptions configures how a message is processed
type processOptions struct {
	SessionKey      string   // Session identifier for history/context
	Channel         string   // Target channel for tool execution
	ChatID          string   // Target chat ID for tool execution
	UserMessage     string   // User message content (may include prefix)
	Media           []string // media:// refs from inbound message
	DefaultResponse string   // Response when LLM returns empty
	EnableSummary   bool     // Whether to trigger summarization
	SendResponse    bool     // Whether to send response via bus
	NoHistory       bool     // If true, don't load session history (for heartbeat)
}

const defaultResponse = "I've completed processing but have no response to give. Increase `max_tool_iterations` in config.json."

func NewAgentLoop(
	cfg *config.Config,
	msgBus *bus.MessageBus,
	provider providers.LLMProvider,
) *AgentLoop {
	registry := NewAgentRegistry(cfg, provider)

	// Register shared tools to all agents (pass nil for teamManager, will be set later)
	registerSharedTools(cfg, msgBus, registry, provider, nil)

	// Set up shared fallback chain
	cooldown := providers.NewCooldownTracker()
	fallbackChain := providers.NewFallbackChain(cooldown)

	// Create state manager using default agent's workspace for channel recording
	defaultAgent := registry.GetDefaultAgent()
	var stateManager *state.Manager
	if defaultAgent != nil {
		stateManager = state.NewManager(defaultAgent.Workspace)
	}

	// Initialize vector memory if enabled
	var vectorStore memory.VectorStore
	var memoryProvider memory.MemoryProvider
	var embeddingService embedding.Service
	var memoryEnabled atomic.Bool // BUG FIX #12: Use atomic.Bool

	if cfg.Memory.Enabled {
		// Initialize embedding service
		embeddingService = initializeEmbeddingService(cfg.Memory.Embedding)

		// Initialize vector store
		if cfg.Memory.VectorStore.Provider != "none" && embeddingService != nil {
			vectorStore = initializeVectorStore(cfg.Memory.VectorStore, cfg.Memory.Embedding.Dimension)
			if vectorStore != nil {
				memoryEnabled.Store(true)
				logger.InfoCF("agent", "Vector memory initialized",
					map[string]any{
						"embedding_provider": cfg.Memory.Embedding.Provider,
						"vector_provider":    cfg.Memory.VectorStore.Provider,
						"dimension":          cfg.Memory.Embedding.Dimension,
					})
			}
		}

		// Initialize Memory Provider (e.g. Mem0)
		if cfg.Memory.MemoryProvider.Provider != "none" && cfg.Memory.MemoryProvider.Provider != "" {
			memoryProvider = initializeMemoryProvider(cfg.Memory.MemoryProvider)
			if memoryProvider != nil {
				memoryEnabled.Store(true)
				logger.InfoCF("agent", "Memory provider initialized",
					map[string]any{"provider": cfg.Memory.MemoryProvider.Provider})
			}
		}
	}

	al := &AgentLoop{
		bus:              msgBus,
		cfg:              cfg,
		registry:         registry,
		state:            stateManager,
		summarizing:      sync.Map{},
		reasoningSem:     make(chan struct{}, 10),
		fallback:         fallbackChain,
		vectorStore:      vectorStore,
		memoryProvider:   memoryProvider,
		embeddingService: embeddingService,
		teamManager:      nil,                            // Will be set via SetTeamManager
		router:           routing.NewRouter(cfg.Routing), // v0.2.1: Initialize model router
	}

	// Copy atomic.Bool state without copying the lock
	al.memoryEnabled.Store(memoryEnabled.Load())

	return al
}

// initializeEmbeddingService creates an embedding service based on configuration
func initializeEmbeddingService(cfg config.EmbeddingConfig) embedding.Service {
	switch cfg.Provider {
	case "openai", "litellm", "deepseek", "mistral", "vllm", "ollama":
		service, err := embedding.NewOpenAIService(embedding.Config{
			Provider:  cfg.Provider,
			Model:     cfg.Model,
			Dimension: cfg.Dimension,
			APIKey:    cfg.APIKey,
			BaseURL:   cfg.BaseURL,
			TimeoutMs: cfg.TimeoutMS,
		})
		if err != nil {
			logger.ErrorCF("agent", fmt.Sprintf("Failed to initialize %s embedding service", cfg.Provider),
				map[string]any{"error": err.Error()})
			return embedding.NewNullService(cfg.Dimension)
		}
		return service

	case "none", "":
		logger.InfoCF("agent", "No embedding provider specified, using null service",
			map[string]any{"dimension": cfg.Dimension})
		return embedding.NewNullService(cfg.Dimension)

	default:
		logger.WarnCF("agent", "Unknown embedding provider, using null service",
			map[string]any{"provider": cfg.Provider, "dimension": cfg.Dimension})
		return embedding.NewNullService(cfg.Dimension)
	}
}

// initializeVectorStore creates a vector store based on configuration
func initializeVectorStore(cfg config.VectorStoreConfig, dimension int) memory.VectorStore {
	switch cfg.Provider {
	case "qdrant":
		// Check if Qdrant is explicitly disabled
		if !cfg.Qdrant.Enabled {
			logger.InfoCF("agent", "Qdrant is disabled in config", nil)
			return nil
		}

		// Check if Qdrant is available (not built with no_qdrant tag)
		store, err := createQdrantStore(cfg, dimension)
		if err != nil {
			logger.ErrorCF("agent", "Failed to initialize Qdrant store",
				map[string]any{"error": err.Error()})
			return nil
		}
		return store

	case "lancedb":
		// Check if LanceDB is explicitly disabled
		if !cfg.LanceDB.Enabled {
			logger.InfoCF("agent", "LanceDB is disabled in config", nil)
			return nil
		}

		// Check if LanceDB is available (not built without CGO)
		store, err := createLanceDBStore(cfg, dimension)
		if err != nil {
			logger.ErrorCF("agent", "Failed to initialize LanceDB store",
				map[string]any{"error": err.Error()})
			return nil
		}
		return store

	case "none", "":
		logger.InfoCF("agent", "No vector store provider specified, vector memory disabled", nil)
		return nil

	default:
		logger.WarnCF("agent", "Unknown vector store provider",
			map[string]any{"provider": cfg.Provider})
		return nil
	}
}

// initializeMemoryProvider creates a memory provider based on configuration
func initializeMemoryProvider(cfg config.MemoryProviderConfig) memory.MemoryProvider {
	switch strings.ToLower(cfg.Provider) {
	case "mem0":
		if !cfg.Mem0.Enabled {
			logger.InfoCF("agent", "Mem0 is disabled in config", nil)
			return nil
		}

		breaker := memory.NewCircuitBreaker(memory.CircuitBreakerConfig{
			MaxFailures:   5,
			ResetTimeoutS: 30,
			HalfOpenMax:   3,
		})

		client, err := memory.NewMem0Client(memory.Mem0Config{
			Enabled:   cfg.Mem0.Enabled,
			URL:       cfg.Mem0.URL,
			APIKey:    cfg.Mem0.APIKey,
			TimeoutMS: cfg.Mem0.TimeoutMS,
		}, breaker)

		if err != nil {
			logger.ErrorCF("agent", "Failed to initialize Mem0 client",
				map[string]any{"error": err.Error()})
			return nil
		}
		return client

	case "none", "":
		logger.InfoCF("agent", "No memory provider specified", nil)
		return nil

	default:
		logger.WarnCF("agent", "Unknown memory provider",
			map[string]any{"provider": cfg.Provider})
		return nil
	}
}

// registerSharedTools registers tools that are shared across all agents (web, message, spawn, team).
func registerSharedTools(
	cfg *config.Config,
	msgBus *bus.MessageBus,
	registry *AgentRegistry,
	provider providers.LLMProvider,
	teamManager tools.TeamManager,
) {
	for _, agentID := range registry.ListAgentIDs() {
		agent, ok := registry.GetAgent(agentID)
		if !ok {
			continue
		}

		// Web tools (v0.2.1: check if enabled)
		if cfg.Tools.WebToolsEnabled {
			searchTool, err := tools.NewWebSearchTool(tools.WebSearchToolOptions{
				BraveAPIKey:          cfg.Tools.Web.Brave.APIKey,
				BraveMaxResults:      cfg.Tools.Web.Brave.MaxResults,
				BraveEnabled:         cfg.Tools.Web.Brave.Enabled,
				TavilyAPIKey:         cfg.Tools.Web.Tavily.APIKey,
				TavilyBaseURL:        cfg.Tools.Web.Tavily.BaseURL,
				TavilyMaxResults:     cfg.Tools.Web.Tavily.MaxResults,
				TavilyEnabled:        cfg.Tools.Web.Tavily.Enabled,
				DuckDuckGoMaxResults: cfg.Tools.Web.DuckDuckGo.MaxResults,
				DuckDuckGoEnabled:    cfg.Tools.Web.DuckDuckGo.Enabled,
				PerplexityAPIKey:     cfg.Tools.Web.Perplexity.APIKey,
				PerplexityMaxResults: cfg.Tools.Web.Perplexity.MaxResults,
				PerplexityEnabled:    cfg.Tools.Web.Perplexity.Enabled,
				SearXNGBaseURL:       cfg.Tools.Web.SearXNG.BaseURL,
				SearXNGMaxResults:    cfg.Tools.Web.SearXNG.MaxResults,
				SearXNGEnabled:       cfg.Tools.Web.SearXNG.Enabled,
				GLMAPIKey:            cfg.Tools.Web.GLM.APIKey,
				GLMMaxResults:        cfg.Tools.Web.GLM.MaxResults,
				GLMEnabled:           cfg.Tools.Web.GLM.Enabled,
				ExaAPIKey:            cfg.Tools.Web.Exa.APIKey,
				ExaMaxResults:        cfg.Tools.Web.Exa.MaxResults,
				ExaEnabled:           cfg.Tools.Web.Exa.Enabled,
				Proxy:                cfg.Tools.Web.Proxy,
			})
			if err != nil {
				logger.ErrorCF("agent", "Failed to create web search tool", map[string]any{"error": err.Error()})
			} else if searchTool != nil {
				agent.Tools.Register(searchTool)
			}
			fetchTool, err := tools.NewWebFetchToolWithProxy(50000, cfg.Tools.Web.Proxy, cfg.Tools.Web.FetchLimitBytes)
			if err != nil {
				logger.ErrorCF("agent", "Failed to create web fetch tool", map[string]any{"error": err.Error()})
			} else {
				agent.Tools.Register(fetchTool)
			}
		}

		// Hardware tools (v0.2.1: check if enabled) - Linux only, returns error on other platforms
		if cfg.Tools.HardwareToolsEnabled {
			agent.Tools.Register(tools.NewI2CTool())
			agent.Tools.Register(tools.NewSPITool())
		}

		// Message tool (v0.2.1: check if enabled)
		if cfg.Tools.MessageToolEnabled {
			messageTool := tools.NewMessageTool()
			messageTool.SetSendCallback(func(channel, chatID, content string) error {
				pubCtx, pubCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer pubCancel()
				return msgBus.PublishOutbound(pubCtx, bus.OutboundMessage{
					Channel: channel,
					ChatID:  chatID,
					Content: content,
				})
			})
			agent.Tools.Register(messageTool)
		}

		// Skill discovery and installation tools (v0.2.1: check if enabled)
		if cfg.Tools.SkillToolsEnabled {
			registryMgr := skills.NewRegistryManagerFromConfig(skills.RegistryConfig{
				MaxConcurrentSearches: cfg.Tools.Skills.MaxConcurrentSearches,
				ClawHub:               skills.ClawHubConfig(cfg.Tools.Skills.Registries.ClawHub),
			})
			searchCache := skills.NewSearchCache(
				cfg.Tools.Skills.SearchCache.MaxSize,
				time.Duration(cfg.Tools.Skills.SearchCache.TTLSeconds)*time.Second,
			)
			agent.Tools.Register(tools.NewFindSkillsTool(registryMgr, searchCache))
			agent.Tools.Register(tools.NewInstallSkillTool(registryMgr, agent.Workspace))
		}

		// Spawn tool (v0.2.1: check if enabled) with allowlist checker
		if cfg.Tools.SpawnToolEnabled {
			subagentManager := tools.NewSubagentManager(provider, agent.Model, agent.Workspace, msgBus)
			subagentManager.SetLLMOptions(agent.MaxTokens, agent.Temperature)
			spawnTool := tools.NewSpawnTool(subagentManager)
			currentAgentID := agentID
			spawnTool.SetAllowlistChecker(func(targetAgentID string) bool {
				return registry.CanSpawnSubagent(currentAgentID, targetAgentID)
			})
			agent.Tools.Register(spawnTool)
		}

		// Team delegation tools (v0.2.1: check if enabled) (if team manager available)
		if cfg.Tools.TeamToolsEnabled && teamManager != nil {
			agent.Tools.Register(tools.NewTeamDelegationTool(teamManager))
			agent.Tools.Register(tools.NewTeamStatusTool(teamManager))
			logger.DebugCF("agent", "Registered team tools",
				map[string]any{"agent_id": agentID})
		}
	}
}

func (al *AgentLoop) Run(ctx context.Context) error {
	al.running.Store(true)

	// Initialize MCP servers for all agents
	if al.cfg.Tools.MCP.Enabled {
		mcpManager := mcp.NewManager()
		// Ensure MCP connections are cleaned up on exit, regardless of initialization success
		// This fixes resource leak when LoadFromMCPConfig partially succeeds then fails
		// The defer ensures cleanup happens even if:
		// 1. LoadFromMCPConfig fails after partial initialization
		// 2. Tool registration fails
		// 3. Main loop exits due to context cancellation
		defer func() {
			if err := mcpManager.Close(); err != nil {
				logger.ErrorCF("agent", "Failed to close MCP manager",
					map[string]any{
						"error": err.Error(),
					})
			}
		}()

		defaultAgent := al.registry.GetDefaultAgent()
		var workspacePath string
		if defaultAgent != nil && defaultAgent.Workspace != "" {
			workspacePath = defaultAgent.Workspace
		} else {
			workspacePath = al.cfg.WorkspacePath()
		}

		if err := mcpManager.LoadFromMCPConfig(ctx, al.cfg.Tools.MCP, workspacePath); err != nil {
			logger.WarnCF("agent", "Failed to load MCP servers, MCP tools will not be available",
				map[string]any{
					"error": err.Error(),
				})
		} else {
			// Register MCP tools for all agents
			servers := mcpManager.GetServers()
			uniqueTools := 0
			totalRegistrations := 0
			agentIDs := al.registry.ListAgentIDs()
			agentCount := len(agentIDs)

			for serverName, conn := range servers {
				uniqueTools += len(conn.Tools)
				for _, tool := range conn.Tools {
					for _, agentID := range agentIDs {
						agent, ok := al.registry.GetAgent(agentID)
						if !ok {
							continue
						}
						mcpTool := tools.NewMCPTool(mcpManager, serverName, tool)
						agent.Tools.Register(mcpTool)
						totalRegistrations++
						logger.DebugCF("agent", "Registered MCP tool",
							map[string]any{
								"agent_id": agentID,
								"server":   serverName,
								"tool":     tool.Name,
								"name":     mcpTool.Name(),
							})
					}
				}
			}
			logger.InfoCF("agent", "MCP tools registered successfully",
				map[string]any{
					"server_count":        len(servers),
					"unique_tools":        uniqueTools,
					"total_registrations": totalRegistrations,
					"agent_count":         agentCount,
				})
		}
	}

	for al.running.Load() {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, ok := al.bus.ConsumeInbound(ctx)
			if !ok {
				continue
			}

			// Process message
			go func(msg bus.InboundMessage) {
				// TODO: Re-enable media cleanup after inbound media is properly consumed by the agent.
				// Currently disabled because files are deleted before the LLM can access their content.
				// defer func() {
				// 	if al.mediaStore != nil && msg.MediaScope != "" {
				// 		if releaseErr := al.mediaStore.ReleaseAll(msg.MediaScope); releaseErr != nil {
				// 			logger.WarnCF("agent", "Failed to release media", map[string]any{
				// 				"scope": msg.MediaScope,
				// 				"error": releaseErr.Error(),
				// 			})
				// 		}
				// 	}
				// }()

				response, err := al.processMessage(ctx, msg)
				if err != nil {
					// Wrap error properly to preserve context
					response = fmt.Sprintf("Error processing message: %v", err)
					logger.ErrorCF("agent", "Message processing failed",
						map[string]any{
							"channel": msg.Channel,
							"chat_id": msg.ChatID,
							"error":   err.Error(),
						})
				}

				if response != "" {
					// Check if the message tool already sent a response during this round.
					// If so, skip publishing to avoid duplicate messages to the user.
					// Use default agent's tools to check (message tool is shared).
					alreadySent := false
					defaultAgent := al.registry.GetDefaultAgent()
					if defaultAgent != nil {
						if tool, ok := defaultAgent.Tools.Get("message"); ok {
							if mt, ok := tool.(*tools.MessageTool); ok {
								alreadySent = mt.HasSentInRound()
							}
						}
					}

					if !alreadySent {
						al.bus.PublishOutbound(ctx, bus.OutboundMessage{
							Channel: msg.Channel,
							ChatID:  msg.ChatID,
							Content: response,
						})
						logger.InfoCF("agent", "Published outbound response",
							map[string]any{
								"channel":     msg.Channel,
								"chat_id":     msg.ChatID,
								"content_len": len(response),
							})
					} else {
						logger.DebugCF(
							"agent",
							"Skipped outbound (message tool already sent)",
							map[string]any{"channel": msg.Channel},
						)
					}
				}
			}(msg)
		}
	}

	return nil
}

func (al *AgentLoop) Stop() {
	al.running.Store(false)

	// Close vector memory connections
	if al.vectorStore != nil {
		if err := al.vectorStore.Close(); err != nil {
			logger.ErrorCF("agent", "Failed to close vector store",
				map[string]any{"error": err.Error()})
		}
	}

	if al.embeddingService != nil {
		if err := al.embeddingService.Close(); err != nil {
			logger.ErrorCF("agent", "Failed to close embedding service",
				map[string]any{"error": err.Error()})
		}
	}
}

func (al *AgentLoop) RegisterTool(tool tools.Tool) {
	for _, agentID := range al.registry.ListAgentIDs() {
		if agent, ok := al.registry.GetAgent(agentID); ok {
			agent.Tools.Register(tool)
		}
	}
}

func (al *AgentLoop) SetChannelManager(cm *channels.Manager) {
	al.channelManager = cm
}

// SetMediaStore injects a MediaStore for media lifecycle management.
func (al *AgentLoop) SetMediaStore(s media.MediaStore) {
	al.mediaStore = s
}

// SetTeamManager sets the team manager and registers team tools
func (al *AgentLoop) SetTeamManager(teamManager tools.TeamManager) {
	al.teamManager = teamManager

	// Re-register shared tools with team manager
	registerSharedTools(al.cfg, al.bus, al.registry, al.registry.GetDefaultAgent().Provider, teamManager)

	logger.InfoCF("agent", "Team manager configured, team tools registered", nil)
}

// GetRegistry returns the agent registry
func (al *AgentLoop) GetRegistry() *AgentRegistry {
	return al.registry
}

// inferMediaType determines the media type ("image", "audio", "video", "file")
// from a filename and MIME content type.
func inferMediaType(filename, contentType string) string {
	ct := strings.ToLower(contentType)
	fn := strings.ToLower(filename)

	if strings.HasPrefix(ct, "image/") {
		return "image"
	}
	if strings.HasPrefix(ct, "audio/") || ct == "application/ogg" {
		return "audio"
	}
	if strings.HasPrefix(ct, "video/") {
		return "video"
	}

	// Fallback: infer from extension
	ext := filepath.Ext(fn)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
		return "image"
	case ".mp3", ".wav", ".ogg", ".m4a", ".flac", ".aac", ".wma", ".opus":
		return "audio"
	case ".mp4", ".avi", ".mov", ".webm", ".mkv":
		return "video"
	}

	return "file"
}

// RecordLastChannel records the last active channel for this workspace.
// This uses the atomic state save mechanism to prevent data loss on crash.
func (al *AgentLoop) RecordLastChannel(channel string) error {
	if al.state == nil {
		return nil
	}
	return al.state.SetLastChannel(channel)
}

// RecordLastChatID records the last active chat ID for this workspace.
// This uses the atomic state save mechanism to prevent data loss on crash.
func (al *AgentLoop) RecordLastChatID(chatID string) error {
	if al.state == nil {
		return nil
	}
	return al.state.SetLastChatID(chatID)
}

func (al *AgentLoop) ProcessDirect(
	ctx context.Context,
	content, sessionKey string,
) (string, error) {
	return al.ProcessDirectWithChannel(ctx, content, sessionKey, "cli", "direct")
}

func (al *AgentLoop) ProcessDirectWithChannel(
	ctx context.Context,
	content, sessionKey, channel, chatID string,
) (string, error) {
	msg := bus.InboundMessage{
		Channel:    channel,
		SenderID:   "cron",
		ChatID:     chatID,
		Content:    content,
		SessionKey: sessionKey,
	}

	return al.processMessage(ctx, msg)
}

// ProcessWithAgent processes a message with a specific agent ID
func (al *AgentLoop) ProcessWithAgent(
	ctx context.Context,
	agentID, content, sessionKey, channel, chatID string,
) (string, error) {
	// Get the specific agent
	agent, ok := al.registry.GetAgent(agentID)
	if !ok {
		return "", fmt.Errorf("agent not found: %s", agentID)
	}

	// Build session key with agent prefix if not already set
	if sessionKey == "" || !strings.HasPrefix(sessionKey, "agent:") {
		sessionKey = fmt.Sprintf("agent:%s:%s", agentID, sessionKey)
	}

	logger.InfoCF("agent", "Processing with specific agent",
		map[string]any{
			"agent_id":    agentID,
			"session_key": sessionKey,
			"channel":     channel,
			"chat_id":     chatID,
		})

	// Run agent loop directly with the specified agent
	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      sessionKey,
		Channel:         channel,
		ChatID:          chatID,
		UserMessage:     content,
		DefaultResponse: defaultResponse,
		EnableSummary:   true,
		SendResponse:    false,
	})
}

// ProcessHeartbeat processes a heartbeat request without session history.
// Each heartbeat is independent and doesn't accumulate context.
func (al *AgentLoop) ProcessHeartbeat(
	ctx context.Context,
	content, channel, chatID string,
) (string, error) {
	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return "", fmt.Errorf("no default agent for heartbeat")
	}
	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      "heartbeat",
		Channel:         channel,
		ChatID:          chatID,
		UserMessage:     content,
		DefaultResponse: defaultResponse,
		EnableSummary:   false,
		SendResponse:    false,
		NoHistory:       true, // Don't load session history for heartbeat
	})
}

func (al *AgentLoop) processMessage(ctx context.Context, msg bus.InboundMessage) (string, error) {
	// Add message preview to log (show full content for error messages)
	var logContent string
	if strings.Contains(msg.Content, "Error:") || strings.Contains(msg.Content, "error") {
		logContent = msg.Content // Full content for errors
	} else {
		logContent = utils.Truncate(msg.Content, 80)
	}
	logger.InfoCF(
		"agent",
		fmt.Sprintf("Processing message from %s:%s: %s", msg.Channel, msg.SenderID, logContent),
		map[string]any{
			"channel":     msg.Channel,
			"chat_id":     msg.ChatID,
			"sender_id":   msg.SenderID,
			"session_key": msg.SessionKey,
		},
	)

	// Route system messages to processSystemMessage
	if msg.Channel == "system" {
		return al.processSystemMessage(ctx, msg)
	}

	// Check for commands
	if response, handled := al.handleCommand(ctx, msg); handled {
		return response, nil
	}

	// Route to determine agent and session key
	route := al.registry.ResolveRoute(routing.RouteInput{
		Channel:    msg.Channel,
		AccountID:  msg.Metadata["account_id"],
		Peer:       extractPeer(msg),
		ParentPeer: extractParentPeer(msg),
		GuildID:    msg.Metadata["guild_id"],
		TeamID:     msg.Metadata["team_id"],
	})

	agent, ok := al.registry.GetAgent(route.AgentID)
	if !ok {
		agent = al.registry.GetDefaultAgent()
	}
	if agent == nil {
		return "", fmt.Errorf("no agent available for route (agent_id=%s)", route.AgentID)
	}

	// Reset message-tool state for this round so we don't skip publishing due to a previous round.
	if tool, ok := agent.Tools.Get("message"); ok {
		if mt, ok := tool.(tools.ContextualTool); ok {
			mt.SetContext(msg.Channel, msg.ChatID)
		}
	}

	// Use routed session key, but honor pre-set agent-scoped keys (for ProcessDirect/cron)
	sessionKey := route.SessionKey
	if msg.SessionKey != "" && strings.HasPrefix(msg.SessionKey, "agent:") {
		sessionKey = msg.SessionKey
	}

	logger.InfoCF("agent", "Routed message",
		map[string]any{
			"agent_id":    agent.ID,
			"session_key": sessionKey,
			"matched_by":  route.MatchedBy,
		})

	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      sessionKey,
		Channel:         msg.Channel,
		ChatID:          msg.ChatID,
		UserMessage:     msg.Content,
		Media:           msg.Media,
		DefaultResponse: defaultResponse,
		EnableSummary:   true,
		SendResponse:    false,
	})
}

func (al *AgentLoop) processSystemMessage(
	ctx context.Context,
	msg bus.InboundMessage,
) (string, error) {
	if msg.Channel != "system" {
		return "", fmt.Errorf(
			"processSystemMessage called with non-system message channel: %s",
			msg.Channel,
		)
	}

	logger.InfoCF("agent", "Processing system message",
		map[string]any{
			"sender_id": msg.SenderID,
			"chat_id":   msg.ChatID,
		})

	// Parse origin channel from chat_id (format: "channel:chat_id")
	var originChannel, originChatID string
	if idx := strings.Index(msg.ChatID, ":"); idx > 0 {
		originChannel = msg.ChatID[:idx]
		originChatID = msg.ChatID[idx+1:]
	} else {
		originChannel = "cli"
		originChatID = msg.ChatID
	}

	// Extract subagent result from message content
	// Format: "Task 'label' completed.\n\nResult:\n<actual content>"
	content := msg.Content
	if idx := strings.Index(content, "Result:\n"); idx >= 0 {
		content = content[idx+8:] // Extract just the result part
	}

	// Skip internal channels - only log, don't send to user
	if constants.IsInternalChannel(originChannel) {
		logger.InfoCF("agent", "Subagent completed (internal channel)",
			map[string]any{
				"sender_id":   msg.SenderID,
				"content_len": len(content),
				"channel":     originChannel,
			})
		return "", nil
	}

	// Use default agent for system messages
	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return "", fmt.Errorf("no default agent for system message")
	}

	// Use the origin session for context
	sessionKey := routing.BuildAgentMainSessionKey(agent.ID)

	return al.runAgentLoop(ctx, agent, processOptions{
		SessionKey:      sessionKey,
		Channel:         originChannel,
		ChatID:          originChatID,
		UserMessage:     fmt.Sprintf("[System: %s] %s", msg.SenderID, msg.Content),
		DefaultResponse: "Background task completed.",
		EnableSummary:   false,
		SendResponse:    true,
	})
}

// runAgentLoop is the core message processing logic.
func (al *AgentLoop) runAgentLoop(
	ctx context.Context,
	agent *AgentInstance,
	opts processOptions,
) (string, error) {
	// 0. Record last channel for heartbeat notifications (skip internal channels)
	if opts.Channel != "" && opts.ChatID != "" {
		// Don't record internal channels (cli, system, subagent)
		if !constants.IsInternalChannel(opts.Channel) {
			channelKey := fmt.Sprintf("%s:%s", opts.Channel, opts.ChatID)
			if err := al.RecordLastChannel(channelKey); err != nil {
				logger.WarnCF(
					"agent",
					"Failed to record last channel",
					map[string]any{"error": err.Error()},
				)
			}
		}
	}

	// 1. Update tool contexts
	al.updateToolContexts(agent, opts.Channel, opts.ChatID)

	// 2. Build messages (skip history for heartbeat)
	var history []providers.Message
	var summary string
	if !opts.NoHistory {
		history = agent.Sessions.GetHistory(opts.SessionKey)
		summary = agent.Sessions.GetSummary(opts.SessionKey)
	}

	// Search vector memory for relevant context
	var memoryContext string
	if al.memoryEnabled.Load() && opts.UserMessage != "" {
		if al.vectorStore != nil {
			memoryContext = al.searchVectorMemory(ctx, opts.UserMessage, opts.SessionKey)
		}

		var extendedContext string
		if al.memoryProvider != nil {
			extendedContext = al.searchMemoryProvider(ctx, opts.UserMessage, opts.SessionKey)
		}

		if extendedContext != "" {
			if memoryContext != "" {
				memoryContext += "\n\nAdditional Personalized Context:\n" + extendedContext
			} else {
				memoryContext = "Personalized Context:\n" + extendedContext
			}
		}
	}

	messages := agent.ContextBuilder.BuildMessages(
		history,
		summary,
		opts.UserMessage,
		opts.Media,
		opts.Channel,
		opts.ChatID,
	)

	// Inject memory context into system prompt if found
	if memoryContext != "" {
		al.injectMemoryContext(messages, memoryContext)
	}

	// Resolve media:// refs to base64 data URLs (streaming)
	maxMediaSize := al.cfg.Agents.Defaults.GetMaxMediaSize()
	messages = resolveMediaRefs(messages, al.mediaStore, maxMediaSize)

	// 3. Save user message to session
	agent.Sessions.AddMessage(opts.SessionKey, "user", opts.UserMessage)

	// 4. Run LLM iteration loop
	finalContent, iteration, err := al.runLLMIteration(ctx, agent, messages, opts)
	if err != nil {
		return "", err
	}

	// Check if context was cancelled during LLM iteration
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	// If last tool had ForUser content and we already sent it, we might not need to send final response
	// This is controlled by the tool's Silent flag and ForUser content

	// 5. Handle empty response
	if finalContent == "" {
		finalContent = opts.DefaultResponse
	} else {
		// 6. Save final assistant message to session
		agent.Sessions.AddMessage(opts.SessionKey, "assistant", finalContent)
		agent.Sessions.Save(opts.SessionKey)

		// Store in vector memory (async)
		if al.memoryEnabled.Load() && !opts.NoHistory {
			// BUG FIX: Pass parent context instead of Background to respect cancellation
			// Create a detached context with timeout to prevent blocking shutdown
			storeCtx, storeCancel := context.WithTimeout(context.Background(), 45*time.Second)
			go func(a *AgentInstance) {
				defer storeCancel()
				al.storeInVectorMemory(storeCtx, a, opts.SessionKey,
					opts.UserMessage, finalContent, opts.Channel)
			}(agent)
		}
	}

	// 7. Optional: summarization (skip if context cancelled)
	if opts.EnableSummary && ctx.Err() == nil {
		al.maybeSummarize(agent, opts.SessionKey, opts.Channel, opts.ChatID)
	}

	// 8. Optional: send response via bus (skip if context cancelled)
	if opts.SendResponse && ctx.Err() == nil {
		al.bus.PublishOutbound(ctx, bus.OutboundMessage{
			Channel: opts.Channel,
			ChatID:  opts.ChatID,
			Content: finalContent,
		})
	}

	// 9. Log response
	responsePreview := utils.Truncate(finalContent, 120)
	logger.InfoCF("agent", fmt.Sprintf("Response: %s", responsePreview),
		map[string]any{
			"agent_id":     agent.ID,
			"session_key":  opts.SessionKey,
			"iterations":   iteration,
			"final_length": len(finalContent),
		})

	return finalContent, nil
}

func (al *AgentLoop) targetReasoningChannelID(channelName string) (chatID string) {
	if al.channelManager == nil {
		return ""
	}
	if ch, ok := al.channelManager.GetChannel(channelName); ok {
		return ch.ReasoningChannelID()
	}
	return ""
}

func (al *AgentLoop) handleReasoning(
	ctx context.Context,
	reasoningContent, channelName, channelID string,
) {
	if reasoningContent == "" || channelName == "" || channelID == "" {
		return
	}

	// Check parent context cancellation before attempting to publish
	if ctx.Err() != nil {
		return
	}

	// Try to acquire semaphore, skip if full (best-effort)
	select {
	case al.reasoningSem <- struct{}{}:
		// Run in goroutine with tracking to prevent accumulation
		go func() {
			// Create timeout context for entire goroutine
			goroutineCtx, goroutineCancel := context.WithTimeout(ctx, 10*time.Second)
			defer goroutineCancel()

			defer func() {
				<-al.reasoningSem // Release semaphore
				if r := recover(); r != nil {
					logger.ErrorCF("agent", "Panic in reasoning handler",
						map[string]any{
							"channel": channelName,
							"panic":   r,
						})
				}
			}()

			// Use shorter timeout for publish operation
			pubCtx, pubCancel := context.WithTimeout(goroutineCtx, 5*time.Second)
			defer pubCancel()

			// Create done channel to track completion
			done := make(chan error, 1)
			go func() {
				done <- al.bus.PublishOutbound(pubCtx, bus.OutboundMessage{
					Channel: channelName,
					ChatID:  channelID,
					Content: reasoningContent,
				})
			}()

			// Wait for completion or timeout
			select {
			case err := <-done:
				if err != nil {
					// Treat context.DeadlineExceeded / context.Canceled as expected
					// (bus full under load, or parent canceled).  Check the error
					// itself rather than ctx.Err(), because pubCtx may time out
					// (5 s) while the parent ctx is still active.
					// Also treat ErrBusClosed as expected — it occurs during normal
					// shutdown when the bus is closed before all goroutines finish.
					if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) ||
						errors.Is(err, bus.ErrBusClosed) {
						logger.DebugCF("agent", "Reasoning publish skipped (timeout/cancel)", map[string]any{
							"channel": channelName,
							"error":   err.Error(),
						})
					} else {
						logger.WarnCF("agent", "Failed to publish reasoning (best-effort)", map[string]any{
							"channel": channelName,
							"error":   err.Error(),
						})
					}
				}
			case <-goroutineCtx.Done():
				logger.WarnCF("agent", "Reasoning goroutine timeout", map[string]any{
					"channel": channelName,
					"error":   goroutineCtx.Err().Error(),
				})
			}
		}()
	default:
		// Semaphore full, skip (best-effort reasoning)
		logger.DebugCF("agent", "Reasoning skipped (too many concurrent)", map[string]any{
			"channel": channelName,
		})
	}
}

// runLLMIteration executes the LLM call loop with tool handling.
func (al *AgentLoop) runLLMIteration(
	ctx context.Context,
	agent *AgentInstance,
	messages []providers.Message,
	opts processOptions,
) (string, int, error) {
	iteration := 0
	var finalContent string

	for iteration < agent.MaxIterations {
		iteration++

		logger.DebugCF("agent", "LLM iteration",
			map[string]any{
				"agent_id":  agent.ID,
				"iteration": iteration,
				"max":       agent.MaxIterations,
			})

		// v0.2.1: Model routing - select model based on complexity
		selectedModel := agent.Model
		if al.router != nil && iteration == 1 {
			// Only route on first iteration (user message)
			// Get last user message for routing decision
			var lastUserMsg string
			for i := len(messages) - 1; i >= 0; i-- {
				if messages[i].Role == "user" {
					lastUserMsg = messages[i].Content
					break
				}
			}

			if lastUserMsg != "" {
				routedModel := al.router.SelectModel(lastUserMsg, messages, agent.Model)
				if routedModel != "" && routedModel != agent.Model {
					selectedModel = routedModel
					logger.InfoCF("agent", "Model routed",
						map[string]any{
							"from": agent.Model,
							"to":   selectedModel,
						})
				}
			}
		}

		// Build tool definitions
		providerToolDefs := agent.Tools.ToProviderDefs()

		// Log LLM request details
		logger.DebugCF("agent", "LLM request",
			map[string]any{
				"agent_id":          agent.ID,
				"iteration":         iteration,
				"model":             selectedModel,
				"messages_count":    len(messages),
				"tools_count":       len(providerToolDefs),
				"max_tokens":        agent.MaxTokens,
				"temperature":       agent.Temperature,
				"system_prompt_len": len(messages[0].Content),
			})

		// Log full messages (detailed)
		logger.DebugCF("agent", "Full LLM request",
			map[string]any{
				"iteration":     iteration,
				"messages_json": formatMessagesForLog(messages),
				"tools_json":    formatToolsForLog(providerToolDefs),
			})

		// Call LLM with fallback chain if candidates are configured.
		var response *providers.LLMResponse
		var err error

		callLLM := func() (*providers.LLMResponse, error) {
			if len(agent.Candidates) > 1 && al.fallback != nil {
				fbResult, fbErr := al.fallback.Execute(
					ctx,
					agent.Candidates,
					func(ctx context.Context, provider, model string) (*providers.LLMResponse, error) {
						return agent.Provider.Chat(
							ctx,
							messages,
							providerToolDefs,
							model,
							map[string]any{
								"max_tokens":       agent.MaxTokens,
								"temperature":      agent.Temperature,
								"prompt_cache_key": agent.ID,
							},
						)
					},
				)
				if fbErr != nil {
					return nil, fbErr
				}
				if fbResult.Provider != "" && len(fbResult.Attempts) > 0 {
					logger.InfoCF(
						"agent",
						fmt.Sprintf("Fallback: succeeded with %s/%s after %d attempts",
							fbResult.Provider, fbResult.Model, len(fbResult.Attempts)+1),
						map[string]any{"agent_id": agent.ID, "iteration": iteration},
					)
				}
				return fbResult.Response, nil
			}
			return agent.Provider.Chat(ctx, messages, providerToolDefs, selectedModel, map[string]any{
				"max_tokens":       agent.MaxTokens,
				"temperature":      agent.Temperature,
				"prompt_cache_key": agent.ID,
			})
		}

		// Retry loop for context/token errors
		maxRetries := 2
		for retry := 0; retry <= maxRetries; retry++ {
			response, err = callLLM()
			if err == nil {
				break
			}

			errMsg := strings.ToLower(err.Error())

			// Check if this is a network/HTTP timeout — not a context window error.
			isTimeoutError := errors.Is(err, context.DeadlineExceeded) ||
				strings.Contains(errMsg, "deadline exceeded") ||
				strings.Contains(errMsg, "client.timeout") ||
				strings.Contains(errMsg, "timed out") ||
				strings.Contains(errMsg, "timeout exceeded")

			// Detect real context window / token limit errors, excluding network timeouts.
			isContextError := !isTimeoutError && (strings.Contains(errMsg, "context_length_exceeded") ||
				strings.Contains(errMsg, "context window") ||
				strings.Contains(errMsg, "maximum context length") ||
				strings.Contains(errMsg, "token limit") ||
				strings.Contains(errMsg, "too many tokens") ||
				strings.Contains(errMsg, "max_tokens") ||
				strings.Contains(errMsg, "invalidparameter") ||
				strings.Contains(errMsg, "prompt is too long") ||
				strings.Contains(errMsg, "request too large"))

			if isTimeoutError && retry < maxRetries {
				backoff := time.Duration(retry+1) * 5 * time.Second
				logger.WarnCF("agent", "Timeout error, retrying after backoff", map[string]any{
					"error":   err.Error(),
					"retry":   retry,
					"backoff": backoff.String(),
				})

				// Respect context cancellation during backoff
				select {
				case <-ctx.Done():
					return "", iteration, ctx.Err()
				case <-time.After(backoff):
					// Continue to retry
				}
				continue
			}

			if isContextError && retry < maxRetries {
				logger.WarnCF(
					"agent",
					"Context window error detected, attempting compression",
					map[string]any{
						"error": err.Error(),
						"retry": retry,
					},
				)

				if retry == 0 && !constants.IsInternalChannel(opts.Channel) {
					al.bus.PublishOutbound(ctx, bus.OutboundMessage{
						Channel: opts.Channel,
						ChatID:  opts.ChatID,
						Content: "Context window exceeded. Compressing history and retrying...",
					})
				}

				al.forceCompression(agent, opts.SessionKey)
				newHistory := agent.Sessions.GetHistory(opts.SessionKey)
				newSummary := agent.Sessions.GetSummary(opts.SessionKey)
				messages = agent.ContextBuilder.BuildMessages(
					newHistory, newSummary, "",
					nil, opts.Channel, opts.ChatID,
				)
				continue
			}
			break
		}

		if err != nil {
			logger.ErrorCF("agent", "LLM call failed",
				map[string]any{
					"agent_id":  agent.ID,
					"iteration": iteration,
					"error":     err.Error(),
				})
			return "", iteration, fmt.Errorf("LLM call failed after retries: %w", err)
		}

		go al.handleReasoning(
			ctx,
			response.Reasoning,
			opts.Channel,
			al.targetReasoningChannelID(opts.Channel),
		)

		// v0.2.1: Handle extended thinking (Anthropic reasoning_content)
		if response.ReasoningContent != "" {
			go al.handleReasoning(
				ctx,
				response.ReasoningContent,
				opts.Channel,
				al.targetReasoningChannelID(opts.Channel),
			)
		}

		logger.DebugCF("agent", "LLM response",
			map[string]any{
				"agent_id":          agent.ID,
				"iteration":         iteration,
				"content_chars":     len(response.Content),
				"tool_calls":        len(response.ToolCalls),
				"reasoning":         response.Reasoning,
				"reasoning_content": response.ReasoningContent,
				"target_channel":    al.targetReasoningChannelID(opts.Channel),
				"channel":           opts.Channel,
			})
		// Check if no tool calls - we're done
		if len(response.ToolCalls) == 0 {
			finalContent = response.Content
			logger.InfoCF("agent", "LLM response without tool calls (direct answer)",
				map[string]any{
					"agent_id":      agent.ID,
					"iteration":     iteration,
					"content_chars": len(finalContent),
				})
			break
		}

		normalizedToolCalls := make([]providers.ToolCall, 0, len(response.ToolCalls))
		for _, tc := range response.ToolCalls {
			normalizedToolCalls = append(normalizedToolCalls, providers.NormalizeToolCall(tc))
		}

		// Log tool calls
		toolNames := make([]string, 0, len(normalizedToolCalls))
		for _, tc := range normalizedToolCalls {
			toolNames = append(toolNames, tc.Name)
		}
		logger.InfoCF("agent", "LLM requested tool calls",
			map[string]any{
				"agent_id":  agent.ID,
				"tools":     toolNames,
				"count":     len(normalizedToolCalls),
				"iteration": iteration,
			})

		// Build assistant message with tool calls
		assistantMsg := providers.Message{
			Role:             "assistant",
			Content:          response.Content,
			ReasoningContent: response.ReasoningContent,
		}
		for _, tc := range normalizedToolCalls {
			argumentsJSON, _ := json.Marshal(tc.Arguments)
			// Copy ExtraContent to ensure thought_signature is persisted for Gemini 3
			extraContent := tc.ExtraContent
			thoughtSignature := ""
			if tc.Function != nil {
				thoughtSignature = tc.Function.ThoughtSignature
			}

			assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, providers.ToolCall{
				ID:   tc.ID,
				Type: "function",
				Name: tc.Name,
				Function: &providers.FunctionCall{
					Name:             tc.Name,
					Arguments:        string(argumentsJSON),
					ThoughtSignature: thoughtSignature,
				},
				ExtraContent:     extraContent,
				ThoughtSignature: thoughtSignature,
			})
		}
		messages = append(messages, assistantMsg)

		// Save assistant message with tool calls to session
		agent.Sessions.AddFullMessage(opts.SessionKey, assistantMsg)

		// v0.2.1: Execute tool calls in parallel for better performance
		type indexedAgentResult struct {
			result *tools.ToolResult
			tc     providers.ToolCall
		}

		agentResults := make([]indexedAgentResult, len(normalizedToolCalls))
		var wg sync.WaitGroup

		for i, tc := range normalizedToolCalls {
			agentResults[i].tc = tc

			wg.Add(1)
			go func(idx int, tc providers.ToolCall) {
				defer wg.Done()

				argsJSON, _ := json.Marshal(tc.Arguments)
				argsPreview := utils.Truncate(string(argsJSON), 200)
				logger.InfoCF("agent", fmt.Sprintf("Tool call: %s(%s)", tc.Name, argsPreview),
					map[string]any{
						"agent_id":  agent.ID,
						"tool":      tc.Name,
						"iteration": iteration,
					})

				// Create async callback for tools that implement AsyncTool
				asyncCallback := func(callbackCtx context.Context, result *tools.ToolResult) {
					if !result.Silent && result.ForUser != "" {
						logger.InfoCF("agent", "Async tool completed, agent will handle notification",
							map[string]any{
								"tool":        tc.Name,
								"content_len": len(result.ForUser),
							})
					}
				}

				toolResult := agent.Tools.ExecuteWithContext(
					ctx,
					tc.Name,
					tc.Arguments,
					opts.Channel,
					opts.ChatID,
					asyncCallback,
				)
				agentResults[idx].result = toolResult
			}(i, tc)
		}
		wg.Wait()

		// Process results in original order (send to user, save to session)
		for _, r := range agentResults {
			// Send ForUser content to user immediately if not Silent
			if !r.result.Silent && r.result.ForUser != "" && opts.SendResponse {
				al.bus.PublishOutbound(ctx, bus.OutboundMessage{
					Channel: opts.Channel,
					ChatID:  opts.ChatID,
					Content: r.result.ForUser,
				})
				logger.DebugCF("agent", "Sent tool result to user",
					map[string]any{
						"tool":        r.tc.Name,
						"content_len": len(r.result.ForUser),
					})
			}

			// If tool returned media refs, publish them as outbound media
			if len(r.result.Media) > 0 && opts.SendResponse {
				parts := make([]bus.MediaPart, 0, len(r.result.Media))
				for _, ref := range r.result.Media {
					part := bus.MediaPart{Ref: ref}
					if al.mediaStore != nil {
						if _, meta, err := al.mediaStore.ResolveWithMeta(ref); err == nil {
							part.Filename = meta.Filename
							part.ContentType = meta.ContentType
							part.Type = inferMediaType(meta.Filename, meta.ContentType)
						}
					}
					parts = append(parts, part)
				}
				al.bus.PublishOutboundMedia(ctx, bus.OutboundMediaMessage{
					Channel: opts.Channel,
					ChatID:  opts.ChatID,
					Parts:   parts,
				})
			}

			// Determine content for LLM based on tool result
			contentForLLM := r.result.ForLLM
			if contentForLLM == "" && r.result.Err != nil {
				contentForLLM = r.result.Err.Error()
			}

			toolResultMsg := providers.Message{
				Role:       "tool",
				Content:    contentForLLM,
				ToolCallID: r.tc.ID,
			}
			messages = append(messages, toolResultMsg)

			// Save tool result message to session
			agent.Sessions.AddFullMessage(opts.SessionKey, toolResultMsg)
		}
	}

	return finalContent, iteration, nil
}

// updateToolContexts updates the context for tools that need channel/chatID info.
func (al *AgentLoop) updateToolContexts(agent *AgentInstance, channel, chatID string) {
	// Use ContextualTool interface instead of type assertions
	if tool, ok := agent.Tools.Get("message"); ok {
		if mt, ok := tool.(tools.ContextualTool); ok {
			mt.SetContext(channel, chatID)
		}
	}
	if tool, ok := agent.Tools.Get("spawn"); ok {
		if st, ok := tool.(tools.ContextualTool); ok {
			st.SetContext(channel, chatID)
		}
	}
	if tool, ok := agent.Tools.Get("subagent"); ok {
		if st, ok := tool.(tools.ContextualTool); ok {
			st.SetContext(channel, chatID)
		}
	}
}

// maybeSummarize triggers summarization if the session history exceeds thresholds.
func (al *AgentLoop) maybeSummarize(agent *AgentInstance, sessionKey, _, _ string) {
	newHistory := agent.Sessions.GetHistory(sessionKey)
	tokenEstimate := al.estimateTokens(newHistory)

	// v0.2.1: Use configurable thresholds from config
	messageThreshold := al.cfg.Session.SummarizationMessageThreshold
	if messageThreshold == 0 {
		messageThreshold = 20 // Fallback to default if not set
	}

	tokenPercent := al.cfg.Session.SummarizationTokenPercent
	if tokenPercent == 0 {
		tokenPercent = 0.75 // Fallback to default if not set
	}

	threshold := int(float64(agent.ContextWindow) * tokenPercent)

	if len(newHistory) > messageThreshold || tokenEstimate > threshold {
		summarizeKey := agent.ID + ":" + sessionKey
		if _, loading := al.summarizing.LoadOrStore(summarizeKey, true); !loading {
			go func() {
				defer func() {
					// Always clean up, even if panic occurs
					al.summarizing.Delete(summarizeKey)
					// Recover from panic to prevent goroutine crash
					if r := recover(); r != nil {
						logger.ErrorCF("agent", "Panic in summarization goroutine",
							map[string]any{
								"agent_id":    agent.ID,
								"session_key": sessionKey,
								"panic":       r,
							})
					}
				}()
				logger.DebugCF("agent", "Memory threshold reached. Optimizing conversation history...",
					map[string]any{
						"message_count":     len(newHistory),
						"message_threshold": messageThreshold,
						"token_estimate":    tokenEstimate,
						"token_threshold":   threshold,
					})
				al.summarizeSession(agent, sessionKey)
			}()
		}
	}
}

// forceCompression aggressively reduces context when the limit is hit.
// It drops the oldest 50% of messages (keeping system prompt and last user message).
func (al *AgentLoop) forceCompression(agent *AgentInstance, sessionKey string) {
	history := agent.Sessions.GetHistory(sessionKey)
	if len(history) <= 4 {
		return
	}

	// Keep system prompt (usually [0]) and the very last message (user's trigger)
	// We want to drop the oldest half of the *conversation*
	// Assuming [0] is system, [1:] is conversation
	conversation := history[1 : len(history)-1]
	if len(conversation) == 0 {
		return
	}

	// Helper to find the mid-point of the conversation
	mid := len(conversation) / 2

	// New history structure:
	// 1. System Prompt (with compression note appended)
	// 2. Second half of conversation
	// 3. Last message

	droppedCount := mid
	keptConversation := conversation[mid:]

	newHistory := make([]providers.Message, 0, 1+len(keptConversation)+1)

	// Strip any existing compression note from the system prompt to prevent accumulation
	systemContent := history[0].Content
	if idx := strings.LastIndex(systemContent, "\n\n[System Note: Emergency compression dropped"); idx >= 0 {
		systemContent = systemContent[:idx]
	}

	compressionNote := fmt.Sprintf(
		"\n\n[System Note: Emergency compression dropped %d oldest messages due to context limit]",
		droppedCount,
	)
	enhancedSystemPrompt := history[0]
	enhancedSystemPrompt.Content = systemContent + compressionNote
	newHistory = append(newHistory, enhancedSystemPrompt)

	newHistory = append(newHistory, keptConversation...)
	newHistory = append(newHistory, history[len(history)-1]) // Last message

	// Update session
	agent.Sessions.SetHistory(sessionKey, newHistory)
	agent.Sessions.Save(sessionKey)

	logger.WarnCF("agent", "Forced compression executed", map[string]any{
		"session_key":  sessionKey,
		"dropped_msgs": droppedCount,
		"new_count":    len(newHistory),
	})
}

// GetStartupInfo returns information about loaded tools and skills for logging.
func (al *AgentLoop) GetStartupInfo() map[string]any {
	info := make(map[string]any)

	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return info
	}

	// Tools info
	toolsList := agent.Tools.List()
	info["tools"] = map[string]any{
		"count": len(toolsList),
		"names": toolsList,
	}

	// Skills info
	info["skills"] = agent.ContextBuilder.GetSkillsInfo()

	// Agents info
	info["agents"] = map[string]any{
		"count": len(al.registry.ListAgentIDs()),
		"ids":   al.registry.ListAgentIDs(),
	}

	return info
}

// formatMessagesForLog formats messages for logging
func formatMessagesForLog(messages []providers.Message) string {
	if len(messages) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, msg := range messages {
		fmt.Fprintf(&sb, "  [%d] Role: %s\n", i, msg.Role)
		if len(msg.ToolCalls) > 0 {
			sb.WriteString("  ToolCalls:\n")
			for _, tc := range msg.ToolCalls {
				fmt.Fprintf(&sb, "    - ID: %s, Type: %s, Name: %s\n", tc.ID, tc.Type, tc.Name)
				if tc.Function != nil {
					fmt.Fprintf(
						&sb,
						"      Arguments: %s\n",
						utils.Truncate(tc.Function.Arguments, 200),
					)
				}
			}
		}
		if msg.Content != "" {
			content := utils.Truncate(msg.Content, 200)
			fmt.Fprintf(&sb, "  Content: %s\n", content)
		}
		if msg.ToolCallID != "" {
			fmt.Fprintf(&sb, "  ToolCallID: %s\n", msg.ToolCallID)
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]")
	return sb.String()
}

// formatToolsForLog formats tool definitions for logging
func formatToolsForLog(toolDefs []providers.ToolDefinition) string {
	if len(toolDefs) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, tool := range toolDefs {
		fmt.Fprintf(&sb, "  [%d] Type: %s, Name: %s\n", i, tool.Type, tool.Function.Name)
		fmt.Fprintf(&sb, "      Description: %s\n", tool.Function.Description)
		if len(tool.Function.Parameters) > 0 {
			fmt.Fprintf(
				&sb,
				"      Parameters: %s\n",
				utils.Truncate(fmt.Sprintf("%v", tool.Function.Parameters), 200),
			)
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// summarizeSession summarizes the conversation history for a session.
func (al *AgentLoop) summarizeSession(agent *AgentInstance, sessionKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	history := agent.Sessions.GetHistory(sessionKey)
	summary := agent.Sessions.GetSummary(sessionKey)

	// Keep last 4 messages for continuity
	if len(history) <= 4 {
		return
	}

	toSummarize := history[:len(history)-4]

	// Oversized Message Guard
	maxMessageTokens := agent.ContextWindow / 2
	validMessages := make([]providers.Message, 0)
	omitted := false

	for _, m := range toSummarize {
		if m.Role != "user" && m.Role != "assistant" {
			continue
		}
		msgTokens := len(m.Content) / 2
		if msgTokens > maxMessageTokens {
			omitted = true
			continue
		}
		validMessages = append(validMessages, m)
	}

	if len(validMessages) == 0 {
		return
	}

	// Multi-Part Summarization
	var finalSummary string
	if len(validMessages) > 10 {
		mid := len(validMessages) / 2
		part1 := validMessages[:mid]
		part2 := validMessages[mid:]

		s1, err1 := al.summarizeBatch(ctx, agent, part1, "")
		s2, err2 := al.summarizeBatch(ctx, agent, part2, "")

		if err1 != nil && err2 != nil {
			logger.ErrorCF("agent", "Failed to summarize both batches", nil)
			return // Abort to prevent WIPING OUT context
		}

		mergePrompt := fmt.Sprintf(
			"Merge these two conversation summaries into one cohesive summary:\n\n1: %s\n\n2: %s",
			s1,
			s2,
		)
		resp, err := agent.Provider.Chat(
			ctx,
			[]providers.Message{{Role: "user", Content: mergePrompt}},
			nil,
			agent.Model,
			map[string]any{
				"max_tokens":       1024,
				"temperature":      0.3,
				"prompt_cache_key": agent.ID,
			},
		)
		if err == nil && resp.Content != "" {
			finalSummary = resp.Content
		} else {
			finalSummary = strings.TrimSpace(s1 + "\n\n" + s2)
		}
	} else {
		summaryResult, err := al.summarizeBatch(ctx, agent, validMessages, summary)
		if err != nil || summaryResult == "" {
			logger.ErrorCF("agent", "Summarization failed", map[string]any{"error": err})
			return // Abort to prevent wipeout
		}
		finalSummary = summaryResult
	}

	if finalSummary == "" {
		return // Guard against empty summary wipeout
	}

	if omitted && finalSummary != "" {
		finalSummary += "\n[Note: Some oversized messages were omitted from this summary for efficiency.]"
	}

	agent.Sessions.SetSummary(sessionKey, finalSummary)
	agent.Sessions.TruncateHistory(sessionKey, 4)
	agent.Sessions.Save(sessionKey)
}

// summarizeBatch summarizes a batch of messages.
func (al *AgentLoop) summarizeBatch(
	ctx context.Context,
	agent *AgentInstance,
	batch []providers.Message,
	existingSummary string,
) (string, error) {
	var sb strings.Builder
	sb.WriteString(
		"Provide a concise summary of this conversation segment, preserving core context and key points.\n",
	)
	if existingSummary != "" {
		sb.WriteString("Existing context: ")
		sb.WriteString(existingSummary)
		sb.WriteString("\n")
	}
	sb.WriteString("\nCONVERSATION:\n")
	for _, m := range batch {
		fmt.Fprintf(&sb, "%s: %s\n", m.Role, m.Content)
	}
	prompt := sb.String()

	response, err := agent.Provider.Chat(
		ctx,
		[]providers.Message{{Role: "user", Content: prompt}},
		nil,
		agent.Model,
		map[string]any{
			"max_tokens":       1024,
			"temperature":      0.3,
			"prompt_cache_key": agent.ID,
		},
	)
	if err != nil {
		return "", err
	}
	return response.Content, nil
}

// estimateTokens estimates the number of tokens in a message list.
// Uses a safe heuristic of 2.5 characters per token to account for CJK and other
// overheads better than the previous 3 chars/token.
func (al *AgentLoop) estimateTokens(messages []providers.Message) int {
	totalChars := 0
	for _, m := range messages {
		totalChars += utf8.RuneCountInString(m.Content)
	}
	// 2.5 chars per token = totalChars * 2 / 5
	return totalChars * 2 / 5
}

func (al *AgentLoop) handleCommand(_ context.Context, msg bus.InboundMessage) (string, bool) {
	content := strings.TrimSpace(msg.Content)
	if !strings.HasPrefix(content, "/") {
		return "", false
	}

	parts := strings.Fields(content)
	if len(parts) == 0 {
		return "", false
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/show":
		if len(args) < 1 {
			return "Usage: /show [model|channel|agents]", true
		}
		switch args[0] {
		case "model":
			defaultAgent := al.registry.GetDefaultAgent()
			if defaultAgent == nil {
				return "No default agent configured", true
			}
			return fmt.Sprintf("Current model: %s", defaultAgent.Model), true
		case "channel":
			return fmt.Sprintf("Current channel: %s", msg.Channel), true
		case "agents":
			agentIDs := al.registry.ListAgentIDs()
			return fmt.Sprintf("Registered agents: %s", strings.Join(agentIDs, ", ")), true
		default:
			return fmt.Sprintf("Unknown show target: %s", args[0]), true
		}

	case "/list":
		if len(args) < 1 {
			return "Usage: /list [models|channels|agents]", true
		}
		switch args[0] {
		case "models":
			return "Available models: configured in config.json per agent", true
		case "channels":
			if al.channelManager == nil {
				return "Channel manager not initialized", true
			}
			channels := al.channelManager.GetEnabledChannels()
			if len(channels) == 0 {
				return "No channels enabled", true
			}
			return fmt.Sprintf("Enabled channels: %s", strings.Join(channels, ", ")), true
		case "agents":
			agentIDs := al.registry.ListAgentIDs()
			return fmt.Sprintf("Registered agents: %s", strings.Join(agentIDs, ", ")), true
		default:
			return fmt.Sprintf("Unknown list target: %s", args[0]), true
		}

	case "/switch":
		if len(args) < 3 || args[1] != "to" {
			return "Usage: /switch [model|channel] to <name>", true
		}
		target := args[0]
		value := args[2]

		switch target {
		case "model":
			defaultAgent := al.registry.GetDefaultAgent()
			if defaultAgent == nil {
				return "No default agent configured", true
			}
			oldModel := defaultAgent.Model
			defaultAgent.Model = value
			return fmt.Sprintf("Switched model from %s to %s", oldModel, value), true
		case "channel":
			if al.channelManager == nil {
				return "Channel manager not initialized", true
			}
			if _, exists := al.channelManager.GetChannel(value); !exists && value != "cli" {
				return fmt.Sprintf("Channel '%s' not found or not enabled", value), true
			}
			return fmt.Sprintf("Switched target channel to %s", value), true
		default:
			return fmt.Sprintf("Unknown switch target: %s", target), true
		}
	}

	return "", false
}

// searchVectorMemory searches for relevant past conversations in vector memory
func (al *AgentLoop) searchVectorMemory(ctx context.Context, query string, sessionKey string) string {
	if al.vectorStore == nil || al.embeddingService == nil {
		return ""
	}

	// Generate embedding for query
	embedding, err := al.embeddingService.Generate(ctx, query)
	if err != nil {
		logger.DebugCF("agent", "Failed to generate embedding for search",
			map[string]any{"error": err.Error()})
		return ""
	}

	// BUG FIX #10: Check if embedding is empty (len() for nil slices is defined as zero)
	if len(embedding) == 0 {
		logger.DebugCF("agent", "Empty embedding result", nil)
		return ""
	}

	// BUG FIX #11: Limit search results to prevent memory leak
	maxResults := 3
	if maxResults > 10 {
		maxResults = 10 // Hard limit
	}

	// Search similar vectors
	results, err := al.vectorStore.Search(ctx, memory.Vector{
		Embedding: embedding,
	}, maxResults)

	if err != nil {
		logger.DebugCF("agent", "Vector search failed",
			map[string]any{"error": err.Error()})
		return ""
	}

	if len(results) == 0 {
		return ""
	}

	// Build context from results
	var contextBuilder strings.Builder
	contextBuilder.WriteString("\n\n[Relevant past context from memory:]\n")
	for i, result := range results {
		if content, ok := result.Metadata["content"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("%d. %s (relevance: %.2f)\n",
				i+1, content, result.Score))
		}
	}

	logger.DebugCF("agent", "Retrieved memory context",
		map[string]any{
			"results": len(results),
			"session": sessionKey,
		})

	return contextBuilder.String()
}

// searchMemoryProvider searches for relevant past conversations in a memory provider (like Mem0)
func (al *AgentLoop) searchMemoryProvider(ctx context.Context, query string, sessionKey string) string {
	if al.memoryProvider == nil {
		return ""
	}

	// Limit search results
	maxResults := 3

	// Create a timeout context specifically for search
	searchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Search via provider
	results, err := al.memoryProvider.Recall(searchCtx, query, maxResults)
	if err != nil {
		logger.DebugCF("agent", "Memory provider search failed",
			map[string]any{"error": err.Error()})
		return ""
	}

	if len(results) == 0 {
		return ""
	}

	// Build context from results
	var contextBuilder strings.Builder
	for _, result := range results {
		contextBuilder.WriteString(fmt.Sprintf("- %s\n", result.Content))
	}

	logger.DebugCF("agent", "Retrieved provider memory context",
		map[string]any{
			"results": len(results),
			"session": sessionKey,
		})

	return contextBuilder.String()
}

// injectMemoryContext injects memory context into the system prompt
func (al *AgentLoop) injectMemoryContext(messages []providers.Message, context string) {
	if len(messages) == 0 {
		return
	}

	// Append to system prompt (first message)
	if messages[0].Role == "system" {
		messages[0].Content += context
	}
}

// extractMemoryFacts uses the LLM to extract permanent facts from a conversation transcript
func (al *AgentLoop) extractMemoryFacts(ctx context.Context, agent *AgentInstance, transcript string) ([]string, error) {
	if agent == nil || agent.Provider == nil {
		return nil, fmt.Errorf("no LLM provider configured on agent")
	}

	systemPrompt := `You are an AI Memory Extractor.
Ignore all pleasantries, greetings, and filler words.
Extract ONLY permanent, valuable facts (user traits, hardware, environment, core problems solved, decisions).
Rewrite each fact as a standalone, 3rd-person declarative sentence.
The output MUST be a strict JSON array of strings (e.g., ["fact 1", "fact 2"]).
Do NOT include any markdown formatting wrappers or conversational text.`

	messages := []providers.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: transcript},
	}

	// Use the configured model for the agent, fallback to provider default
	model := agent.Model
	if model == "" {
		model = agent.Provider.GetDefaultModel()
	}

	// Call the LLM with explicit options to prevent JSON truncation
	options := map[string]any{
		"max_tokens":  2048,
		"temperature": 0.1,
	}
	resp, err := agent.Provider.Chat(ctx, messages, nil, model, options)
	if err != nil {
		return nil, fmt.Errorf("LLM fact extraction failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("LLM returned nil response")
	}

	// Robust JSON parsing: locate array bounds to strip markdown
	content := strings.TrimSpace(resp.Content)
	startIdx := strings.Index(content, "[")
	endIdx := strings.LastIndex(content, "]")

	if startIdx == -1 || endIdx == -1 || startIdx > endIdx {
		return nil, fmt.Errorf("no valid JSON array found in LLM response: %s", utils.Truncate(content, 100))
	}

	jsonStr := content[startIdx : endIdx+1]

	var facts []string
	if err := json.Unmarshal([]byte(jsonStr), &facts); err != nil {
		return nil, fmt.Errorf("failed to parse facts JSON: %w (content: %s)", err, jsonStr)
	}

	return facts, nil
}

// storeInVectorMemory stores a conversation in vector memory (async)
func (al *AgentLoop) storeInVectorMemory(ctx context.Context, agent *AgentInstance, sessionKey, userMsg, assistantMsg, channel string) {
	if al.vectorStore == nil && al.memoryProvider == nil {
		return
	}

	// Hard skip for heartbeat messages to prevent memory pollution
	if sessionKey == "heartbeat" || assistantMsg == "HEARTBEAT_OK" || strings.Contains(userMsg, "Heartbeat Check") {
		return
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // Increased timeout to allow LLM extraction
	defer cancel()

	// Combine user and assistant messages for embedding
	combinedText := fmt.Sprintf("User: %s\nAssistant: %s", userMsg, assistantMsg)

	// Extract clean, standalone facts
	facts, err := al.extractMemoryFacts(ctx, agent, combinedText)
	if err != nil {
		logger.WarnCF("agent", "Failed to extract memory facts", map[string]any{"error": err.Error()})
		return
	}

	if len(facts) == 0 {
		logger.DebugCF("agent", "No extractable memory facts found in conversation turn", nil)
		return // Nothing valuable to store
	}

	var vectors []memory.Vector
	timestamp := time.Now().Unix()
	createdAt := time.Now().Format(time.RFC3339)

	for i, fact := range facts {
		// Generate embedding for each individual fact ONLY if vectorStore is enabled
		if al.vectorStore != nil && al.embeddingService != nil {
			embedding, err := al.embeddingService.Generate(ctx, fact)
			if err != nil {
				logger.WarnCF("agent", "Failed to generate embedding for fact", map[string]any{"error": err.Error(), "fact": utils.Truncate(fact, 50)})
			} else {
				// Prevent Duplicate Vectors: Search before Upsert
				// Create a temporary vector to search with
				queryVec := memory.Vector{
					Embedding: embedding,
				}

				// Search for Top 5 most similar vectors (in case there are cross-session matches)
				searchResults, searchErr := al.vectorStore.Search(ctx, queryVec, 5)

				// Check if the most similar vector is practically identical AND belongs to the same session
				isDuplicate := false
				if searchErr == nil && len(searchResults) > 0 {
					for _, res := range searchResults {
						// Handle both Qdrant (Cosine Similarity where 1.0 is exact) 
						// and LanceDB (L2/Cosine Distance where 0.0 is exact)
						isHighlySimilar := false
						if res.Score > 0.95 || res.Score < 0.05 {
							isHighlySimilar = true
						}

						if isHighlySimilar {
							// Only deduplicate if the fact belongs to the SAME session
							if resSession, ok := res.Metadata["session"].(string); ok && resSession == sessionKey {
								isDuplicate = true
								logger.DebugCF("agent", "Skipped duplicate fact (semantic deduplication)",
									map[string]any{
										"fact":       utils.Truncate(fact, 50),
										"score":      res.Score,
										"matched_id": res.ID,
									})
								break
							}
						}
					}
				} else if searchErr != nil {
					logger.WarnCF("agent", "Failed to check duplicates in vector store", map[string]any{"error": searchErr.Error()})
				}

				if !isDuplicate {
					vec := memory.Vector{
						ID:        fmt.Sprintf("%s:%d:%d", sessionKey, time.Now().UnixNano(), i),
						Embedding: embedding,
						Metadata: map[string]interface{}{
							"session":    sessionKey,
							"channel":    channel,
							"timestamp":  timestamp,
							"user_msg":   userMsg,
							"content":    fact,
							"created_at": createdAt,
							"original":   combinedText,
						},
					}
					vectors = append(vectors, vec)
				}
			}
		}

		// Store in memory provider (Mem0) if enabled
		if al.memoryProvider != nil {
			_, err := al.memoryProvider.Store(ctx, fact, map[string]interface{}{
				"session":  sessionKey,
				"channel":  channel,
				"original": combinedText,
			})
			if err != nil {
				logger.WarnCF("agent", "Failed to store fact in memory provider",
					map[string]any{"error": err.Error(), "fact": utils.Truncate(fact, 50)})
			}
		}
	}

	if len(vectors) == 0 && al.memoryProvider != nil {
		logger.DebugCF("agent", "Stored facts in memory provider", map[string]any{"count": len(facts)})
		return
	} else if len(vectors) == 0 {
		return
	}

	// Upsert to vector store
	if al.vectorStore != nil {
		err = al.vectorStore.Upsert(ctx, vectors)
		if err != nil {
			logger.WarnCF("agent", "Failed to store in vector memory",
				map[string]any{"error": err.Error()})
			return
		}
	}

	logger.DebugCF("agent", "Stored extracted facts in vector memory",
		map[string]any{
			"session": sessionKey,
			"count":   len(vectors),
		})
}

// extractPeer extracts the routing peer from the inbound message's structured Peer field.
func extractPeer(msg bus.InboundMessage) *routing.RoutePeer {
	if msg.Peer.Kind == "" {
		return nil
	}
	peerID := msg.Peer.ID
	if peerID == "" {
		if msg.Peer.Kind == "direct" {
			peerID = msg.SenderID
		} else {
			peerID = msg.ChatID
		}
	}
	return &routing.RoutePeer{Kind: msg.Peer.Kind, ID: peerID}
}

// extractParentPeer extracts the parent peer (reply-to) from inbound message metadata.
func extractParentPeer(msg bus.InboundMessage) *routing.RoutePeer {
	parentKind := msg.Metadata["parent_peer_kind"]
	parentID := msg.Metadata["parent_peer_id"]
	if parentKind == "" || parentID == "" {
		return nil
	}
	return &routing.RoutePeer{Kind: parentKind, ID: parentID}
}
