---
name: supermemory
description: "Use when you need to store, search, or manage memories across sessions. Provides persistent memory via the Supermemory API."
---

# Supermemory — Persistent Memory for Codex

Store and retrieve facts, preferences, and project knowledge across sessions.

## When to Use

- User shares preferences or project decisions worth remembering
- You need context from previous sessions
- Information becomes outdated and should be forgotten

## Quick Reference

```bash
# Store a memory
sm add "This project uses Vitest for unit tests"

# Search memories
sm search "testing framework"

# List recent memories
sm list --containerTag project-name

# Forget outdated info
sm forget "old fact" --containerTag project-name

# Add a document
sm doc add README.md --containerTag project-name
```

## Setup

```bash
curl -sSL https://raw.githubusercontent.com/franksde/supermemory-cli/main/install.sh | bash
sm config set-key <api-key>
sm config set-tag <default-container-tag>
```

## Notes

- All output is JSON
- Use `--containerTag` to scope memories per project
- `sm add` uses v4 API (instant), `sm search` uses v3 by default (use `--v4` for memory search)
