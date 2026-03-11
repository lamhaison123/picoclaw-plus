# v0.2.1 Integration Verification Checklist

## 🎯 Purpose
Verify that all v0.2.1 integrated features work correctly in production.

## ✅ Pre-Deployment Checks

### Build & Compilation
- [x] `go build` completes without errors
- [x] No diagnostics errors in modified files
- [x] All imports resolved correctly
- [x] Binary size reasonable (~50-100MB)

### Code Quality
- [x] No race conditions (verified with sync.Map, atomic.Bool)
- [x] Proper error handling in all new code
- [x] Logging added for debugging
- [x] Comments and documentation complete

### Backward Compatibility
- [x] Old config.json files still work
- [x] JSON sessions still loadable
- [x] Default behavior unchanged (opt-in features)
- [x] No breaking API changes

## 🧪 Feature Testing

### 1. JSONL Memory Store
- [ ] Create new session → saves as .jsonl
- [ ] Load existing JSON session → auto-migrates
- [ ] Crash during write → data not corrupted
- [ ] Multiple sessions → no file conflicts
- [ ] Large session → performance acceptable

**Test Commands:**
```bash
# Test new session
picoclaw agent
> Hello
> Exit

# Verify JSONL file created
ls ~/.picoclaw/sessions/*.jsonl

# Test migration
# (Copy old JSON session, restart, verify migration)
```

### 2. Vision/Image Support
- [ ] Send image to GPT-4V → describes correctly
- [ ] Send image to Claude 3 → describes correctly
- [ ] Large image (>10MB) → handled or rejected gracefully
- [ ] Invalid image → error message clear
- [ ] Multiple images → all processed

**Test Commands:**
```bash
picoclaw agent
> Describe this image: /path/to/image.jpg
```

### 3. Model Routing
- [ ] Simple query → uses cheap model
- [ ] Code query → uses expensive model
- [ ] Image query → uses expensive model
- [ ] Routing disabled → uses primary model
- [ ] Logging shows routing decisions

**Test Commands:**
```bash
# Enable routing in config.json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {"name": "cheap", "models": ["gpt-4o-mini"]},
      {"name": "expensive", "models": ["gpt-5.2"]}
    ]
  }
}

# Test simple query
picoclaw agent
> What is 2+2?  # Should use cheap model

# Test complex query
> Write a Python function to sort a list  # Should use expensive model
```

### 4. Extended Thinking
- [ ] Claude with thinking → shows reasoning
- [ ] OpenAI o1 → shows reasoning
- [ ] Reasoning channel configured → output sent there
- [ ] Thinking preserved in history
- [ ] No thinking → works normally

**Test Commands:**
```bash
# Use Claude 3.5 with extended thinking
picoclaw agent --model claude-sonnet-4.6
> Solve this complex problem: ...
```

### 5. Parallel Tool Execution
- [ ] Multiple tool calls → executed in parallel
- [ ] Results returned in correct order
- [ ] Error in one tool → others still execute
- [ ] Performance faster than sequential
- [ ] No race conditions

**Test Commands:**
```bash
picoclaw agent
> Search for "AI" and "ML" and "DL" simultaneously
```

### 6. Environment Configuration
- [ ] .env file loaded
- [ ] .env.local overrides .env
- [ ] Environment variables override config
- [ ] PICOCLAW_HOME/.env loaded
- [ ] Invalid .env → graceful error

**Test Commands:**
```bash
# Create .env
echo "PICOCLAW_ROUTING_ENABLED=true" > .env

# Verify loaded
picoclaw agent
# Check logs for routing enabled
```

### 7. Tool Enable/Disable
- [ ] Disable file tools → file operations fail
- [ ] Disable web tools → web search unavailable
- [ ] Enable all → all tools available
- [ ] Config change → takes effect on restart
- [ ] Default (all enabled) → works

**Test Commands:**
```bash
# Disable web tools in config
{
  "tools": {
    "web_tools_enabled": false
  }
}

# Test
picoclaw agent
> Search the web for "AI"  # Should fail or skip
```

### 8. Configurable Summarization
- [ ] 20 messages → triggers summarization
- [ ] Custom threshold → respected
- [ ] Token percentage → works correctly
- [ ] Summarization quality → acceptable
- [ ] Performance → not too slow

