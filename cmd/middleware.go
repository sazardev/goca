package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// validMiddlewareTypes lists all supported middleware type keys.
var validMiddlewareTypes = []string{
	"cors", "logging", "auth", "rate-limit", "recovery", "request-id", "timeout",
}

var middlewareCmd = &cobra.Command{
	Use:   "middleware <name>",
	Short: "Generate middleware package for HTTP handlers",
	Long: `Generate a dedicated internal/middleware/ package with composable HTTP middleware.

Supported middleware types:
  cors        — Configurable CORS headers (origins, methods, headers)
  logging     — Structured request logging with duration and status
  auth        — JWT token validation with claims extraction
  rate-limit  — Token bucket rate limiter per IP
  recovery    — Panic recovery returning JSON 500 responses
  request-id  — Inject X-Request-ID into context and response
  timeout     — Per-request context deadline

Use --types to select which middleware to generate (comma-separated).
Default: cors,logging,recovery

A middleware.go file with a Chain() helper is always generated for composing
multiple middleware functions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		validator := NewFieldValidator()
		if err := validator.ValidateEntityName(name); err != nil {
			return err
		}

		typesStr, _ := cmd.Flags().GetString("types")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")

		types := parseMiddlewareTypes(typesStr)
		if err := validateMiddlewareTypes(types); err != nil {
			return err
		}

		sm := NewSafetyManager(dryRun, force, backup)

		ui.Header("Goca Middleware — Package Generation")
		ui.Blank()
		ui.KeyValue("Name", name)
		ui.KeyValue("Types", strings.Join(types, ", "))
		ui.Blank()

		if err := generateMiddlewarePackage(name, types, sm); err != nil {
			return err
		}

		if dryRun {
			sm.PrintSummary()
			return nil
		}

		ui.Blank()
		ui.Success("Middleware package generated successfully!")
		ui.Blank()
		ui.Info("Next steps:")
		ui.Step(1, "Import the middleware package in your routes")
		ui.Step(2, "Use middleware.Chain() to compose middleware functions")
		return nil
	},
}

func init() {
	middlewareCmd.Flags().String("types", "cors,logging,recovery", "Middleware types to generate (comma-separated)")
	middlewareCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	middlewareCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	middlewareCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}

// parseMiddlewareTypes splits a comma-separated types string and normalizes them.
func parseMiddlewareTypes(typesStr string) []string {
	raw := strings.Split(typesStr, ",")
	var types []string
	for _, t := range raw {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" {
			types = append(types, t)
		}
	}
	return types
}

// validateMiddlewareTypes checks that all requested types are supported.
func validateMiddlewareTypes(types []string) error {
	valid := make(map[string]bool, len(validMiddlewareTypes))
	for _, v := range validMiddlewareTypes {
		valid[v] = true
	}
	for _, t := range types {
		if !valid[t] {
			return fmt.Errorf("unsupported middleware type %q; valid types: %s", t, strings.Join(validMiddlewareTypes, ", "))
		}
	}
	return nil
}
