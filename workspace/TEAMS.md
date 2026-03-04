# Team Agents

## What are Team Agents?

Team agents are specialized AI agents that work together in coordinated teams. Each team has multiple roles (developer, reviewer, tester, etc.) with specific capabilities.

## How to Use Teams

### Check Available Teams

To see what teams are available, use:
```bash
picoclaw team list
```

### Execute Tasks on Teams

When you need specialized help, delegate to a team:

```bash
picoclaw team execute <team-id> -t "your task description"
```

The team coordinator will automatically:
1. Analyze your task
2. Select the best role for the job
3. Execute with the specialized agent
4. Return results

### Example

```bash
# List teams
picoclaw team list

# Execute coding task
picoclaw team execute dev-team-001 -t "Create a REST API endpoint for user login"

# Check team status
picoclaw team status dev-team-001
```

## Team Agent IDs

Team agents have IDs in this format: `{team-id}-{role}`

Examples:
- `dev-team-001-developer` - Writes code
- `dev-team-001-reviewer` - Reviews code quality
- `dev-team-001-tester` - Writes and runs tests

## When to Use Teams

Use teams when you need:
- **Specialized expertise** - Different roles for different tasks
- **Quality assurance** - Code review and testing
- **Complex workflows** - Multi-step processes (design → implement → test → review)
- **Collaboration** - Multiple agents working together

## Common Team Roles

| Role | What They Do |
|------|--------------|
| **architect** | Design system architecture, create technical specs |
| **developer** | Write code, implement features |
| **tester** | Write tests, validate functionality |
| **reviewer** | Review code quality, suggest improvements |
| **researcher** | Research solutions, analyze problems |
| **manager** | Plan tasks, coordinate workflow |

## Tips

1. **Let the coordinator choose** - Use `team execute` instead of picking roles manually
2. **Check team list first** - Make sure the team exists before delegating
3. **Use for complex tasks** - Teams are best for multi-step or specialized work
4. **Review team memory** - Check `picoclaw team memory <team-id>` to see what the team accomplished

## Learn More

See `docs/MULTI_AGENT_GUIDE.md` for complete documentation on multi-agent collaboration.
