/*
Package cmd provides the command-line interface for the Goca CLI tool.

This package contains all the CLI commands and their implementations for
generating Go Clean Architecture projects. It uses the Cobra library for
command-line parsing and organization.

# Available Commands

The cmd package provides these main commands:

- init: Initialize a new Clean Architecture project
- feature: Generate complete feature with all layers
- entity: Generate domain entities with validation
- usecase: Generate use cases and business logic
- repository: Generate data access layer
- handler: Generate interface adapters
- di: Generate dependency injection container
- interfaces: Generate interfaces for TDD
- messages: Generate error messages and responses
- version: Show version information

Each command is implemented in its own file and provides specific functionality
for code generation following Clean Architecture principles.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goca",
	Short: "Go Clean Architecture Code Generator",
	Long: `Goca es un potente generador de código CLI para Go que te ayuda a crear 
proyectos con Clean Architecture siguiendo las mejores prácticas.

Genera código limpio y bien estructurado por capas, permitiéndote 
enfocarte en la lógica de negocio en lugar de tareas repetitivas de configuración.`,
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println(`Goca - Go Clean Architecture Code Generator

USAGE:
  goca [command]

AVAILABLE COMMANDS:
  help        Ayuda sobre cualquier comando
  version     Muestra la versión de Goca
  init        Inicializa un nuevo proyecto con Clean Architecture
  feature     Genera un feature completo con todas las capas
  entity      Genera entidades de dominio puras
  usecase     Genera casos de uso con DTOs
  handler     Genera handlers para diferentes protocolos
  repository  Genera repositorios con interfaces
  messages    Genera mensajes y constantes
  di          Genera contenedor de inyección de dependencias
  interfaces  Genera solo interfaces para TDD

Use "goca [command] --help" para más información sobre un comando.`)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(featureCmd)
	rootCmd.AddCommand(entityCmd)
	rootCmd.AddCommand(usecaseCmd)
	rootCmd.AddCommand(handlerCmd)
	rootCmd.AddCommand(repositoryCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(diCmd)
	rootCmd.AddCommand(interfacesCmd)
}
