---
name: supermemory
description: "Use when you need to store, search, or manage memories across sessions. Provides persistent memory for AI coding agents via the Supermemory API."
---

# Supermemory — Persistent Memory for AI Agents

Store facts, preferences, and project knowledge that persist across sessions.

## When to Use

- **Store**: When the user shares preferences, project facts, or decisions worth remembering
- **Search**: When you need context from previous sessions (e.g., "what stack does this project use?")
- **Forget**: When information becomes outdated or incorrect

## Commands

### Store a memory
```bash
sm add "Frank prefers TypeScript over JavaScript"
sm add "This project uses Next.js 15 with App Router" --containerTag project-name
echo "deployment uses Cloudflare Workers" | sm add -
```

### Search memories
```bash
sm search "user preferences"
sm search "database schema" --containerTag project-name
sm search "authentication" --v4 --limit 5
```

### List recent memories
```bash
sm list --containerTag project-name
```

### Forget outdated info
```bash
sm forget "old incorrect fact" --containerTag project-name
sm forget mem_abc123 --containerTag project-name
```

### Document operations (for larger content)
```bash
sm doc add README.md --containerTag project-name
sm doc get <document-id>
sm doc delete <document-id>
```

## Common Patterns

### Session Start — Recall Context
At the start of a session, search for relevant project context:
```bash
sm search "project architecture and conventions" --containerTag $(basename $(pwd))
```

### Session End — Save Key Learnings
After making important decisions or discoveries:
```bash
sm add "Decided to use Prisma ORM for database layer" --containerTag project-name
sm add "Auth flow: JWT tokens with 24h expiry, refresh via /api/auth/refresh"
```

### User Preferences
When the user expresses a preference:
```bash
sm add "User prefers dark mode in all UIs"
sm add "User wants minimal comments, self-documenting code"
```

## Setup

If `sm` is not installed or configured:
```bash
curl -sSL https://raw.githubusercontent.com/franksde/supermemory-cli/main/install.sh | bash
sm config set-key <api-key>
sm config set-tag <default-container-tag>
```

## Notes

- All commands output JSON (agent-friendly)
- `--containerTag` scopes memories; set a default with `sm config set-tag`
- The `sm add` command uses the v4 memories API (immediate, no ingestion delay)
- The `sm search` command defaults to v3 document search; use `--v4` for memory search
