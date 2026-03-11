# Còn Gì Để Integrate từ v0.2.1?

**Date**: 2026-03-09  
**Status**: Đã giải quyết xung đột, sẵn sàng cho integration tiếp theo

## ✅ Đã Có Sẵn (Không cần integrate)

### 1. MCP Tools Support ✅
- **Status**: ĐÃ CÓ trong codebase
- **Files**: `pkg/mcp/manager.go`, `pkg/tools/mcp_tool.go`
- **Action**: Không cần làm gì

### 2. Vision/Image Support ✅
- **Status**: ĐÃ CÓ `resolveMediaRefs()` trong `pkg/agent/loop_media.go`
- **Files**: `pkg/agent/loop_media.go`, `pkg/providers/types.go`
- **Action**: Không cần làm gì

### 3. Parallel Tool Execution ✅
- **Status**: ĐÃ DÙNG v0.2.1's inline implementation
- **Files**: `pkg/agent/loop.go`
- **Action**: Đã resolved conflict

### 4. Tool Enable/Disable ✅
- **Status**: ĐÃ IMPLEMENT
- **Files**: `pkg/config/config.go`, `pkg/agent/loop.go`
- **Action**: Cần verify với v0.2.1

### 5. Configurable Summarization ✅
- **Status**: ĐÃ IMPLEMENT
- **Files**: `pkg/config/config.go`, `pkg/agent/loop.go`
- **Action**: Cần verify với v0.2.1

### 6. .env File Loading ✅
- **Status**: ĐÃ IMPLEMENT
- **Files**: `pkg/config/env.go`, `.env.example`
- **Action**: Cần compare với v0.2.1

---

## 🔥 CÒN CÓ THỂ INTEGRATE (High Value)

### 1. ⭐ JSONL Memory Store (#732)
**Priority**: 🔥 HIGHEST  
**Effort**: High  
**Impact**: Critical for data safety

**Mô tả**:
- Thay thế JSON sessions bằng JSONL append-only format
- Crash-safe session storage với fsync
- Physical compaction để giảm disk usage
- Migration từ legacy JSON sessions

**Lợi ích**:
- ✅ Crash-safe - Không mất data khi crash
- ✅ Better performance với large sessions
- ✅ Disk space optimization với compaction
- ✅ Backward compatible với migration

**Implementation**:
```
Files cần tạo:
- pkg/memory/jsonl/store.go (new)
- pkg/memory/jsonl/compaction.go (new)
- pkg/memory/jsonl/migration.go (new)

Files cần sửa:
- pkg/session/manager.go (use Store interface)
- pkg/config/config.go (add JSONL config)
```

**Commits tham khảo**:
- `c8178f4` - Main PR
- `b464687` - Compact method
- `32ec8ca` - Store interface
- `9f36e50` - JSONL implementation
- `9036812` - Migration support
- `f9f726c` - fsync for durability
- `529622b` - Tests

**Estimated Time**: 2-3 days

---

### 2. ⭐ Model Routing (#994)
**Priority**: 🔥 HIGH  
**Effort**: Medium  
**Impact**: Cost optimization + Performance

**Mô tả**:
- Tự động route messages đến models phù hợp dựa trên complexity
- Language-agnostic complexity scorer
- CJK character support

**Lợi ích**:
- ✅ Cost savings (dùng cheap models cho simple queries)
- ✅ Performance optimization
- ✅ Smart model selection

**Implementation**:
```
Files cần tạo:
- pkg/routing/complexity.go (new)
- pkg/routing/scorer.go (new)

Files cần sửa:
- pkg/agent/loop.go (integrate routing)
- pkg/config/config.go (add RoutingConfig)
```

**Commits tham khảo**:
- `9b1e73d` - Main PR
- `1943c3e` - Complexity scorer
- `02e8192` - Wire into agent loop
- `c5a21b2` - RoutingConfig
- `b84adac` - CJK support

**Estimated Time**: 1-2 days

---

### 3. ⭐ Extended Thinking Support (#1076)
**Priority**: 🔥 MEDIUM  
**Effort**: Low  
**Impact**: Better reasoning với Claude

**Mô tả**:
- Hỗ trợ extended thinking mode cho Anthropic models
- Preserve reasoning_content in history
- Fallback to reasoning content

**Lợi ích**:
- ✅ Better reasoning với Claude models
- ✅ Preserve thinking process
- ✅ Improved multi-turn conversations

**Implementation**:
```
Files cần sửa:
- pkg/providers/anthropic.go (add extended thinking)
- pkg/session/manager.go (preserve reasoning)
- pkg/agent/loop.go (fallback logic)
```

**Commits tham khảo**:
- `204038e` - Main feature
- `a4e5c39` - Preserve reasoning_content
- `9efdde2` - Multi-turn preservation
- `66e6fb6` - Fallback logic

**Estimated Time**: 0.5-1 day

---

### 4. 🔍 New Search Providers
**Priority**: ⚠️ MEDIUM  
**Effort**: Low per provider  
**Impact**: More search options

