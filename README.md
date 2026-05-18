# `sm` — Blazing-fast Supermemory CLI for AI Agents

> Forked from [shaoyanji/supermemory-cli](https://github.com/shaoyanji/supermemory-cli)

Give your AI agents (Claude Code, Codex, Gemini CLI, etc.) free, persistent memory across sessions — powered by the [Supermemory](https://supermemory.ai) API.

## Why?

The official Supermemory plugins require a **Pro plan**. The API free tier is generous enough for personal agent use. `sm` wraps that API into a single binary that any agent can call.

## Quick Start

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

### Claude Code

Copy `skills/claude-code/SKILL.md` into your project or global skills directory. The skill teaches Claude Code when and how to use `sm` for persistent memory.

### OpenAI Codex

Copy `skills/codex/SKILL.md` to `~/.codex/skills/supermemory/`. Codex auto-discovers skills from that directory.

### Any Agent

All commands output JSON and use exit code 1 on failure — designed for machine consumption. Add `sm` to your agent's allowed tools and teach it:

- `sm add "fact"` to store
- `sm search "query"` to recall
- `sm forget "outdated fact"` to remove

## License

MIT — see [LICENSE](LICENSE).