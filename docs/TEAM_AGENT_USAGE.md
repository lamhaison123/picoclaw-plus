# Using Team Agents in PicoClaw

## Overview

PicoClaw supports multi-agent collaboration through teams. When teams are created, specialized team agents become available for delegation.

## How to Check Available Teams

Use the `picoclaw team list` command to see all active teams:

```bash
picoclaw team list
```

This will show:
- Team ID
- Team name
- Number of agents
- Status (active/dissolved)

## Team Agent Naming Convention

Team agents follow this naming pattern:
```
{team-id}-{role}
```

Examples:
- `dev-team-001-developer` - Developer role in dev-team-001
- `dev-team-001-reviewer` - Reviewer role in dev-team-001
- `research-team-001-researcher` - Researcher role in research-team-001

## Using Team Agents

### Option 1: Execute Task on Team (Recommended)

Use the team execute command to automatically route tasks to the right role:

```bash
picoclaw team execute dev-team-001 -t "Create a hello world function"
```

The team coordinator will:
1. Analyze the task
2. Determine the best role (developer, reviewer, tester, etc.)
3. Assign to the appropriate team agent
4. Execute and return results

### Option 2: Direct Agent Delegation (Advanced)

If you know the specific role needed, you can delegate directly to a team agent using the `spawn` tool:

```markdown
I need to delegate this coding task to the development team's developer agent.

spawn(
  agent_id="dev-team-001-developer",
  task="Implement a binary search function in Python with unit tests"
)
```

## Team Roles and Capabilities

Common team roles and their specializations:

| Role | Capabilities | Best For |
|------|-------------|----------|
| **architect** | System design, architecture planning | High-level design, technical specs |
| **developer** | Code generation, implementation | Writing code, implementing features |
| **tester** | Testing, validation, QA | Writing tests, finding bugs |
| **reviewer** | Code review, quality inspection | Reviewing code, suggesting improvements |
| **researcher** | Research, analysis, investigation | Gathering information, analysis |
| **manager** | Planning, coordination, delegation | Task planning, workflow management |

## Checking Team Status

To see detailed information about a team:

```bash
picoclaw team status dev-team-001
```

This shows:
- Team configuration
- Active agents and their roles
- Current status
- Coordinator information

## Creating New Teams

Create a team from a configuration file:

```bash
picoclaw team create templates/teams/development-team.json
```

See `templates/teams/` for example configurations.

## Best Practices

1. **Use team execute for automatic routing** - Let the coordinator choose the right role
2. **Check team list before delegating** - Ensure the team exists and is active
3. **Use correct agent IDs** - Follow the `{team-id}-{role}` format exactly
4. **Dissolve teams when done** - Clean up resources: `picoclaw team dissolve <team-id>`

## Example Workflow

```bash
# 1. List available teams
picoclaw team list

# 2. Check team details
picoclaw team status dev-team-001

# 3. Execute task (automatic role selection)
picoclaw team execute dev-team-001 -t "Implement user authentication"

# 4. Or delegate to specific role
# (Use spawn tool with agent_id="dev-team-001-developer")

# 5. Check team memory after completion
picoclaw team memory dev-team-001

# 6. Dissolve team when project is done
picoclaw team dissolve dev-team-001
```

## Troubleshooting

**"Team not found" error:**
- Run `picoclaw team list` to verify team ID
- Check if team was dissolved
- Ensure team was created successfully

**"Agent not found" error:**
- Verify agent ID format: `{team-id}-{role}`
- Check team status to see available roles
- Ensure team is active (not dissolved)

**Task execution fails:**
- Check team status for agent health
- Review team memory for previous errors
- Verify role has required capabilities

## Related Documentation

- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Complete guide to multi-agent collaboration
- [Team Tool Access](TEAM_TOOL_ACCESS.md) - Tool permissions and restrictions
- [Model Selection](MULTI_AGENT_MODEL_SELECTION.md) - Per-role model configuration
