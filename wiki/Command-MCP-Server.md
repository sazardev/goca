# goca mcp-server

Start the **Model Context Protocol (MCP) server** that exposes all Goca commands as AI-callable tools. Once configured, AI assistants (GitHub Copilot, Claude Desktop, Cursor, Zed) can generate Clean Architecture scaffolding directly from chat.

## Syntax

```bash
goca mcp-server [flags]
```

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `--print-config` | `string` | Print a ready-to-copy client config snippet. One of: `vscode`, `claude`, `cursor`, `zed` |

## Quick Setup

```bash
# Print VS Code / GitHub Copilot config and paste into .vscode/mcp.json
goca mcp-server --print-config vscode

# Print Claude Desktop config
goca mcp-server --print-config claude
```

## Available MCP Tools

The server exposes **16 tools** — one per Goca command:

`goca_feature`, `goca_entity`, `goca_usecase`, `goca_repository`, `goca_handler`, `goca_di`, `goca_integrate`, `goca_interfaces`, `goca_messages`, `goca_mocks`, `goca_init`, `goca_doctor`, `goca_analyze`, `goca_ci`, `goca_middleware`, `goca_upgrade`

## See Also

- Full docs → https://sazardev.github.io/goca/commands/mcp-server
- AI integration guide → https://sazardev.github.io/goca/guide/mcp-integration
