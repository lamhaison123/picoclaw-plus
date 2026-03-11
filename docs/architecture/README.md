# Architecture Documentation

This directory contains comprehensive documentation about PicoClaw's system architecture and design.

## 📑 Contents

### Core Architecture
- **[ARCHITECTURE_OVERVIEW.md](./ARCHITECTURE_OVERVIEW.md)** - Complete system architecture with diagrams
  - Message Bus design
  - Agent Loop execution
  - Team coordination
  - Memory system
  - Design patterns

- **[COMPONENT_DETAILS.md](./COMPONENT_DETAILS.md)** - Detailed component specifications
  - Struct definitions
  - Key methods
  - Code examples
  - Component interactions

- **[DATA_FLOW.md](./DATA_FLOW.md)** - Step-by-step flow diagrams
  - Message processing (21 steps)
  - Team delegation flow
  - Mention cascading flow
  - ASCII diagrams

### Codebase Structure
- **[CODEBASE_OVERVIEW.md](./CODEBASE_OVERVIEW.md)** - Code organization and structure
- **[REPOSITORY_OVERVIEW.md](./REPOSITORY_OVERVIEW.md)** - Repository layout (English)
- **[REPOSITORY_OVERVIEW.vi.md](./REPOSITORY_OVERVIEW.vi.md)** - Repository layout (Vietnamese)

### Legacy
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - Original architecture document (may be outdated)

## 🎯 Reading Guide

### For New Developers
1. Start with **ARCHITECTURE_OVERVIEW.md** to understand the big picture
2. Read **DATA_FLOW.md** to see how messages flow through the system
3. Dive into **COMPONENT_DETAILS.md** for implementation details
4. Reference **CODEBASE_OVERVIEW.md** when navigating the code

### For System Designers
1. Review **ARCHITECTURE_OVERVIEW.md** for design patterns
2. Study **COMPONENT_DETAILS.md** for component interfaces
3. Analyze **DATA_FLOW.md** for integration points

### For Contributors
1. Read **REPOSITORY_OVERVIEW.md** to understand project structure
2. Check **CODEBASE_OVERVIEW.md** for coding conventions
3. Reference **COMPONENT_DETAILS.md** when modifying components

## 🔑 Key Concepts

### Message Bus
Central hub for all communication between components. Uses buffered channels with context cancellation support.

### Agent Loop
Core execution engine that processes messages, manages sessions, executes tools, and coordinates with LLM providers.

### Team Manager
Orchestrates multi-agent collaboration with role-based task delegation and shared context management.

### Memory System
Dual-layer memory with session-based history (JSON files) and semantic search (vector database).

### Collaborative Chat
Multi-agent conversation system with @mention detection, cascade prevention, and context compaction.

## 📊 Architecture Diagrams

All architecture documents include ASCII diagrams for:
- System overview
- Component relationships
- Data flow sequences
- State transitions

## 🔗 Related Documentation

- [Developer Guide](../development/DEVELOPER_GUIDE.md) - Development workflow
- [API Reference](../reference/API_REFERENCE.md) - API documentation
- [Troubleshooting](../development/troubleshooting.md) - Common issues

## 📝 Contributing

When updating architecture documentation:
1. Keep diagrams in sync with code
2. Update all related documents
3. Add examples for complex concepts
4. Maintain both English and Vietnamese versions where applicable

---

**Last Updated**: 2026-03-09
