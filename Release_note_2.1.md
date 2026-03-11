# PicoClaw Release Notes

## Overview
This release brings massive additions to PicoClaw, integrating advanced memory management, multi-modal capabilities, search enhancements, and enterprise-grade resilience.

## 🚀 Major Features

### 1. Vector Memory System & Resiliency
- **Qdrant Integration:** Production-ready vector storage with automatic retry logic, context-aware timeouts, and a robust `QdrantConfig` schema.
- **LanceDB Integration:** Embedded, CGO-based local vector storage without requiring an external server.
- **Circuit Breaker Pattern:** Automatic failure detection and recovery (5 failures / 30s reset / 3 test requests) to protect system stability during downstream failures.
- **Exponential Backoff:** Intelligent retry logic to handle transient network issues smoothly.

### 2. Advanced Memory Providers
- **Mem0 Integration:** Personalized memory provider supporting store, recall, update, and delete operations with token authentication and circuit breaker protection.
- **MindGraph Integration:** Knowledge graph memory implementation supporting robust data operations, bearer token authentication, and cloud/self-hosted deployments.

### 3. Multi-Modal Vision Support
- **Image Understanding:** Streaming base64 encoding for high memory efficiency.
- **Model Support:** Integrated OpenAI Vision (GPT-4o, GPT-4V) and Anthropic Vision (Claude 3+).
- **Format Handling:** Automatic MIME type detection with configurable max file size (default 20MB).

### 4. Extended Thinking & Routing
- **Extended Thinking:** Extracts thinking blocks from Claude 3.5+ responses (Anthropic `reasoning_content`) and preserves them in session history.
- **Complexity-Based Model Routing:** Language-agnostic feature extraction to automatically select between cheap, medium, and expensive models for cost optimization.

### 5. Enhanced Search Capabilities
- **New Search Providers:** Added SearXNG (privacy-focused metasearch), GLM Search (Chinese search by Zhipu AI), and Exa AI (semantic search with autoprompt).
- **Priority Chain:** Perplexity > Exa > GLM > Brave > Tavily > SearXNG > DuckDuckGo.

## 🛠️ System & Architecture Improvements

- **PICOCLAW_HOME:** Custom home directory configuration via environment variables, supporting multi-tenant/multi-user deployments and Docker environments (fallback to `~/.picoclaw`).
- **Granular Tool Control:** Individual feature flags for file, shell, web, message, spawn, team, skill, and hardware tools.
- **Safe State Management:** Default storage backend upgraded from JSON to JSONL for crash-safe, append-only session tracking.
- **Environment Parity:** Improved `.env` file support taking precedence over `config.json`.
- **Concurrency & Idempotency Fixes:** Addressed context stacking leaks, implemented `TryMarkDispatched()` for atomic message tracking, and resolved various race conditions under high load.

## 📊 Key Metrics
- Optimized for **<50MB RAM budget** under deep multi-agent cascades.
- Search latency p50 **<100ms**.
- Test coverage expanded to **~85%** across core memory modules.
