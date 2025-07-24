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
	Short: "Generar mensajes y constantes",
	Long: `Crea archivos de mensajes de error, respuestas y constantes 
organizados por feature para mantener consistencia en la aplicación.`,
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

		fmt.Printf("Generando mensajes para entidad '%s'\n", entity)

		if errors {
			fmt.Println("✓ Generando mensajes de error")
		}
		if responses {
			fmt.Println("✓ Generando mensajes de respuesta")
		}
		if constants {
			fmt.Println("✓ Generando constantes")
		}

		generateMessages(entity, errors, responses, constants)
		fmt.Printf("\n✅ Mensajes para '%s' generados exitosamente!\n", entity)
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
		generateErrorMessages(messagesDir, entity)
	}

	if responses {
		generateResponseMessages(messagesDir, entity)
	}

	if constants {
		generateConstants(constantsDir, entity)
	}
}

func generateErrorMessages(dir, entity string) {
	filename := filepath.Join(dir, "errors.go")
	entityLower := strings.ToLower(entity)

	// Check if file exists and read existing content
	var existingContent strings.Builder
	if _, err := os.Stat(filename); err == nil {
		// File exists, read it
		if content, err := os.ReadFile(filename); err == nil {
			existing := string(content)
			// Remove the closing parenthesis and const block end
			if strings.Contains(existing, ")\n") {
				existing = strings.Replace(existing, ")\n", "", -1)
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

	// Add new entity error messages
	existingContent.WriteString(fmt.Sprintf("\t// %s errors\n", entity))
	existingContent.WriteString(fmt.Sprintf("\tErr%sNotFound        = \"%s not found\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErr%sAlreadyExists   = \"%s already exists\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErrInvalid%sData     = \"invalid %s data\"\n", entity, entityLower))

	// Field-specific errors
	existingContent.WriteString(fmt.Sprintf("\tErr%sEmailRequired   = \"%s email is required\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErr%sNameRequired    = \"%s name is required\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErr%sAgeInvalid      = \"%s age must be positive\"\n", entity, entityLower))

	// Business logic errors
	existingContent.WriteString(fmt.Sprintf("\tErr%sAccessDenied    = \"access denied to %s\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErr%sUpdateFailed    = \"failed to update %s\"\n", entity, entityLower))
	existingContent.WriteString(fmt.Sprintf("\tErr%sDeleteFailed    = \"failed to delete %s\"\n", entity, entityLower))

	// Close the const block
	existingContent.WriteString(")\n")

	writeGoFile(filename, existingContent.String())
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
				existing = strings.Replace(existing, ")\n", "", -1)
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

	writeGoFile(filename, existingContent.String())
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

	writeGoFile(filename, content.String())
}

func init() {
	messagesCmd.Flags().BoolP("errors", "e", false, "Generar mensajes de error")
	messagesCmd.Flags().BoolP("responses", "r", false, "Generar mensajes de respuesta")
	messagesCmd.Flags().BoolP("constants", "c", false, "Generar constantes del feature")
	messagesCmd.Flags().BoolP("all", "a", false, "Generar todos los tipos de mensajes")
}
