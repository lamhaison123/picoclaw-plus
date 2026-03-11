# PicoClaw v2.0.5 Release Notes (Turbo Patch)

**Release Date**: March 7, 2026  
**Type**: Critical Performance & Depth Upgrade  
**Status**: ✅ CODEBASE-VERIFIED (Static Analysis + Code Audit)

---

## 🎯 Overview

Version 2.0.5 (Turbo Patch) is a major upgrade to the multi-agent collaboration core. It introduces support for flexible conversation depths, a robust idempotency system, and improved resource management to maintain PicoClaw's signature low-memory footprint even under heavy load.

---

## ✨ What's New

### 🚀 Flexible Depth Support (Turbo Mode)

**Increased Conversation Depth:**
- **Default Max Depth**: Increased from 3 to **20 levels**.
- **Configurable**: Users can adjust `max_mention_depth` in configuration to suit their hardware constraints.
- **Context Preservation**: Optimized context chain to ensure rich information flow even at high depths.

### 🛡️ Resource Management (Safety v2.0)

**To support higher depths while keeping RAM <50MB:**
- **Per-Role Queues**: Buffered channels of **20 tasks per role** to prevent resource exhaustion.
- **Rate Limiting**: Integrated **2-second cooldown** per role to prevent API spam and manage load.
- **Cycle Detection**: Advanced cascade tracking to prevent infinite mention loops.

### 🔄 Idempotency & Reliability Fixes

**Fixed critical P0/P1 issues:**
- **Idempotency System**: Implemented a robust **Check-Execute-Mark** pattern using `atomic` flags and `sync.Map` to prevent duplicate message execution.
- **Retry Mechanism**: Improved retry logic with **3 attempts** and **exponential backoff** (1s base).
- **Timeout Hierarchy**: Implemented a progressive timeout strategy to ensure clean context propagation.

### 🧹 Infrastructure Improvements

**Environment Bypass:**
- Successfully implemented `/tmp` redirection for Go toolchain (`GOPATH`, `GOCACHE`) to support read-only filesystems in sandboxed environments.

---

## 📊 Statistics (v2.0.5 vs v1.2.0)

| Metric | v1.2.0 | **v2.0.5 (Turbo)** |
|--------|--------|--------------------|
| Max Depth | 3 | **20 (Default)** |
| RAM Budget | 10MB | **50MB (Active)** |
| Queue Size | 10 | **20 per role** |
| Idempotency | Basic | **Verified (Atomic)** |
| Rate Limit | None | **2 Seconds** |

---

## 🚀 Upgrade Guide

### Configuration Changes

Update your `config.yaml` or `config.json` to leverage new features:

```json
{
  "channels": {
    "telegram": {
      "collaborative_chat": {
        "max_mention_depth": 20,
        "mention_queue_size": 20,
        "mention_rate_limit": "2s",
        "mention_max_retries": 3
      }
    }
  }
}
```

---

## 🔒 Verification Status

### ✅ CODEBASE-VERIFIED
- Comprehensive static analysis completed.
- Thread-safety patterns (Atomic + Mutex) verified against `pkg/collaborative/manager_improved.go`.
- Idempotency logic (Check-Execute-Mark) confirmed in `pkg/collaborative/dispatch.go`.
- Integration chain (Manager→Tracker→Metrics) validated.

---

## 🙏 Credits

### Contributors
- @manager: Coordination & Release Management
- @developer: Turbo Patch Implementation & Environment Bypass
- @architect: Concurrency Design & Safety Architecture
- @tester: Depth Verification & Metrics Validation

---

**Happy coding! 🦞**
