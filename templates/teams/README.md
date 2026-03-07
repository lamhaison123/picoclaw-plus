# Team Configuration Templates

Pre-built team configurations for different use cases. Copy and customize these templates for your collaborative chat needs.

## Quick Start

```bash
# 1. Copy a team template
cp templates/teams/collaborative-dev-team.json ~/.picoclaw/teams/dev-team.json

# 2. Update your config to use this team
nano ~/.picoclaw/config.json
# Set: "default_team_id": "dev-team"

# 3. Start gateway
picoclaw gateway

# 4. Use in Telegram
# @architect @developer Help me design this feature
```

## Available Templates

### Development Teams

#### `minimal-team.json`
**Roles:** 2 (architect, developer)  
**Pattern:** Parallel  
**Best for:** Simple projects, testing, learning

```bash
cp templates/teams/minimal-team.json ~/.picoclaw/teams/my-team.json
```

**Usage:**
```
@architect Design a REST API
@developer Implement the endpoints
```

#### `collaborative-dev-team.json`
**Roles:** 6 (architect, developer, tester, manager, devops, designer)  
**Pattern:** Parallel  
**Best for:** Full-featured development, complex projects

```bash
cp templates/teams/collaborative-dev-team.json ~/.picoclaw/teams/dev-team.json
```

**Usage:**
```
@architect @developer @tester Plan authentication system
@devops How should we deploy this?
@designer Create the login UI
@manager What's the timeline?
```

#### `development-team.json`
**Roles:** 3 (architect, developer, tester)  
**Pattern:** Sequential  
**Best for:** Structured workflows, step-by-step development

```bash
cp templates/teams/development-team.json ~/.picoclaw/teams/seq-team.json
```

**Usage:**
```
# Sequential execution: architect → developer → tester
picoclaw team execute seq-team "Build user authentication"
```

### Research Teams

#### `research-team-collab.json`
**Roles:** 3 (researcher, analyst, writer)  
**Pattern:** Parallel  
**Best for:** Information gathering, data analysis, report writing

```bash
cp templates/teams/research-team-collab.json ~/.picoclaw/teams/research.json
```

**Usage:**
```
@researcher Find information about AI trends
@analyst Analyze the market data
@writer Create a summary report
```

#### `research-team.json`
**Roles:** 3 (researcher, analyst, synthesizer)  
**Pattern:** Parallel  
**Best for:** Academic research, literature review

```bash
cp templates/teams/research-team.json ~/.picoclaw/teams/academic.json
```

### Support Teams

#### `support-team-collab.json`
**Roles:** 3 (support, technical, manager)  
**Pattern:** Hierarchical  
**Best for:** Customer support, help desk, escalation workflows

```bash
cp templates/teams/support-team-collab.json ~/.picoclaw/teams/support.json
```

**Usage:**
```
@support User can't login, help!
@technical Investigate the authentication service
@manager Should we escalate this?
```

### Creative Teams

#### `creative-team-collab.json`
**Roles:** 4 (writer, designer, strategist, reviewer)  
**Pattern:** Parallel  
**Best for:** Content creation, marketing, branding

```bash
cp templates/teams/creative-team-collab.json ~/.picoclaw/teams/creative.json
```

**Usage:**
```
@strategist Plan a product launch campaign
@writer Create the announcement post
@designer Design the visuals
@reviewer Review everything
```

### Analysis Teams

#### `analysis-team.json`
**Roles:** 3 (data_analyst, code_reviewer, security_auditor)  
**Pattern:** Hierarchical  
**Best for:** Code review, security audits, quality assurance

```bash
cp templates/teams/analysis-team.json ~/.picoclaw/teams/qa.json
```

## Team Patterns

### Parallel
All agents execute simultaneously. Fast but independent.

```json
{
  "pattern": "parallel"
}
```

**Use when:**
- Agents work independently
- Speed is important
- Tasks can be done in any order

**Example:** Multiple agents analyzing different aspects

### Sequential
Agents execute one after another. Structured workflow.

```json
{
  "pattern": "sequential"
}
```

**Use when:**
- Output of one agent feeds into next
- Strict order required
- Step-by-step process

**Example:** Design → Implement → Test

### Hierarchical
Coordinator delegates to specialized agents.

```json
{
  "pattern": "hierarchical"
}
```

**Use when:**
- Complex task breakdown needed
- Dynamic delegation
- Coordinator makes decisions

**Example:** Manager delegates to specialists

## Customizing Teams

### 1. Change Team ID

```json
{
  "team_id": "my-custom-team",  // Must match filename
  "name": "My Custom Team"
}
```

### 2. Add/Remove Roles

```json
{
  "roles": [
    {
      "name": "new-role",
      "description": "What this role does",
      "capabilities": ["skill1", "skill2"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    }
  ]
}
```

### 3. Change Models

```json
{
  "roles": [
    {
      "name": "architect",
      "model": "gpt-4-turbo-preview"  // Use GPT-4
    },
    {
      "name": "tester",
      "model": "claude-3-haiku-20240307"  // Faster/cheaper
    }
  ]
}
```

