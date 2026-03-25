package cmd

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerCoreTools registers the five primary code-generation tools:
// goca_feature, goca_entity, goca_usecase, goca_repository, goca_handler.
func registerCoreTools(s *server.MCPServer) {
	s.AddTool(toolFeature(), handleFeature)
	s.AddTool(toolEntity(), handleEntity)
	s.AddTool(toolUsecase(), handleUsecase)
	s.AddTool(toolRepository(), handleRepository)
	s.AddTool(toolHTTPHandler(), handleHandler)
}

// ─── goca_feature ────────────────────────────────────────────────────────────

func toolFeature() mcp.Tool {
	return mcp.NewTool("goca_feature",
		mcp.WithDescription("Generate a complete Clean Architecture feature (entity + usecase + repository + handler) for a Go project managed by Goca."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. Product, OrderItem"),
			mcp.Required(),
		),
		mcp.WithString("fields",
			mcp.Description(`Comma-separated field definitions, e.g. "Name:string,Price:float64,Stock:int"`),
			mcp.Required(),
		),
		mcp.WithString("database",
			mcp.Description("Database driver: postgres, mysql, sqlite, mongodb (optional)"),
		),
		mcp.WithBoolean("validation",
			mcp.Description("Add input validation to the use-case layer"),
		),
		mcp.WithBoolean("business_rules",
			mcp.Description("Generate business-rule stubs in the entity"),
		),
		mcp.WithBoolean("handlers",
			mcp.Description("Generate HTTP handler (default: true when feature is generated)"),
		),
		mcp.WithBoolean("integration_tests",
			mcp.Description("Generate integration test files"),
		),
		mcp.WithBoolean("mocks",
			mcp.Description("Generate testify mock stubs for repository and use-case interfaces"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview files that would be created without writing to disk (recommended before committing)"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files without prompting"),
		),
	)
}

func handleFeature(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"feature", name}
	args = appendIfSet(args, req.GetString("fields", ""), "--fields")
	args = appendIfSet(args, req.GetString("database", ""), "--database")
	args = appendIfTrue(args, req.GetBool("validation", false), "--validation")
	args = appendIfTrue(args, req.GetBool("business_rules", false), "--business-rules")
	args = appendIfTrue(args, req.GetBool("integration_tests", false), "--integration-tests")
	args = appendIfTrue(args, req.GetBool("mocks", false), "--mocks")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_entity ─────────────────────────────────────────────────────────────

func toolEntity() mcp.Tool {
	return mcp.NewTool("goca_entity",
		mcp.WithDescription("Generate a domain entity for a Go Clean Architecture project."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. User"),
			mcp.Required(),
		),
		mcp.WithString("fields",
			mcp.Description(`Comma-separated field definitions, e.g. "Email:string,Age:int"`),
		),
		mcp.WithBoolean("validation",
			mcp.Description("Add validation methods to the entity"),
		),
		mcp.WithBoolean("business_rules",
			mcp.Description("Generate business-rule stubs"),
		),
		mcp.WithBoolean("timestamps",
			mcp.Description("Add CreatedAt / UpdatedAt timestamp fields"),
		),
		mcp.WithBoolean("soft_delete",
			mcp.Description("Add DeletedAt soft-delete field"),
		),
		mcp.WithBoolean("tests",
			mcp.Description("Generate unit tests for the entity"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleEntity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"entity", name}
	args = appendIfSet(args, req.GetString("fields", ""), "--fields")
	args = appendIfTrue(args, req.GetBool("validation", false), "--validation")
	args = appendIfTrue(args, req.GetBool("business_rules", false), "--business-rules")
	args = appendIfTrue(args, req.GetBool("timestamps", false), "--timestamps")
	args = appendIfTrue(args, req.GetBool("soft_delete", false), "--soft-delete")
	args = appendIfTrue(args, req.GetBool("tests", false), "--tests")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_usecase ────────────────────────────────────────────────────────────

func toolUsecase() mcp.Tool {
	return mcp.NewTool("goca_usecase",
		mcp.WithDescription("Generate a use-case (application service) for a domain entity."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name the use-case operates on"),
			mcp.Required(),
		),
		mcp.WithString("fields",
			mcp.Description(`Comma-separated DTO field definitions, e.g. "Title:string,Done:bool"`),
		),
		mcp.WithBoolean("validation",
			mcp.Description("Add input validation in the use-case"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleUsecase(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"usecase", name}
	args = appendIfSet(args, req.GetString("fields", ""), "--fields")
	args = appendIfTrue(args, req.GetBool("validation", false), "--validation")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_repository ─────────────────────────────────────────────────────────

func toolRepository() mcp.Tool {
	return mcp.NewTool("goca_repository",
		mcp.WithDescription("Generate a repository interface and its concrete implementation for an entity."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. Order"),
			mcp.Required(),
		),
		mcp.WithString("database",
			mcp.Description("Database driver: postgres, mysql, sqlite, mongodb (optional)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleRepository(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"repository", name}
	args = appendIfSet(args, req.GetString("database", ""), "--database")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_handler ────────────────────────────────────────────────────────────

func toolHTTPHandler() mcp.Tool {
	return mcp.NewTool("goca_handler",
		mcp.WithDescription("Generate an HTTP handler (delivery layer) for a domain entity."),
		mcp.WithString("name",
			mcp.Description("PascalCase entity name, e.g. Product"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description("Handler type: http, grpc, graphql (default: http)"),
		),
		mcp.WithBoolean("validation",
			mcp.Description("Add request validation in the handler"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview without writing files"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing files"),
		),
	)
}

func handleHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcpErr(err), nil
	}

	args := []string{"handler", name}
	args = appendIfSet(args, req.GetString("type", ""), "--type")
	args = appendIfTrue(args, req.GetBool("validation", false), "--validation")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}
