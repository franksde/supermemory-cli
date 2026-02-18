# Supermemory CLI

A Go CLI for the Supermemory API. Supports document ingestion (single/batch), conversation archiving, search, and container tag management.

## Build

```bash
go build -o supermemory ./main.go
```

## Install

Copy the binary to a directory in your PATH, e.g.:

```bash
sudo cp supermemory /usr/local/bin/
```

Or to your user bin:

```bash
mkdir -p ~/.local/bin
cp supermemory ~/.local/bin/
```

## Configuration

Set your API key via environment:

```bash
export SUPERMEMORY_API_KEY="sm_..."
# Optional base URL (default: https://api.supermemory.ai)
export SUPERMEMORY_API_BASE="https://api.supermemory.ai"
```

Or create a config file at `~/.config/supermemory/config.json`:

```json
{
  "apiKey": "sm_...",
  "apiBase": "https://api.supermemory.ai"
}
```

## Commands

- `supermemory doc add <file|-` — add a document (read from file or stdin)
- `supermemory doc batch <manifest.json> [--dry-run]` — batch add (max 600 per request)
- `supermemory conv ingest <conversation.json>` — ingest a conversation
- `supermemory search docs <query> [--limit N] [--containerTag tag]` — search documents
- `supermemory container get <tag>` — get container settings
- `supermemory container set <tag> --settings key=value,...` — update settings
- `supermemory container merge <sourceTag> <targetTag>` — merge containers
- `supermemory container delete <tag>` — delete container

All commands output JSON.

## Example: Batch manifest

```json
[
  {
    "content": "Markdown or plain text content",
    "containerTag": "my-project",
    "customId": "doc-001",
    "metadata": { "source": "clipper" }
  }
]
```

## Notes

- Uses Cobra; `--help` for command flags.
- Batch chunking is automatic; respects the 600 document limit.
- Errors are printed to stderr; exit code 1 on failure.