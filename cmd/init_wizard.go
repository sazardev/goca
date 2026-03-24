package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// runInitWizard launches an interactive form to gather project init parameters.
// When projectName is empty, also prompts for the project name.
// Returns name, module, database, api, auth, config values or an error if cancelled.
func runInitWizard(projectName string) (name, module, database, api string, auth, config bool, err error) {
	name = projectName
	modulePlaceholder := "github.com/user/myproject"
	if projectName != "" {
		module = fmt.Sprintf("github.com/user/%s", projectName)
		modulePlaceholder = module
	}
	database = "sqlite"
	api = "rest"
	auth = false
	config = true

	var formGroups []*huh.Group

	// When no project name was provided on the CLI, ask for it first
	if projectName == "" {
		formGroups = append(formGroups, huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("Directory name for your new project").
				Placeholder("myproject").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name is required")
					}
					return nil
				}).
				Value(&name),
		).Title("New Project"))
	}

	formGroups = append(formGroups,
		huh.NewGroup(
			huh.NewInput().
				Title("Go module path").
				Description("The import path for your Go module").
				Placeholder(modulePlaceholder).
				Value(&module),

			huh.NewSelect[string]().
				Title("Database backend").
				Options(
					huh.NewOption("SQLite", "sqlite"),
					huh.NewOption("PostgreSQL", "postgres"),
					huh.NewOption("MySQL", "mysql"),
					huh.NewOption("MongoDB", "mongodb"),
					huh.NewOption("SQL Server", "sqlserver"),
					huh.NewOption("DynamoDB", "dynamodb"),
					huh.NewOption("Elasticsearch", "elasticsearch"),
				).
				Value(&database),

			huh.NewSelect[string]().
				Title("API style").
				Options(
					huh.NewOption("REST", "rest"),
					huh.NewOption("gRPC", "grpc"),
					huh.NewOption("GraphQL", "graphql"),
				).
				Value(&api),
		).Title("Project Setup"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Include authentication scaffold?").
				Value(&auth),

			huh.NewConfirm().
				Title("Generate .goca.yaml config file?").
				Affirmative("Yes").
				Negative("No").
				Value(&config),
		).Title("Options"),
	)

	err = huh.NewForm(formGroups...).Run()
	if err != nil {
		return "", "", "", "", false, false, err
	}

	// When projectName was collected interactively, use it as module default if left blank
	if module == "" && name != "" {
		module = fmt.Sprintf("github.com/user/%s", name)
	}

	return name, module, database, api, auth, config, nil
}
