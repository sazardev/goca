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
	Long: `Goca is a powerful CLI code generator for Go that helps you create 
Clean Architecture projects following best practices.

It generates clean, well-structured layered code, allowing you to 
focus on business logic instead of repetitive configuration tasks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
	rootCmd.AddCommand(integrateCmd)
	rootCmd.AddCommand(interfacesCmd)
}