**Test Commands:**
```bash
# Set low threshold for testing
{
  "session": {
    "summarization_message_threshold": 5
  }
}

# Send 6 messages, verify summarization
```

### 9. PICOCLAW_HOME
- [ ] Set PICOCLAW_HOME → uses custom path
- [ ] Unset → uses ~/.picoclaw
- [ ] Config file in custom path → loaded
- [ ] Sessions in custom path → saved/loaded
- [ ] Multiple users → isolated data

**Test Commands:**
```bash
# Test custom home
export PICOCLAW_HOME=/tmp/picoclaw-test
picoclaw agent
> Hello

# Verify files in /tmp/picoclaw-test
ls /tmp/picoclaw-test
```

## 🔍 Integration Testing

### End-to-End Scenarios

#### Scenario 1: New User Setup
```bash
# Fresh installation
rm -rf ~/.picoclaw
picoclaw agent
> Hello, I'm a new user
> Exit

# Verify:
# - JSONL session created
# - Config generated
# - No errors
```

#### Scenario 2: Existing User Upgrade
```bash
# Simulate old installation
mkdir -p ~/.picoclaw/sessions
echo '{"messages":[]}' > ~/.picoclaw/sessions/test.json

# Run new version
picoclaw agent
# Load old session

# Verify:
# - JSON migrated to JSONL
# - Old data preserved
# - No data loss
```

#### Scenario 3: Multi-Modal Workflow
```bash
picoclaw agent
> Describe this image: /path/to/chart.png
> Now analyze the data you see
> Create a summary report

# Verify:
# - Image processed
# - Context maintained
# - Routing works
# - Session saved
```

#### Scenario 4: Cost Optimization
```bash
# Enable routing
picoclaw agent
> What's 2+2?  # Cheap model
> Write a complex algorithm  # Expensive model
> Hello  # Cheap model

# Verify:
# - Correct model selection
# - Cost savings
# - Quality maintained
```

## 📊 Performance Testing

### Metrics to Check
- [ ] Session load time < 100ms
- [ ] JSONL write time < 10ms
- [ ] Image encoding < 1s for 5MB image
- [ ] Parallel tools 2x faster than sequential
- [ ] Memory usage reasonable (<500MB)

### Load Testing
```bash
# Create 100 sessions
for i in {1..100}; do
  picoclaw agent --session "test-$i" <<< "Hello"
done

# Verify:
# - All sessions created
# - No file corruption
# - Performance acceptable
```

## 🔒 Security Testing

### Security Checks
- [ ] .env not committed to git
- [ ] API keys not in logs
- [ ] File permissions correct (600 for sensitive files)
- [ ] No SQL injection (N/A)
- [ ] No path traversal vulnerabilities

### Test Commands:
```bash
# Verify .env ignored
git status  # Should not show .env

# Check file permissions
ls -la ~/.picoclaw/auth.json  # Should be 600 or 644
```

## 📝 Documentation Testing

### Documentation Checks
- [ ] README accurate
- [ ] CHANGELOG complete
- [ ] .env.example comprehensive
- [ ] Integration guides clear
- [ ] API docs updated

### User Testing
- [ ] New user can follow quick start
- [ ] Configuration examples work
- [ ] Migration guide accurate
- [ ] Troubleshooting helpful

## 🚀 Deployment Checklist

### Pre-Deployment
- [ ] All tests passing
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] Version bumped
- [ ] Release notes written

### Deployment
- [ ] Binary built for all platforms
- [ ] Docker image updated
- [ ] Systemd service tested
- [ ] Backup procedures documented
- [ ] Rollback plan ready

### Post-Deployment
- [ ] Monitor logs for errors
- [ ] Check performance metrics
- [ ] Verify user feedback
- [ ] Update documentation if needed
- [ ] Plan next iteration

## ✅ Sign-Off

### Development Team
- [ ] All features implemented
- [ ] Code reviewed
- [ ] Tests passing
- [ ] Documentation complete

### QA Team
- [ ] All tests executed
- [ ] No critical bugs
- [ ] Performance acceptable
- [ ] Security verified

### Product Team
- [ ] Features meet requirements
- [ ] User experience good
- [ ] Documentation clear
- [ ] Ready for release

---

**Verification Date**: _____________  
**Verified By**: _____________  
**Status**: ⬜ Pending / ⬜ In Progress / ⬜ Complete  
**Notes**: _____________
