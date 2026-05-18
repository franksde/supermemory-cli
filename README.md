# `sm` — Blazing-fast Supermemory CLI for AI Agents

> Forked from [shaoyanji/supermemory-cli](https://github.com/shaoyanji/supermemory-cli)

Give your AI agents (Claude Code, Codex, Gemini CLI, etc.) persistent memory across sessions — powered by the [Supermemory](https://supermemory.ai) API. 

Designed for speed, token efficiency, and universal compatibility across all AI Agent environments.

## 🤖 For AI Agents (Self-Installation)

If you are an AI Agent reading this, you can install the CLI and its required skill file automatically by running:

```bash
# 1. Install the binary
curl -sSL https://raw.githubusercontent.com/franksde/supermemory-cli/main/install.sh | bash

# 2. Configure the API key (Ask the user if you don't have one)
sm config set-key <api-key>

# 3. Download the universal SKILL.md to your active workspace or skills directory
curl -sSL https://raw.githubusercontent.com/franksde/supermemory-cli/main/skills/SKILL.md -o supermemory_skill.md
```

## Why `sm`?

When integrating persistent memory into autonomous workflows, you need a tool that is:
1. **Extremely Fast**: Zero startup time, direct Go binary execution.
2. **Token Efficient**: JSON-only output tailored for LLMs. No bloated CLI UI.
3. **Explicit Control**: The agent explicitly decides what to store and when to recall, keeping context windows clean instead of auto-capturing everything.
4. **Universal**: Works with *any* agent that can execute shell commands.

## Quick Start (For Humans)

### Install

```bash
curl -sSL https://raw.githubusercontent.com/franksde/supermemory-cli/main/install.sh | bash
```

Or build from source:

```bash
git clone https://github.com/franksde/supermemory-cli.git
cd supermemory-cli
go build -o sm .
sudo mv sm /usr/local/bin/
```

### Configure

```bash
# Get your API key from https://console.supermemory.ai/keys
sm config set-key sm_xxxxx

# Set a default container tag (scopes your memories)
sm config set-tag my-project
```

Or use environment variables:

```bash
export SUPERMEMORY_API_KEY="sm_xxxxx"
export SUPERMEMORY_API_BASE="https://api.supermemory.ai"  # optional
```

### 30-Second Demo

```bash
sm add "This project uses Go with Cobra for CLI"
sm add "Frank prefers minimal comments in code"
sm search "project stack"
sm list
```

## CLI Reference

### Memory Operations (v4 API)

| Command | Description |
|---------|-------------|
| `sm add <text\|->` | Add a memory (from args or stdin) |
| `sm search <query>` | Search memories and documents |
| `sm list` | List recent memories |
| `sm forget <id\|content>` | Forget (soft delete) a memory |

### Document Operations (v3 API)

| Command | Description |
|---------|-------------|
| `sm doc add <file\|->` | Add a document |
| `sm doc batch <manifest.json>` | Batch add documents |
| `sm doc get <id>` | Get document by ID |
| `sm doc delete <id>` | Delete document by ID |

### Other Commands

| Command | Description |
|---------|-------------|
| `sm conv ingest <file.json>` | Ingest a conversation |
| `sm container get\|set\|merge\|delete` | Container tag management |
| `sm config set-key <key>` | Set API key |
| `sm config set-tag <tag>` | Set default container tag |
| `sm config show` | Show current configuration |

### Global Flags

- `--containerTag` — Override the default container tag for any command
- `--help` — Help for any command

## Configuration

Config file: `~/.config/sm/config.json`

```json
{
  "api_key": "sm_xxxxx",
  "base_url": "https://api.supermemory.ai",
  "default_container_tag": "my-project"
}
```

Priority: environment variables > config file > defaults.

Also reads legacy config from `~/.config/supermemory/config.json` for backward compatibility.

## Agent Integration

`sm` provides a universal `SKILL.md` file located at `skills/SKILL.md`.

This single skill file works for **any** agent (Claude Code, Codex, Gemini, etc.). It teaches the agent when to store memories, how to search for context at the start of a session, and how to manage the API.

Just provide the `skills/SKILL.md` to your agent, ensure `sm` is in the PATH, and the agent will handle the rest. All commands output JSON and use standard exit codes for machine consumption.

## License

MIT — see [LICENSE](LICENSE).