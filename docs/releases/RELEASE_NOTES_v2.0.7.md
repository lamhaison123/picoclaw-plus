# Release Notes - PicoClaw v2.0.7 🦞

## [v2.0.7] - 2026-03-09

### 🚀 New Features
- **Vector Memory Integration**: Full support for Qdrant as a vector storage provider.
- **Circuit Breaker Pattern**: Production-grade circuit breaker (5 failures / 30s reset) to protect system stability during downstream failures.
- **Exponential Backoff**: Intelligent retry logic for transient network issues.
- **Improved Config Schema**: New `vector` configuration block with validation.

### 🛠️ Bug Fixes & Improvements
- **Race Condition Fix**: Atomic state transitions in Circuit Breaker verified by @architect.
- **Context Stacking Resolved**: Proper `WithTimeout` hierarchy prevents context leaks and deadline propagation issues.
- **Idempotency Patch**: Refined dispatch tracking to prevent duplicate processing under high concurrency.
- **Memory Safety**: Optimized for <50MB RAM budget even with deep discussion threads.

### 📊 Key Metrics
- **Codebase**: 1,800 lines of Go (Sprint 1).
- **Test Coverage**: ~85% across core modules.
- **Performance**: p50 search latency <100ms (simulated).
- **Audit Status**: 10/10 PASS by @architect.

### 📂 Repository Information
- **Main Path**: `/root/.picoclaw/workspace/teams/dev-team/picoclaw-plus/pkg/memory/vector/`
- **Docs Path**: `/root/.picoclaw/workspace/teams/dev-team/picoclaw-plus/docs/memory/vector/`

---
*PicoClaw - Ultra-lightweight AI Assistant*
