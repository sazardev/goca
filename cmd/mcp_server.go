package cmd

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// mcpServerCmd is the cobra command for "goca mcp-server".
// It is registered in root.go via rootCmd.AddCommand(mcpServerCmd).
var mcpServerCmd = newMCPServerCommand()

func newMCPServerCommand() *cobra.Command {
	var printConfig string

	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Start the Goca MCP server for AI assistant integration",
		Long: `Start a Model Context Protocol (MCP) server that exposes all Goca
code-generation commands as AI-callable tools.

Compatible clients: GitHub Copilot (VS Code), Claude Desktop, Cursor, Zed.

Quickstart — add to your client config and run:

  goca mcp-server

To print a ready-to-use client configuration snippet:

  goca mcp-server --print-config vscode
  goca mcp-server --print-config claude
  goca mcp-server --print-config cursor
  goca mcp-server --print-config zed`,
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if printConfig != "" {
				return printClientConfig(printConfig)
			}
			return runMCPServer()
		},
	}

	cmd.Flags().StringVar(
		&printConfig,
		"print-config",
		"",
		`Print client configuration snippet. One of: vscode, claude, cursor, zed`,
	)

	return cmd
}

// printClientConfig writes a ready-to-copy configuration snippet for the
// given MCP client to stdout and exits cleanly.
func printClientConfig(client string) error {
	configs := map[string]string{
		"vscode": `// .vscode/mcp.json
{
  "servers": {
    "goca": {
      "type": "stdio",
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}`,
		"claude": `// ~/Library/Application Support/Claude/claude_desktop_config.json
// (macOS) or %APPDATA%\Claude\claude_desktop_config.json (Windows)
{
  "mcpServers": {
    "goca": {
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}`,
		"cursor": `// .cursor/mcp.json  (project-level)
// or  ~/.cursor/mcp.json  (global)
{
  "mcpServers": {
    "goca": {
      "command": "goca",
      "args": ["mcp-server"]
    }
  }
}`,
		"zed": `// ~/.config/zed/settings.json  (add inside the root object)
{
  "context_servers": {
    "goca": {
      "command": {
        "path": "goca",
        "args": ["mcp-server"]
      }
    }
  }
}`,
	}

	snippet, ok := configs[strings.ToLower(client)]
	if !ok {
		return fmt.Errorf("unknown client %q — choose one of: vscode, claude, cursor, zed", client)
	}

	fmt.Println(snippet)
	return nil
}

// runMCPServer creates the MCP server, registers all tools and resources, and
// blocks serving on stdio until the client disconnects.
func runMCPServer() error {
	s := server.NewMCPServer(
		"goca",
		Version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
	)

	registerCoreTools(s)
	registerUtilTools(s)
	registerResources(s)

	return server.ServeStdio(s)
}