### 4. Restrict Tools

```json
{
  "roles": [
    {
      "name": "safe-role",
      "tools": ["file_read", "web_search"]  // Limited tools
    },
    {
      "name": "power-role",
      "tools": ["*"]  // All tools
    }
  ]
}
```

### 5. Adjust Timeouts

```json
{
  "settings": {
    "agent_timeout_seconds": 600,  // 10 minutes
    "max_delegation_depth": 5
  }
}
```

## Role Capabilities

Common capability categories:

### Development
- `design`, `architecture`, `planning`
- `coding`, `implementation`, `debugging`
- `testing`, `qa`, `validation`
- `deployment`, `ci_cd`, `monitoring`

### Research
- `web_search`, `data_analysis`, `information_gathering`
- `statistical_analysis`, `pattern_recognition`
- `documentation`, `report_writing`, `summarization`

### Support
- `customer_service`, `troubleshooting`
- `technical_troubleshooting`, `system_analysis`
- `escalation_handling`, `decision_making`

### Creative
- `copywriting`, `storytelling`, `content_creation`
- `visual_design`, `ui_ux`, `branding`
- `strategy`, `planning`, `campaign_design`

## Tool Permissions

Available tools:

### File Operations
- `file_read` - Read files
- `file_write` - Write files
- `file_list` - List directory contents

### Execution
- `exec_safe` - Safe command execution
- `exec` - Full command execution (use with caution)

### Web
- `web_search` - Search the web
- `web_fetch` - Fetch web pages

### Special
- `*` - All tools (use for trusted roles)

## Example Workflows

### Code Review Workflow

```json
{
  "team_id": "review-team",
  "pattern": "sequential",
  "roles": [
    {
      "name": "reviewer",
      "capabilities": ["code_review"],
      "tools": ["file_read"]
    },
    {
      "name": "tester",
      "capabilities": ["testing"],
      "tools": ["file_read", "exec_safe"]
    },
    {
      "name": "approver",
      "capabilities": ["approval"],
      "tools": ["file_read"]
    }
  ]
}
```

**Usage:**
```bash
picoclaw team execute review-team "Review PR #123"
```

### Content Creation Workflow

```json
{
  "team_id": "content-team",
  "pattern": "parallel",
  "roles": [
    {
      "name": "writer",
      "capabilities": ["writing"],
      "tools": ["file_write", "web_search"]
    },
    {
      "name": "editor",
      "capabilities": ["editing"],
      "tools": ["file_read", "file_write"]
    }
  ]
}
```

**Usage in Telegram:**
```
@writer Create a blog post about AI
@editor Review and improve it
```

### Research & Analysis Workflow

```json
{
  "team_id": "research-team",
  "pattern": "hierarchical",
  "roles": [
    {
      "name": "coordinator",
      "capabilities": ["coordination"],
      "tools": ["*"]
    },
    {
      "name": "researcher",
      "capabilities": ["research"],
      "tools": ["web_search", "file_write"]
    },
    {
      "name": "analyst",
      "capabilities": ["analysis"],
      "tools": ["file_read", "file_write"]
    }
  ]
}
```

## Best Practices

### 1. Clear Role Descriptions
```json
{
  "name": "architect",
  "description": "System architect - focuses on high-level design, architecture patterns, and system structure"
}
```

### 2. Specific Capabilities
```json
{
  "capabilities": [
    "system_design",
    "architecture_patterns",
    "scalability",
    "performance"
  ]
}
```

### 3. Appropriate Tool Access
```json
{
  "name": "reviewer",
  "tools": ["file_read"]  // Read-only for reviewers
}
```

### 4. Reasonable Timeouts
```json
{
  "settings": {
    "agent_timeout_seconds": 300  // 5 minutes for most tasks
  }
}
```

### 5. Choose Right Pattern
- **Parallel**: Independent tasks, speed matters
- **Sequential**: Dependencies between steps
- **Hierarchical**: Complex coordination needed

## Testing Your Team

### 1. Validate JSON
```bash
cat ~/.picoclaw/teams/my-team.json | jq .
```

### 2. List Teams
```bash
picoclaw team list
```

### 3. Check Team Status
```bash
picoclaw team status my-team
```

### 4. Test Execution
```bash
picoclaw team execute my-team "Simple test task"
```

### 5. Test in Telegram
```
@role1 @role2 Test message
```

## Troubleshooting

### Team not found
- Check filename matches `team_id`
- Verify file is in `~/.picoclaw/teams/`
- Validate JSON syntax

### Role not responding
- Check role name matches exactly
- Verify model is available
- Check tool permissions

### Timeout errors
- Increase `agent_timeout_seconds`
- Use faster models
- Simplify tasks

## See Also

- [Collaborative Chat Guide](../../docs/COLLABORATIVE_CHAT.md)
- [Multi-Agent Guide](../../docs/MULTI_AGENT_GUIDE.md)
- [Configuration Examples](../../config/README.md)
- [Tool Access Control](../../docs/TEAM_TOOL_ACCESS.md)
