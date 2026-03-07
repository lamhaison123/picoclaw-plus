# Queue Architecture Diagram

## Current Architecture (ManagerV2 - NO Queue)

```
User Message: "@architect @developer help me"
         |
         v
  HandleMentions()
         |
         +------------------+
         |                  |
         v                  v
   go execute()      go execute()
   @architect        @developer
         |                  |
         v                  v
    LLM API           LLM API
    (no limit)        (no limit)
         |                  |
         v                  v
    Response          Response
         |                  |
         +------------------+
                  |
                  v
         Check for mentions
                  |
         +--------+--------+
         |                 |
         v                 v
   go execute()      go execute()
   @tester           @manager
   (unlimited)       (unlimited)

❌ Problems:
- Unlimited concurrent goroutines
- No rate limiting
- No retry on failure
- No metrics
- Potential resource exhaustion
```

---

## Proposed Architecture (ManagerV2 + Queue)

```
User Message: "@architect @developer help me"
         |
         v
  HandleMentions()
         |
         +------------------+
         |                  |
         v                  v
    Enqueue()          Enqueue()
    @architect         @developer
         |                  |
         v                  v
   +---------+        +---------+
   | Queue   |        | Queue   |
   | Size:20 |        | Size:20 |
   | [====  ]|        | [====  ]|
   +---------+        +---------+
         |                  |
         v                  v
    Worker              Worker
    (rate limit)        (rate limit)
    (2s delay)          (2s delay)
         |                  |
         v                  v
    Execute             Execute
    (with retry)        (with retry)
         |                  |
         v                  v
    LLM API             LLM API
    (controlled)        (controlled)
         |                  |
         v                  v
    Response            Response
         |                  |
         +------------------+
                  |
                  v
         Check for mentions
                  |
         +--------+--------+
         |                 |
         v                 v
    Enqueue()         Enqueue()
    @tester           @manager
         |                 |
         v                 v
   +---------+       +---------+
   | Queue   |       | Queue   |
   | [==    ]|       | [==    ]|
   +---------+       +---------+
         |                 |
         v                 v
    Worker            Worker
    (rate limit)      (rate limit)

✅ Benefits:
- Controlled concurrency (1 worker per role)
- Rate limiting (2s between executions)
- Retry on failure (3 attempts with backoff)
- Metrics (queue length, dropped, processed)
- Resource protection
```

---

## Queue Manager Structure

```
QueueManager
    |
    +-- Queue[@architect]
    |       |
    |       +-- Worker (goroutine)
    |       +-- Rate Limiter (2s)
    |       +-- Retry Logic (3x)
    |       +-- Metrics
    |
    +-- Queue[@developer]
    |       |
    |       +-- Worker (goroutine)
    |       +-- Rate Limiter (2s)
    |       +-- Retry Logic (3x)
    |       +-- Metrics
    |
    +-- Queue[@tester]
    |       |
    |       +-- Worker (goroutine)
    |       +-- Rate Limiter (2s)
    |       +-- Retry Logic (3x)
    |       +-- Metrics
    |
    +-- Queue[@manager]
            |
            +-- Worker (goroutine)
            +-- Rate Limiter (2s)
            +-- Retry Logic (3x)
            +-- Metrics
```

---

## Execution Flow with Queue

```
1. User sends: "@architect design a system"
   
2. HandleMentions() receives mention
   
3. Create MentionRequest:
   {
     Role: "architect",
     Prompt: "design a system",
     SessionID: "abc123",
     ChatID: 12345,
     Depth: 0,
     ExecuteFunc: executeMentionRequest
   }
   
4. Enqueue to @architect queue:
   Queue[@architect].Enqueue(req)
   
5. Worker picks up request:
   - Check rate limit (last exec + 2s)
   - If too soon, sleep
   - Execute with retry:
     * Attempt 1: Execute
     * If fail, wait 1s
     * Attempt 2: Execute
     * If fail, wait 2s
     * Attempt 3: Execute
     * If fail, return error
   
6. Execute calls executeAgentAndCascade():
   - Check depth limit
   - Check dispatch tracker (idempotency)
   - Mark agent in cascade
   - Execute LLM
   - Send response
   - Unmark agent
   - Check for new mentions
   
7. If response has mentions:
   "@developer implement this"
   
8. Create new MentionRequest:
   {
     Role: "developer",
     Prompt: "implement this",
     Depth: 1,  // Incremented
     ...
   }
   
9. Enqueue to @developer queue
   
10. Repeat from step 5
```

---

## Rate Limiting Mechanism

