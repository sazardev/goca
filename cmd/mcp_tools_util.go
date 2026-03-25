package cmd

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerUtilTools registers supporting infrastructure tools:
// goca_di, goca_integrate, goca_interfaces, goca_messages,
// goca_mocks, goca_init, goca_doctor, goca_upgrade.
func registerUtilTools(s *server.MCPServer) {
	s.AddTool(toolDI(), handleDI)
	s.AddTool(toolIntegrate(), handleIntegrate)
	s.AddTool(toolInterfaces(), handleInterfaces)
	s.AddTool(toolMessages(), handleMessages)
	s.AddTool(toolMocks(), handleMocks)
	s.AddTool(toolInit(), handleInit)
	s.AddTool(toolDoctor(), handleDoctor)
	s.AddTool(toolUpgrade(), handleUpgrade)
}

// ─── goca_di ─────────────────────────────────────────────────────────────────

func toolDI() mcp.Tool {
	return mcp.NewTool("goca_di",
		mcp.WithDescription("Generate a dependency injection (DI) container that wires all layers together."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity/module name to include in the DI container"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleDI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"di", name}
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_integrate ──────────────────────────────────────────────────────────

func toolIntegrate() mcp.Tool {
	return mcp.NewTool("goca_integrate",
		mcp.WithDescription("Wire all generated features into the main application entry point (router, DI, migrations)."),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview changes without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleIntegrate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"integrate"}
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_interfaces ─────────────────────────────────────────────────────────

func toolInterfaces() mcp.Tool {
	return mcp.NewTool("goca_interfaces",
		mcp.WithDescription("Generate Go interface contracts for an entity to support TDD and mocking."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. Product"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleInterfaces(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"interfaces", name}
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_messages ───────────────────────────────────────────────────────────

func toolMessages() mcp.Tool {
	return mcp.NewTool("goca_messages",
		mcp.WithDescription("Generate typed response/error message structs for an entity."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. Order"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleMessages(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"messages", name}
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_mocks ──────────────────────────────────────────────────────────────

func toolMocks() mcp.Tool {
	return mcp.NewTool("goca_mocks",
		mcp.WithDescription("Generate testify/mock stubs for repository and use-case interfaces of an entity."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name (optional — omit to generate mocks for all entities)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleMocks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"mocks"}
	if name := req.GetString("name", ""); name != "" {
		args = append(args, name)
	}
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_init ───────────────────────────────────────────────────────────────

func toolInit() mcp.Tool {
	return mcp.NewTool("goca_init",
		mcp.WithDescription("Initialise a new Go Clean Architecture project scaffold with the given module name."),
		mcp.WithString("name",
			mcp.Description("Go module name, e.g. github.com/acme/myapp"),
			mcp.Required(),
		),
		mcp.WithString("database",
			mcp.Description("Database driver to configure: postgres, mysql, sqlite, mongodb (optional)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleInit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"init", name}
	args = appendIfSet(args, req.GetString("database", ""), "--database")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_doctor ─────────────────────────────────────────────────────────────

func toolDoctor() mcp.Tool {
	return mcp.NewTool("goca_doctor",
		mcp.WithDescription("Diagnose the project for Clean Architecture issues and suggest fixes."),
		mcp.WithBoolean("fix",
			mcp.Description("Auto-apply safe fixes where possible"),
		),
	)
}

func handleDoctor(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"doctor"}
	args = appendIfTrue(args, req.GetBool("fix", false), "--fix")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_upgrade ────────────────────────────────────────────────────────────

func toolUpgrade() mcp.Tool {
	return mcp.NewTool("goca_upgrade",
		mcp.WithDescription("Check whether a newer version of Goca is available and optionally update the binary."),
		mcp.WithBoolean("update",
			mcp.Description("Download and install the latest version of Goca"),
		),
	)
}

func handleUpgrade(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"upgrade"}
	args = appendIfTrue(args, req.GetBool("update", false), "--update")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}
