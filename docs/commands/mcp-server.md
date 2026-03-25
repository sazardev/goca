---
layout: doc
title: goca mcp-server
titleTemplate: Goca CLI Reference
description: Start the Goca MCP server so AI assistants (GitHub Copilot, Claude Desktop, Cursor, Zed) can invoke code-generation commands directly from chat.
---

# `goca mcp-server`

Start the **Model Context Protocol (MCP) server** that exposes all Goca commands as AI-callable tools. Once configured, AI assistants can generate Clean Architecture scaffolding without you ever leaving the chat window.

## Synopsis

```bash
goca mcp-server [flags]
```

## Flags

| Flag | Type | Description |
|---|---|---|
| `--print-config` | `string` | Print a ready-to-copy client configuration snippet. One of: `vscode`, `claude`, `cursor`, `zed` |

## Quick Setup

### 1. Confirm `goca` is in your PATH

```bash
which goca   # should print a path
goca version # should print the current version
```

### 2. Print the client config and paste it

```bash
# GitHub Copilot in VS Code
goca mcp-server --print-config vscode

# Claude Desktop
goca mcp-server --print-config claude

# Cursor
goca mcp-server --print-config cursor

# Zed
goca mcp-server --print-config zed
```

### 3. Restart your AI client and look for the Goca tools

In VS Code: open the chat panel → the hammer icon will list all `goca_*` tools.

See the full setup guide → [AI Integration](/guide/mcp-integration)

## Available Tools

The MCP server exposes **13 tools** — one per Goca command:

### Code Generation

| Tool | Equivalent CLI | Description |
|---|---|---|
| `goca_feature` | `goca feature <Name>` | Generate a complete feature (entity + use-case + repository + handler) |
| `goca_entity` | `goca entity <Name>` | Generate a domain entity |
| `goca_usecase` | `goca usecase <Name>` | Generate a use-case (application service) |
| `goca_repository` | `goca repository <Name>` | Generate a repository interface + implementation |
| `goca_handler` | `goca handler <Name>` | Generate an HTTP handler |

### Infrastructure & Support

| Tool | Equivalent CLI | Description |
|---|---|---|
| `goca_di` | `goca di <Name>` | Generate a dependency injection container |
| `goca_integrate` | `goca integrate` | Wire all features into the application entry point |
| `goca_interfaces` | `goca interfaces <Name>` | Generate interface contracts for TDD |
| `goca_messages` | `goca messages <Name>` | Generate typed response/error message structs |
| `goca_mocks` | `goca mocks [Name]` | Generate testify mock stubs |
| `goca_init` | `goca init <Name>` | Initialize a new project scaffold |
| `goca_doctor` | `goca doctor` | Diagnose the project for Architecture issues |
| `goca_upgrade` | `goca upgrade` | Check for and install Goca updates |

## MCP Resources

The server also exposes two read-only **resources** that give AI assistants project context:

| Resource URI | Description |
|---|---|
| `goca://config` | Contents of `.goca.yaml` — module name, database, enabled features |
| `goca://structure` | Directory tree of `internal/` — shows which layers already exist |

AI clients that support MCP resources (like Claude Desktop and Cursor) will automatically use these to avoid regenerating files that already exist.

## Example Chat Interactions

Once configured, you can ask naturally:

> **"Create a Product feature with fields Name:string, Price:float64, Stock:int using postgres"**

The AI will call `goca_feature` with the appropriate parameters.

> **"Show me what files would be created for a User entity with email and timestamps"**

The AI will call `goca_entity` with `dry_run: true` and show you the preview.

> **"What entities does this project already have?"**

The AI will read the `goca://structure` resource and answer from the live directory tree.

## How It Works

`goca mcp-server` starts a [Model Context Protocol](https://modelcontextprotocol.io) server over **stdio**. The AI client launches `goca mcp-server` as a subprocess and communicates via JSON-RPC messages on stdin/stdout. When you ask the AI to generate code, it calls the appropriate MCP tool, which in turn runs the corresponding `goca` subcommand in your project directory.

```
AI Client ──JSON-RPC──▶ goca mcp-server ──subprocess──▶ goca feature Product ...
                                                                │
                                                                ▼
                                                    internal/domain/product.go
                                                    internal/usecase/product_service.go
                                                    internal/repository/product_repository.go
                                                    internal/handler/http/product_handler.go
```

## Security

- The server runs **locally** with the same filesystem permissions as your terminal session.
- All user-supplied strings (entity names, field types) are validated by `CommandValidator` before use.
- Arguments are passed as separate elements to `exec.Command` — no shell interpolation.
- The server never opens network ports; it communicates exclusively over stdio.

## Troubleshooting

**"goca not found" error in client logs**

Make sure the full path to `goca` is in the `PATH` used by your client. On macOS/Linux, add `export PATH="$PATH:/usr/local/bin"` (or wherever goca is installed) to your shell profile and restart the client.

**Tools don't appear in the client**

1. Check the client's MCP server log for startup errors.
2. Run `goca mcp-server` manually in a terminal — it should block waiting for input with no output.
3. Press Ctrl-C to stop.

**"No .goca.yaml found" when reading resources**

The MCP server resolves paths relative to `cwd` at startup. Configure your client to set `cwd` to the project root:

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

Most clients default `cwd` to the workspace/project folder automatically.