```
Timeline for @architect queue:

T=0s    Request 1 arrives → Execute immediately
        lastExecTime = 0s

T=0.5s  Request 2 arrives → Queue
        Worker checks: now - lastExecTime = 0.5s < 2s
        Sleep for 1.5s

T=2s    Request 2 executes
        lastExecTime = 2s

T=3s    Request 3 arrives → Queue
        Worker checks: now - lastExecTime = 1s < 2s
        Sleep for 1s

T=4s    Request 3 executes
        lastExecTime = 4s

Result: Minimum 2 seconds between executions
```

---

## Retry Mechanism with Exponential Backoff

```
Request execution fails:

Attempt 1 (T=0s):
    Execute → FAIL
    Wait 1s (2^0 * 1s)

Attempt 2 (T=1s):
    Execute → FAIL
    Wait 2s (2^1 * 1s)

Attempt 3 (T=3s):
    Execute → FAIL
    Wait 4s (2^2 * 1s)

Attempt 4 (T=7s):
    Execute → SUCCESS or FINAL FAIL

Total retry time: 1s + 2s + 4s = 7s max
```

---

## Queue Overflow Handling

```
Queue Size: 5 (for example)

Current state: [Req1][Req2][Req3][Req4][Req5]
                ↑                           ↑
              front                        back

New request arrives:
    Try to enqueue → FAIL (queue full)
    
    Metrics.DroppedCount++
    
    Return error to user:
    "⚠️ Queue full for @architect, please try again later"
    
User sees immediate feedback instead of silent failure
```

---

## Metrics Dashboard (Proposed)

```
=== Queue Metrics ===

@architect:
  Queue Length:    3 / 20
  Processed:       145
  Dropped:         2
  Retry Count:     8
  Failure Count:   1
  Avg Wait Time:   1.2s

@developer:
  Queue Length:    1 / 20
  Processed:       98
  Dropped:         0
  Retry Count:     3
  Failure Count:   0
  Avg Wait Time:   0.8s

@tester:
  Queue Length:    0 / 20
  Processed:       67
  Dropped:         0
  Retry Count:     1
  Failure Count:   0
  Avg Wait Time:   0.5s

@manager:
  Queue Length:    2 / 20
  Processed:       34
  Dropped:         1
  Retry Count:     5
  Failure Count:   0
  Avg Wait Time:   1.5s
```

---

## Cascade with Queue

```
User: "@architect design a system"
         |
         v
    Enqueue @architect (depth=0)
         |
         v
    Execute @architect
         |
         v
    Response: "Sure! @developer please implement the backend"
         |
         v
    Detect mention: @developer
         |
         v
    Enqueue @developer (depth=1)
         |
         v
    Execute @developer
         |
         v
    Response: "@tester please test this"
         |
         v
    Detect mention: @tester
         |
         v
    Enqueue @tester (depth=2)
         |
         v
    Execute @tester
         |
         v
    Response: "Tests passed! @architect review please"
         |
         v
    Detect mention: @architect
         |
         v
    Check: @architect in cascade? NO (already unmarked)
    Check: depth < max? YES (2 < 3)
         |
         v
    Enqueue @architect (depth=3)
         |
         v
    Execute @architect
         |
         v
    Response: "@developer make changes"
         |
         v
    Detect mention: @developer
         |
         v
    Check: depth < max? NO (3 >= 3)
         |
         v
    STOP - Max depth reached

✅ Cascade works with queue
✅ Depth limit enforced
✅ Cycle detection works
✅ Rate limiting applied to all
```

---

## Comparison: Before vs After

### Before (No Queue)

```
Mention @architect
    ↓ (instant)
Spawn goroutine
    ↓ (instant)
Execute LLM
    ↓ (instant)
Mention @developer
    ↓ (instant)
Spawn goroutine
    ↓ (instant)
Execute LLM
    ↓ (instant)
...unlimited...

Problems:
- Goroutine explosion
- No rate control
- No retry
- No metrics
```

### After (With Queue)

```
Mention @architect
    ↓
Enqueue (instant)
    ↓
Queue [====  ]
    ↓
Worker (rate limited)
    ↓ (wait if needed)
Execute LLM (with retry)
    ↓
Mention @developer
    ↓
Enqueue (instant)
    ↓
Queue [====  ]
    ↓
Worker (rate limited)
    ↓ (wait if needed)
Execute LLM (with retry)

Benefits:
- Controlled concurrency
- Rate limiting
- Retry on failure
- Full metrics
- Resource protection
```

---

## Configuration Example

```json
{
  "channels": {
    "telegram": {
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "my-team",
        "max_context_length": 50,
        "queue_size": 20,
        "rate_limit_seconds": 2,
        "max_retries": 3,
        "retry_backoff_ms": 1000
      }
    }
  }
}
```

---

**Kết luận:** Queue architecture provides controlled, reliable, and observable mention execution.