**Providers mới**:
- **SearXNG** (#534) - Privacy-focused metasearch
- **GLM Search** (#1057) - 智谱 Chinese search
- **Exa AI** (#4b7e8d9) - AI-powered search

**Lợi ích**:
- ✅ More search options
- ✅ Privacy-focused (SearXNG)
- ✅ Chinese market support (GLM)
- ✅ AI-powered search (Exa)

**Implementation**:
```
Files cần tạo:
- pkg/tools/searxng.go (new)
- pkg/tools/glm_search.go (new)
- pkg/tools/exa_search.go (new)

Files cần sửa:
- pkg/config/config.go (add configs)
- pkg/agent/loop.go (register tools)
```

**Estimated Time**: 0.5 day per provider

---

### 5. 🏠 PICOCLAW_HOME Environment Variable (#1155)
**Priority**: ⚠️ LOW  
**Effort**: Very Low  
**Impact**: Convenience

**Mô tả**:
- Honor PICOCLAW_HOME env var cho config/auth/workspace paths
- Flexible installation paths
- Multi-user support

**Lợi ích**:
- ✅ Flexible installation
- ✅ Multi-user support
- ✅ Better DevOps

**Implementation**:
```
Files cần sửa:
- pkg/config/defaults.go (check PICOCLAW_HOME)
- pkg/auth/*.go (use PICOCLAW_HOME)
```

**Commits tham khảo**:
- `651cb2e` - Main PR
- `51e8479` - Implementation

**Estimated Time**: 0.5 day

---

## 📱 Optional: New Channels (Nếu cần)

### IRC Channel (#1138)
- **Effort**: Medium
- **Use Case**: IRC communities
- **Files**: `pkg/channels/irc/`

### Matrix Channel (#1220)
- **Effort**: Medium
- **Use Case**: Matrix/Element users
- **Files**: `pkg/channels/matrix/`

### WeCom AI Bot (#6caee42)
- **Effort**: High
- **Use Case**: Enterprise WeCom
- **Files**: `pkg/channels/wecom/aibot.go`

---

## 🌐 Optional: New LLM Providers (Nếu cần)

### Minimax (#1273)
- Chinese LLM provider
- **Effort**: Low

### Avian (#844)
- Data analytics provider
- **Effort**: Low

### Kimi/Moonshot (#ec54031)
- Chinese LLM provider
- **Effort**: Low

---

## 📊 Recommended Priority Order

### Phase 1: Critical (Week 1-2)
1. **JSONL Memory Store** ⭐⭐⭐
   - Crash safety là critical
   - Foundation cho production use
   - **Effort**: 2-3 days

2. **Model Routing** ⭐⭐
   - Cost optimization
   - Performance improvement
   - **Effort**: 1-2 days

### Phase 2: Important (Week 3)
3. **Extended Thinking** ⭐
   - Better reasoning
   - Easy to implement
   - **Effort**: 0.5-1 day

4. **New Search Providers** ⚠️
   - More options
   - Easy to add incrementally
   - **Effort**: 0.5 day each

### Phase 3: Nice-to-Have (Week 4+)
5. **PICOCLAW_HOME** ⚠️
   - Convenience feature
   - Very easy
   - **Effort**: 0.5 day

6. **New Channels/Providers** (as needed)
   - Only if you need them
   - **Effort**: Varies

---

## 🎯 Recommended Action Plan

### This Week (Week 1)
- [ ] **Day 1-2**: Verify existing implementations
  - Tool Enable/Disable vs v0.2.1
  - Configurable Summarization vs v0.2.1
  - .env Loading vs v0.2.1

- [ ] **Day 3-5**: Implement JSONL Memory Store
  - Study v0.2.1 implementation
  - Create Store interface
  - Implement JSONL store
  - Add migration logic
  - Write tests

### Next Week (Week 2)
- [ ] **Day 1-2**: Implement Model Routing
  - Create complexity scorer
  - Add routing logic
  - Configure model tiers
  - Test routing

- [ ] **Day 3**: Implement Extended Thinking
  - Add reasoning_content support
  - Preserve in history
  - Add fallback logic

- [ ] **Day 4-5**: Add Search Providers
  - SearXNG
  - GLM Search (if needed)
  - Exa AI (if needed)

### Week 3+
- [ ] PICOCLAW_HOME support
- [ ] New channels (if needed)
- [ ] New providers (if needed)
- [ ] Testing and documentation

---

## 💡 Key Insights

### What We Learned
1. **Always check upstream first** - Saved us from more duplicate work
2. **v0.2.1 is feature-rich** - Most features already implemented
3. **Focus on high-value features** - JSONL and Model Routing are most important

### What's Most Valuable
1. **JSONL Memory Store** - Critical for production reliability
2. **Model Routing** - Significant cost savings
3. **Extended Thinking** - Better AI reasoning

### What's Optional
1. New channels - Only if you need specific platforms
2. New providers - Only if you need specific LLMs
3. Search providers - Nice to have, not critical

---

## 📋 Summary

### Must Integrate (High ROI)
1. ✅ **JSONL Memory Store** - Crash safety + performance
2. ✅ **Model Routing** - Cost optimization
3. ✅ **Extended Thinking** - Better reasoning

### Should Integrate (Medium ROI)
4. ⚠️ **New Search Providers** - More options
5. ⚠️ **PICOCLAW_HOME** - Convenience

### Optional (Low ROI unless needed)
6. ❌ New channels - Platform-specific
7. ❌ New providers - Provider-specific

---

## 🎯 Next Steps

### Immediate
1. Verify existing implementations match v0.2.1
2. Plan JSONL Memory Store implementation
3. Study v0.2.1's Model Routing

### Short Term
1. Implement JSONL Memory Store
2. Implement Model Routing
3. Add Extended Thinking

### Long Term
1. Add search providers as needed
2. Add channels/providers as needed
3. Continuous integration with upstream

---

**Total Estimated Time**: 4-6 days for high-priority features  
**Highest Value**: JSONL Memory Store (crash safety)  
**Easiest Win**: Extended Thinking (0.5-1 day)  
**Best ROI**: Model Routing (cost savings)

