package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var messagesCmd = &cobra.Command{
	Use:   "messages <entity>",
	Short: "Generate messages and constants",
	Long: `Creates error message, response and constant files 
organized by feature to maintain consistency in the application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		errors, _ := cmd.Flags().GetBool("errors")
		responses, _ := cmd.Flags().GetBool("responses")
		constants, _ := cmd.Flags().GetBool("constants")
		all, _ := cmd.Flags().GetBool("all")

		// If all is true, enable all options
		if all {
			errors = true
			responses = true
			constants = true
		}

		// If no specific flags, generate all by default
		if !errors && !responses && !constants {
			errors = true
			responses = true
			constants = true
		}

		fmt.Printf("Generating messages for entity '%s'\n", entity)

		if errors {
			fmt.Println("✓ Generating error messages")
		}
		if responses {
			fmt.Println("✓ Generating response messages")
		}
		if constants {
			fmt.Println("✓ Generating constants")
		}

		generateMessages(entity, errors, responses, constants)
		fmt.Printf("\nMessages for '%s' generated successfully!\n", entity)
	},
}

func generateMessages(entity string, errors, responses, constants bool) {
	// Create messages directory
	messagesDir := filepath.Join("internal", "messages")
	_ = os.MkdirAll(messagesDir, 0755)

	// Create constants directory
	constantsDir := filepath.Join("internal", "constants")
	_ = os.MkdirAll(constantsDir, 0755)

	if errors {
		generateUseCaseMessages(entity)
	}

	if responses {
		generateResponseMessages(messagesDir, entity)
	}

	if constants {
		generateConstants(constantsDir, entity)
	}
}

func generateUseCaseMessages(entity string) {
	// Create messages directory and file in internal/messages
	messagesDir := filepath.Join("internal", "messages")
	if err := os.MkdirAll(messagesDir, 0755); err != nil {
		fmt.Printf("Error creating messages directory: %v\n", err)
		return
	}

	filename := filepath.Join(messagesDir, "messages.go")
	entityLower := strings.ToLower(entity)

	// Check if file exists and read existing content
	var existingContent strings.Builder
	if _, err := os.Stat(filename); err == nil {
		// File exists, read it
		if content, err := os.ReadFile(filename); err == nil {
			existing := string(content)
			// Remove the closing parenthesis and const block end
			if strings.Contains(existing, ")\n") {
				existing = strings.ReplaceAll(existing, ")\n", "")
				existingContent.WriteString(existing)
			} else {
				// Start fresh if format is unexpected
				existingContent.WriteString("package messages\n\nconst (\n")
			}
		} else {
			// Error reading, start fresh
			existingContent.WriteString("package messages\n\nconst (\n")
		}
	} else {
		// File doesn't exist, start fresh
		existingContent.WriteString("package messages\n\nconst (\n")
	}

	// Add new entity messages
	existingContent.WriteString(fmt.Sprintf("\t// %s messages\n", entity))
	existingContent.WriteString(fmt.Sprintf("\t%sCreated = \"%s created successfully\"\n", entity, entity))
	existingContent.WriteString(fmt.Sprintf("\t%sNotFound = \"%s not found\"\n", entity, entity))
	existingContent.WriteString(fmt.Sprintf("\t%sUpdated = \"%s updated successfully\"\n", entity, entity))
	existingContent.WriteString(fmt.Sprintf("\t%sDeleted = \"%s deleted successfully\"\n", entity, entity))
	existingContent.WriteString(fmt.Sprintf("\t%sInvalid = \"Invalid %s data\"\n", entity, entityLower))
	existingContent.WriteString(")\n")

	if err := writeGoFile(filename, existingContent.String()); err != nil {
		fmt.Printf("Error writing messages file: %v\n", err)
		return
	}
}

func generateResponseMessages(dir, entity string) {
	filename := filepath.Join(dir, "responses.go")
	entityLower := strings.ToLower(entity)

	// Check if file exists and read existing content
	var existingContent strings.Builder
	if _, err := os.Stat(filename); err == nil {
		// File exists, read it
		if content, err := os.ReadFile(filename); err == nil {
			existing := string(content)
			// Remove the closing parenthesis and const block end
			if strings.Contains(existing, ")\n") {
				existing = strings.ReplaceAll(existing, ")\n", "")
				existingContent.WriteString(existing)
			} else {
				// Start fresh if format is unexpected
				existingContent.WriteString("package messages\n\nconst (\n")
			}
		} else {
			// Error reading, start fresh
			existingContent.WriteString("package messages\n\nconst (\n")
		}
	} else {
		// File doesn't exist, start fresh
		existingContent.WriteString("package messages\n\nconst (\n")
	}

	// Add new entity messages
	existingContent.WriteString(fmt.Sprintf("\t// %s success messages\n", entity))
	existingContent.WriteString(fmt.Sprintf("\t%sCreatedSuccessfully = \"%s created successfully\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%sUpdatedSuccessfully = \"%s updated successfully\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%sDeletedSuccessfully = \"%s deleted successfully\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%sFoundSuccessfully   = \"%s found successfully\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%ssListedSuccessfully = \"%ss listed successfully\"\n", entity, entityLower))

	// Operation messages
	existingContent.WriteString(fmt.Sprintf("\t%sProcessingStarted   = \"%s processing started\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%sProcessingCompleted = \"%s processing completed\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\t%sValidationPassed    = \"%s validation passed\"\n", entity, entityLower))

	// Close the const block
	existingContent.WriteString(")\n")

	if err := writeGoFile(filename, existingContent.String()); err != nil {
		fmt.Printf("Error creating response messages file: %v\n", err)
	}
}

func generateConstants(dir, entity string) {
	filename := filepath.Join(dir, "constants.go")
	entityLower := strings.ToLower(entity)

	var content strings.Builder
	content.WriteString("package constants\n\n")
	content.WriteString("const (\n")

	// Entity constants
	content.WriteString(fmt.Sprintf("\t// %s constants\n", entity))

	// Validation constants
	content.WriteString(fmt.Sprintf("\tMin%sAge        = 0\n", entity))
	content.WriteString(fmt.Sprintf("\tMax%sAge        = 150\n", entity))
	content.WriteString(fmt.Sprintf("\tMin%sNameLength = 2\n", entity))
	content.WriteString(fmt.Sprintf("\tMax%sNameLength = 100\n", entity))

	// Database constants
	content.WriteString(fmt.Sprintf("\t%sTableName     = \"%ss\"\n", entity, entityLower))
	content.WriteString(fmt.Sprintf("\t%sIDColumn      = \"id\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sNameColumn    = \"name\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sEmailColumn   = \"email\"\n", entity))

	// Cache constants
	content.WriteString(fmt.Sprintf("\t%sCachePrefix   = \"%s:\"\n", entity, entityLower))
	content.WriteString(fmt.Sprintf("\t%sCacheTTL      = 300 // 5 minutes\n", entity))

	// API constants
	content.WriteString(fmt.Sprintf("\t%sAPIVersion    = \"v1\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sEndpoint      = \"/%ss\"\n", entity, entityLower))
	content.WriteString(fmt.Sprintf("\tMax%sPerPage    = 100\n", entity))
	content.WriteString(fmt.Sprintf("\tDefault%sPerPage = 20\n", entity))

	content.WriteString(")\n\n")

	// Status constants
	content.WriteString("// Status constants\n")
	content.WriteString("const (\n")
	content.WriteString(fmt.Sprintf("\t%sStatusActive   = \"active\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sStatusInactive = \"inactive\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sStatusPending  = \"pending\"\n", entity))
	content.WriteString(fmt.Sprintf("\t%sStatusDeleted  = \"deleted\"\n", entity))
	content.WriteString(")\n")

	if err := writeGoFile(filename, content.String()); err != nil {
		fmt.Printf("Error creating constants file: %v\n", err)
	}
}

func init() {
	messagesCmd.Flags().BoolP("errors", "e", false, "Generate error messages")
	messagesCmd.Flags().BoolP("responses", "r", false, "Generate response messages")
	messagesCmd.Flags().BoolP("constants", "c", false, "Generate feature constants")
	messagesCmd.Flags().BoolP("all", "a", false, "Generate all message types")
}
