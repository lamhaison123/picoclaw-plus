# Integration Checklist - PicoClaw v0.2.1 Features

## 🎯 Phase 1: Critical Features (Week 1-2)

### 1. JSONL Memory Store ⭐ HIGH PRIORITY
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 JSONL implementation
  - [ ] Read `pkg/memory/jsonl/store.go`
  - [ ] Understand Store interface
  - [ ] Review compaction logic
  - [ ] Check migration code

- [ ] Create JSONL store implementation
  - [ ] Create `pkg/memory/jsonl/` directory
  - [ ] Implement `Store` interface
  - [ ] Add append-only write
  - [ ] Add fsync for durability
  - [ ] Implement compaction

- [ ] Update session manager
  - [ ] Modify `pkg/session/manager.go`
  - [ ] Use Store interface
  - [ ] Add migration from JSON
  - [ ] Test backward compatibility

- [ ] Testing
  - [ ] Unit tests for JSONL store
  - [ ] Concurrency tests
  - [ ] Crash recovery tests
  - [ ] Migration tests
  - [ ] Benchmark tests

**Files to Create**:
- `pkg/memory/jsonl/store.go`
- `pkg/memory/jsonl/store_test.go`
- `pkg/memory/jsonl/compaction.go`
- `pkg/memory/jsonl/migration.go`

**Files to Modify**:
- `pkg/session/manager.go`
- `pkg/config/config.go` (add JSONL config)

**Estimated Time**: 3-4 days

---

### 2. Vision/Image Support ⭐ HIGH PRIORITY
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 vision implementation
  - [ ] Check Media field in Message struct
  - [ ] Review resolveMediaRefs function
  - [ ] Check base64 encoding logic
  - [ ] Review filetype detection

- [ ] Update provider types
  - [ ] Add `Media` field to `providers.Message`
  - [ ] Update message serialization
  - [ ] Support vision models

- [ ] Implement media resolution
  - [ ] Add `resolveMediaRefs()` to agent loop
  - [ ] Implement base64 streaming
  - [ ] Add filetype detection
  - [ ] Handle media:// refs

- [ ] Update providers
  - [ ] OpenAI: Add vision support
  - [ ] Anthropic: Add vision support
  - [ ] Other providers as needed

- [ ] Testing
  - [ ] Unit tests for media resolution
  - [ ] Integration tests with vision models
  - [ ] Test different image formats
  - [ ] Test large images

**Files to Create**:
- `pkg/media/resolver.go` (if not exists)
- `pkg/media/resolver_test.go`

**Files to Modify**:
- `pkg/providers/types.go`
- `pkg/providers/openai_compat.go`
- `pkg/providers/anthropic.go`
- `pkg/agent/loop.go`

**Estimated Time**: 2-3 days

---

### 3. Parallel Tool Execution ⭐ HIGH PRIORITY
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 parallel execution
  - [ ] Check goroutine pool implementation
  - [ ] Review error handling
  - [ ] Check result aggregation

- [ ] Design parallel execution
  - [ ] Plan goroutine pool
  - [ ] Design error handling
  - [ ] Plan result collection

- [ ] Implement parallel execution
  - [ ] Modify `runLLMIteration()` in agent loop
  - [ ] Add goroutine pool
  - [ ] Implement concurrent tool calls
  - [ ] Add timeout handling
  - [ ] Aggregate results

- [ ] Handle edge cases
  - [ ] Tool dependencies
  - [ ] Shared state
  - [ ] Error propagation
  - [ ] Context cancellation

- [ ] Testing
  - [ ] Unit tests for parallel execution
  - [ ] Race condition tests
  - [ ] Performance benchmarks
  - [ ] Timeout tests

**Files to Modify**:
- `pkg/agent/loop.go` (runLLMIteration)
- `pkg/tools/registry.go` (if needed)

**Estimated Time**: 2-3 days

---

### 4. Model Routing ⭐ HIGH PRIORITY
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 routing implementation
  - [ ] Check complexity scorer
  - [ ] Review CJK support
  - [ ] Check model tier logic

- [ ] Create routing module
  - [ ] Create `pkg/routing/complexity.go`
  - [ ] Implement complexity scorer
  - [ ] Add CJK character support
  - [ ] Add model tier configuration

- [ ] Integrate into agent loop
  - [ ] Add routing logic before LLM call
  - [ ] Select model based on complexity
  - [ ] Add fallback logic
  - [ ] Log routing decisions

- [ ] Configuration
  - [ ] Add RoutingConfig to config
  - [ ] Define model tiers (cheap/expensive)
  - [ ] Add complexity thresholds

- [ ] Testing
  - [ ] Unit tests for complexity scorer
  - [ ] Test CJK character handling
  - [ ] Test model selection
  - [ ] Integration tests

**Files to Create**:
- `pkg/routing/complexity.go`
- `pkg/routing/complexity_test.go`
- `pkg/routing/model_selector.go`

**Files to Modify**:
- `pkg/config/config.go`
- `pkg/agent/loop.go`

**Estimated Time**: 2-3 days

---

## 🎨 Phase 2: Important Features (Week 3-4)

### 5. Environment Variable Configuration
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 env implementation
  - [ ] Check .env file loading
  - [ ] Review env var overrides
  - [ ] Check provider env support

- [ ] Implement .env loading
  - [ ] Create `pkg/config/env.go`
  - [ ] Add .env file parser
  - [ ] Implement env var overrides
  - [ ] Add validation

- [ ] Update config loading
  - [ ] Load .env before config
  - [ ] Override config with env vars
  - [ ] Support provider-specific env vars

- [ ] Documentation
  - [ ] Document env var names
  - [ ] Add .env.example
  - [ ] Update configuration guide

