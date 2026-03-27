package cmd

import (
	"context"
	"fmt"

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
	s.AddTool(toolCI(), handleCI)
	s.AddTool(toolMiddleware(), handleMiddleware)
	s.AddTool(toolAnalyze(), handleAnalyze)
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

// ─── goca_ci ─────────────────────────────────────────────────────────────────

func toolCI() mcp.Tool {
	return mcp.NewTool("goca_ci",
		mcp.WithDescription("Generate CI/CD pipeline configuration (GitHub Actions workflows for test, build, and deploy)."),
		mcp.WithString("provider",
			mcp.Description("CI provider: github-actions (default)"),
		),
		mcp.WithBoolean("with_docker",
			mcp.Description("Include Docker build step in the build workflow"),
		),
		mcp.WithBoolean("with_deploy",
			mcp.Description("Generate a tag-triggered deploy workflow"),
		),
		mcp.WithString("go_version",
			mcp.Description("Go version for the CI matrix (reads from go.mod if omitted)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview files without writing to disk"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing workflow files"),
		),
	)
}

func handleCI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"ci"}
	args = appendIfSet(args, req.GetString("provider", ""), "--provider")
	args = appendIfTrue(args, req.GetBool("with_docker", false), "--with-docker")
	args = appendIfTrue(args, req.GetBool("with_deploy", false), "--with-deploy")
	args = appendIfSet(args, req.GetString("go_version", ""), "--go-version")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_middleware ─────────────────────────────────────────────────────────

func toolMiddleware() mcp.Tool {
	return mcp.NewTool("goca_middleware",
		mcp.WithDescription("Generate a dedicated internal/middleware/ package with composable HTTP middleware (CORS, logging, auth, rate-limit, recovery, request-id, timeout)."),
		mcp.WithString("name",
			mcp.Description("Middleware package name (e.g. the project or feature name)"),
			mcp.Required(),
		),
		mcp.WithString("types",
			mcp.Description("Comma-separated middleware types: cors,logging,auth,rate-limit,recovery,request-id,timeout (default: cors,logging,recovery)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Preview files without writing to disk"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Overwrite existing middleware files"),
		),
	)
}

func handleMiddleware(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.GetString("name", "")
	if name == "" {
		return mcpErr(fmt.Errorf("name is required")), nil
	}
	args := []string{"middleware", name}
	args = appendIfSet(args, req.GetString("types", ""), "--types")
	args = appendIfTrue(args, req.GetBool("dry_run", false), "--dry-run")
	args = appendIfTrue(args, req.GetBool("force", false), "--force")

	out, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(out), nil
}

// ─── goca_analyze ─────────────────────────────────────────────────────────────

func toolAnalyze() mcp.Tool {
	return mcp.NewTool("goca_analyze",
		mcp.WithDescription("Deep self-analysis of the generated project: architecture layer boundaries, security (OWASP), code quality, naming standards, test coverage, and dependency hygiene."),
		mcp.WithBoolean("arch",
			mcp.Description("Only run architecture checks"),
		),
		mcp.WithBoolean("quality",
			mcp.Description("Only run code quality checks"),
		),
		mcp.WithBoolean("security",
			mcp.Description("Only run security checks (OWASP A03, hardcoded secrets, unsafe)"),
		),
		mcp.WithBoolean("standards",
			mcp.Description("Only run Go standards checks (naming, context propagation, go.mod)"),
		),
		mcp.WithBoolean("tests",
			mcp.Description("Only run test coverage and pattern checks"),
		),
		mcp.WithBoolean("deps",
			mcp.Description("Only run dependency hygiene checks"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: text (default) or json"),
		),
	)
}

func handleAnalyze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"analyze"}
	args = appendIfTrue(args, req.GetBool("arch", false), "--arch")
	args = appendIfTrue(args, req.GetBool("quality", false), "--quality")
	args = appendIfTrue(args, req.GetBool("security", false), "--security")
	args = appendIfTrue(args, req.GetBool("standards", false), "--standards")
	args = appendIfTrue(args, req.GetBool("tests", false), "--tests")
	args = appendIfTrue(args, req.GetBool("deps", false), "--deps")
	if out := req.GetString("output", ""); out != "" {
		args = append(args, "--output", out)
	}

	result, runErr := runGocaSubcommand(ctx, args)
	if runErr != nil {
		return mcpErr(runErr), nil
	}
	return mcpText(result), nil
}
