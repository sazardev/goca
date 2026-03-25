---
layout: doc
title: AI Integration — MCP Server
titleTemplate: Goca Guide
description: Set up Goca as an MCP server so GitHub Copilot, Claude Desktop, Cursor, and Zed can generate Clean Architecture code directly from your AI chat.
---

# AI Integration (MCP Server)

Goca ships a built-in **Model Context Protocol (MCP) server** that lets any MCP-compatible AI assistant invoke Goca commands as tools. Instead of remembering every flag and typing commands in a terminal, you describe what you want to build and the AI generates the scaffolding for you.

::: tip Supported clients
GitHub Copilot (VS Code), Claude Desktop, Cursor, Zed — and any other client that implements the [MCP standard](https://modelcontextprotocol.io).
:::

## Prerequisites

1. **Goca installed** and available in your `PATH`:

   ```bash
   goca version   # e.g. Goca v1.18.7
   ```

2. **An MCP-compatible AI client** (see setup sections below).

3. A Goca-managed Go project (or start one with `goca init <module>`).

## Client Setup

### GitHub Copilot in VS Code

Create (or update) `.vscode/mcp.json` in your workspace root:

```json
{
  "servers": {
    "goca": {
      "type": "stdio",
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}
```

Then: open the Copilot Chat panel → click the **hammer icon** → you should see all `goca_*` tools listed.

You can also print this snippet directly:

```bash
goca mcp-server --print-config vscode
```

### Claude Desktop

Open `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows) and add:

```json
{
  "mcpServers": {
    "goca": {
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}
```

Restart Claude Desktop. A **Goca** section will appear in the tools drawer.

```bash
# Print snippet
goca mcp-server --print-config claude
```

### Cursor

Create `.cursor/mcp.json` in your project (project-level), or `~/.cursor/mcp.json` (global):

```json
{
  "mcpServers": {
    "goca": {
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}
```

In Cursor: **Settings → MCP** → confirm the server appears as **Connected**.

```bash
goca mcp-server --print-config cursor
```

### Zed

In `~/.config/zed/settings.json`, add inside the root object:

```json
{
  "context_servers": {
    "goca": {
      "command": {
        "path": "goca",
        "args": ["mcp-server"]
      }
    }
  }
}
```

```bash
goca mcp-server --print-config zed
```

## Example Conversations

Once Goca is connected, you can talk to the AI naturally:

---

> **"Create a complete Product feature for an e-commerce app. Fields: Name (string), Price (float64), Stock (int). Use postgres."**

The AI calls `goca_feature` with the right parameters and generates all four layers.

---

> **"Show me what files a User entity with email, age, and soft-delete would create — don't write to disk yet."**

The AI calls `goca_entity` with `dry_run: true` and shows you a preview.

---

> **"What entities already exist in this project?"**

The AI reads the `goca://structure` resource and describes the current `internal/` tree.

---

> **"Generate testify mocks for all my interfaces."**

The AI calls `goca_mocks` with no name parameter, targeting every interface.

---

## Available Tools Reference

| Tool | What it does |
|---|---|
| `goca_feature` | Full feature: entity + use-case + repository + handler |
| `goca_entity` | Domain entity |
| `goca_usecase` | Application service / use-case |
| `goca_repository` | Repository interface + implementation |
| `goca_handler` | HTTP handler |
| `goca_di` | Dependency injection container |
| `goca_integrate` | Wire everything into the app entry point |
| `goca_interfaces` | Interface contracts for TDD |
| `goca_messages` | Typed response/error message structs |
| `goca_mocks` | testify/mock stubs |
| `goca_init` | Scaffold a new project |
| `goca_doctor` | Diagnose architecture issues |
| `goca_upgrade` | Check for / install updates |

Full parameter documentation → [goca mcp-server command reference](/commands/mcp-server)

## Project Context Resources

Two read-only **MCP resources** give the AI automatic project context:

| Resource | Contents |
|---|---|
| `goca://config` | Your `.goca.yaml` — module path, database, features |
| `goca://structure` | Live `internal/` directory tree |

Clients that support resource reads (Claude Desktop, Cursor) use these automatically to avoid regenerating files that already exist.

## How It Works

```
You ─────────────────────────────────────────────────────────────────────┐
                                                                          │
AI Client (Copilot / Claude / Cursor)                                    │
  1. Receives your request                                                │
  2. Decides which goca_* tool to call                                   │
  3. Sends JSON-RPC CallTool over stdio ──────────────────────────────┐  │
                                                                       │  │
goca mcp-server (subprocess in your project dir)                      │  │
  4. Receives the CallTool request                                     │  │
  5. Runs: goca feature Product --fields ... --database postgres       │  │
  6. Returns captured stdout back as tool result ───────────────────┐  │  │
                                                                    │  │  │
AI Client                                                           │  │  │
  7. Formats the result and shows it to you ◄───────────────────────┘  │  │
                                                                        │  │
Generated files (written in your project directory) ◄──────────────────┘  │
                                                                           │
You see the result in chat ◄───────────────────────────────────────────────┘
```

## Security Considerations

- The MCP server runs **locally** with the same permissions as your user account.
- No network connections are made to external servers.
- All input strings are validated against `^[A-Za-z][A-Za-z0-9]*$` before use in file paths or templates.
- Arguments are passed as separate elements to `exec.Command` — no shell interpolation or injection risk.

## Troubleshooting

**The server doesn't appear in my client**

- Confirm `goca` is in the same `PATH` your client uses. Try: `which goca` and `goca version`.
- On macOS with Homebrew, add `/opt/homebrew/bin` to the client's environment if needed.

**"context deadline exceeded" or slow responses**

Large projects with many files in `internal/` may slow down the `goca://structure` resource. This is a read-only walk and does not affect code generation tools.

**Files aren't being written**

Check if `dry_run` was set to `true`. Ask the AI: *"Run that again but actually write the files."*

**I want to undo generated files**

Use `--backup` in the tool call, or restore from git: `git checkout -- .`