- [ ] Testing
  - [ ] Unit tests for env loading
  - [ ] Test override precedence
  - [ ] Test malformed .env

**Files to Create**:
- `pkg/config/env.go`
- `pkg/config/env_test.go`
- `.env.example`

**Files to Modify**:
- `pkg/config/config.go`

**Estimated Time**: 1-2 days

---

### 6. Tool Enable/Disable Configuration
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 implementation
  - [ ] Check enabled flag usage
  - [ ] Review config structure

- [ ] Add enabled flags
  - [ ] Add `enabled` field to tool configs
  - [ ] Update config structs
  - [ ] Add default values

- [ ] Implement enable/disable logic
  - [ ] Check enabled flag before tool execution
  - [ ] Skip disabled tools in registration
  - [ ] Log disabled tools

- [ ] Update configuration
  - [ ] Add enabled flags to config.example.json
  - [ ] Document in tools_configuration.md

- [ ] Testing
  - [ ] Test tool enable/disable
  - [ ] Test default behavior
  - [ ] Test config validation

**Files to Modify**:
- `pkg/config/config.go`
- `pkg/tools/*.go` (各个 tool files)
- `config/config.example.json`
- `docs/reference/tools_configuration.md`

**Estimated Time**: 1 day

---

### 7. Extended Thinking Support (Anthropic)
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 implementation
  - [ ] Check reasoning_content handling
  - [ ] Review history preservation

- [ ] Add reasoning_content support
  - [ ] Update Message struct
  - [ ] Add to Anthropic provider
  - [ ] Preserve in session history

- [ ] Update session manager
  - [ ] Store reasoning_content
  - [ ] Load reasoning_content
  - [ ] Include in context building

- [ ] Testing
  - [ ] Test with Claude models
  - [ ] Test reasoning preservation
  - [ ] Test multi-turn conversations

**Files to Modify**:
- `pkg/providers/types.go`
- `pkg/providers/anthropic.go`
- `pkg/session/manager.go`
- `pkg/agent/context.go`

**Estimated Time**: 1 day

---

### 8. Configurable Summarization Thresholds
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 implementation
  - [ ] Check config fields
  - [ ] Review threshold logic

- [ ] Add config fields
  - [ ] Add message_threshold
  - [ ] Add token_percent
  - [ ] Add default values

- [ ] Update summarization logic
  - [ ] Use config thresholds
  - [ ] Add logging
  - [ ] Test different thresholds

- [ ] Documentation
  - [ ] Document config options
  - [ ] Add examples
  - [ ] Explain thresholds

- [ ] Testing
  - [ ] Test different thresholds
  - [ ] Test edge cases
  - [ ] Performance tests

**Files to Modify**:
- `pkg/config/config.go`
- `pkg/agent/loop.go` (maybeSummarize)
- `config/config.example.json`

**Estimated Time**: 0.5 day

---

## 🔧 Phase 3: Nice-to-Have (Week 5+)

### 9. New Search Providers (Optional)
**Status**: ⬜ Not Started

**Tasks**:
- [ ] SearXNG provider
  - [ ] Create `pkg/tools/searxng.go`
  - [ ] Implement search logic
  - [ ] Add tests

- [ ] GLM Search provider
  - [ ] Create `pkg/tools/glm_search.go`
  - [ ] Implement search logic
  - [ ] Add tests

- [ ] Exa AI provider
  - [ ] Create `pkg/tools/exa_search.go`
  - [ ] Implement search logic
  - [ ] Add tests

**Estimated Time**: 1-2 days per provider

---

### 10. PICOCLAW_HOME Environment Variable
**Status**: ⬜ Not Started

**Tasks**:
- [ ] Review v0.2.1 implementation
- [ ] Add PICOCLAW_HOME support
- [ ] Update path resolution
- [ ] Test multi-user scenarios
- [ ] Document usage

**Estimated Time**: 0.5 day

---

## 📊 Progress Tracking

### Overall Progress
- Phase 1: ⬜⬜⬜⬜ 0/4 (0%)
- Phase 2: ⬜⬜⬜⬜ 0/4 (0%)
- Phase 3: ⬜⬜ 0/2 (0%)
- **Total**: 0/10 (0%)

### Time Estimates
- Phase 1: 9-13 days
- Phase 2: 4-5 days
- Phase 3: 2-3 days
- **Total**: 15-21 days (3-4 weeks)

---

## 🎯 Success Criteria

### Phase 1 Complete When:
- [ ] JSONL store working with migration
- [ ] Vision support working with GPT-4V/Claude 3
- [ ] Parallel tool execution 2x faster
- [ ] Model routing saving costs
- [ ] All tests passing
- [ ] Documentation updated

### Phase 2 Complete When:
- [ ] .env file loading working
- [ ] Tools can be enabled/disabled
- [ ] Extended thinking working with Claude
- [ ] Summarization thresholds configurable
- [ ] All tests passing
- [ ] Documentation updated

### Phase 3 Complete When:
- [ ] New search providers working (if implemented)
- [ ] PICOCLAW_HOME working (if implemented)
- [ ] All tests passing
- [ ] Documentation updated

---

## 📝 Notes

### Dependencies
- Phase 2 can start after Phase 1 item 1 (JSONL) is done
- Phase 3 can start anytime (independent)

### Risks
- JSONL migration might have edge cases
- Vision support needs testing with multiple models
- Parallel execution might have race conditions
- Model routing needs tuning

### Mitigation
- Extensive testing for JSONL
- Test with multiple vision models
- Use race detector for parallel execution
- Monitor routing decisions in production

---

**Created**: 2026-03-09  
**Target Completion**: 2026-04-06 (4 weeks)  
**Last Updated**: 2026-03-09
