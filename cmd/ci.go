package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Generate CI/CD pipeline configuration",
	Long: `Generate Continuous Integration pipeline files for your project.

Currently supported providers:
  - github-actions  GitHub Actions workflows (.github/workflows/)

Generated workflows include:
  - test.yml   — Run tests, vet, and build on every push/PR
  - build.yml  — Build binary artifacts, optionally build Docker image
  - deploy.yml — Tag-triggered deployment (optional, use --with-deploy)

The command reads go.mod to detect the Go version and .goca.yaml to detect
the database driver, then generates provider-specific configuration files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, _ := cmd.Flags().GetString("provider")
		withDocker, _ := cmd.Flags().GetBool("with-docker")
		withDeploy, _ := cmd.Flags().GetBool("with-deploy")
		goVersion, _ := cmd.Flags().GetString("go-version")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")
		backup, _ := cmd.Flags().GetBool("backup")

		sm := NewSafetyManager(dryRun, force, backup)

		data := buildCITemplateData(goVersion)
		data.WithDocker = withDocker
		data.WithDeploy = withDeploy

		ui.Header("Goca CI — Pipeline Generation")
		ui.Blank()
		ui.KeyValue("Provider", provider)
		ui.KeyValue("Go version", data.GoVersion)
		ui.KeyValue("Docker", fmt.Sprintf("%v", withDocker))
		ui.KeyValue("Deploy", fmt.Sprintf("%v", withDeploy))
		if data.Database != "" {
			ui.KeyValue("Database", data.Database)
		}
		ui.Blank()

		if err := generateCIPipeline(provider, data, sm); err != nil {
			return err
		}

		if dryRun {
			sm.PrintSummary()
			return nil
		}

		ui.Blank()
		ui.Success("CI pipeline generated successfully!")
		ui.Blank()
		ui.Info("Next steps:")
		ui.Step(1, "Commit the generated workflow files")
		ui.Step(2, "Push to trigger your first CI run")
		return nil
	},
}

func init() {
	ciCmd.Flags().String("provider", "github-actions", "CI provider (github-actions)")
	ciCmd.Flags().Bool("with-docker", false, "Include Docker build step in build workflow")
	ciCmd.Flags().Bool("with-deploy", false, "Generate deploy workflow (tag-triggered)")
	ciCmd.Flags().String("go-version", "", "Go version for CI matrix (default: read from go.mod)")
	ciCmd.Flags().Bool("dry-run", false, "Preview changes without creating files")
	ciCmd.Flags().Bool("force", false, "Overwrite existing files without asking")
	ciCmd.Flags().Bool("backup", false, "Backup existing files before overwriting")
}
