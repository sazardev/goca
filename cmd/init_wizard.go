package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// runInitWizard launches an interactive form to gather project init parameters.
// Returns module, database, api, auth, config values or an error if cancelled.
func runInitWizard(projectName string) (module, database, api string, auth, config bool, err error) {
	module = fmt.Sprintf("github.com/user/%s", projectName)
	database = "sqlite"
	api = "rest"
	auth = false
	config = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Go module path").
				Description("The import path for your Go module").
				Placeholder(fmt.Sprintf("github.com/user/%s", projectName)).
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

	err = form.Run()
	if err != nil {
		return "", "", "", false, false, err
	}

	return module, database, api, auth, config, nil
}
