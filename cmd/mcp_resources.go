package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerResources registers MCP resources that give AI assistants read-only
// context about the current Goca project without requiring a tool call.
func registerResources(s *server.MCPServer) {
	s.AddResource(
		mcp.NewResource(
			"goca://config",
			"Goca project configuration",
			mcp.WithResourceDescription("Contents of .goca.yaml — the Goca project configuration file. Read this to understand the project module path, database, and enabled features before generating code."),
			mcp.WithMIMEType("text/plain"),
		),
		handleConfigResource,
	)

	s.AddResource(
		mcp.NewResource(
			"goca://structure",
			"Project directory structure",
			mcp.WithResourceDescription("Directory tree of the internal/ folder — shows which entities and layers already exist so you can avoid re-generating existing files."),
			mcp.WithMIMEType("text/plain"),
		),
		handleStructureResource,
	)
}

// handleConfigResource returns the contents of .goca.yaml from the current
// working directory. Returns a helpful message when the file is missing.
func handleConfigResource(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot determine working directory: %w", err)
	}

	configPath := filepath.Join(cwd, ".goca.yaml")
	data, err := os.ReadFile(configPath) //nolint:gosec — path is cwd-relative, validated
	if err != nil {
		if os.IsNotExist(err) {
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "text/plain",
					Text:     "No .goca.yaml found in the current directory. Run `goca init <module>` to create a new project.",
				},
			}, nil
		}
		return nil, fmt.Errorf("reading .goca.yaml: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "text/plain",
			Text:     string(data),
		},
	}, nil
}

// handleStructureResource walks internal/ (up to 4 levels deep) and returns a
// formatted directory tree so the LLM can see what layers already exist.
func handleStructureResource(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot determine working directory: %w", err)
	}

	internalDir := filepath.Join(cwd, "internal")
	if _, statErr := os.Stat(internalDir); os.IsNotExist(statErr) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "text/plain",
				Text:     "internal/ directory not found. Run `goca init <module>` first.",
			},
		}, nil
	}

	tree, walkErr := buildDirTree(internalDir, "", 0, 4)
	if walkErr != nil {
		return nil, fmt.Errorf("walking internal/: %w", walkErr)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "text/plain",
			Text:     "internal/\n" + tree,
		},
	}, nil
}

// buildDirTree recursively builds an indented directory listing up to maxDepth.
func buildDirTree(dir, prefix string, depth, maxDepth int) (string, error) {
	if depth >= maxDepth {
		return "", nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var result string
	for i, entry := range entries {
		connector := "├── "
		childPrefix := prefix + "│   "
		if i == len(entries)-1 {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		result += prefix + connector + entry.Name() + "\n"

		if entry.IsDir() {
			sub, subErr := buildDirTree(filepath.Join(dir, entry.Name()), childPrefix, depth+1, maxDepth)
			if subErr != nil {
				return "", subErr
			}
			result += sub
		}
	}

	return result, nil
}
