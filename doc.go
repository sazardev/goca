// Package goca provides a CLI tool for generating Go Clean Architecture projects.
//
// Goca is a comprehensive CLI tool that generates well-structured Go projects following
// Clean Architecture principles. It creates projects with zero compilation errors,
// proper file organization, and clean code that passes all linting checks.
//
// # Features
//
//   - Generate complete Go Clean Architecture projects
//   - Support for multiple databases (PostgreSQL, MySQL, MongoDB)
//   - Built-in authentication and authorization
//   - REST API and gRPC support
//   - Comprehensive validation and business rules
//   - Zero-error code generation
//   - Clean and properly formatted output
//
// # Installation
//
//	go install github.com/jorgefuertes/goca@latest
//
// # Quick Start
//
//	# Initialize a new project
//	goca init myproject --module=github.com/user/myproject --database=postgres --auth
//
//	# Generate a complete feature
//	goca feature User --fields="name:string,email:string,age:int" --validation --business-rules
//
//	# Generate individual components
//	goca entity Product --fields="name:string,price:float64" --validation --timestamps
//	goca usecase User --operations="create,read,update,delete,list" --validation
//	goca repository User --database=postgres --cache --transactions
//	goca handler User --protocol=http --validation --middleware
//
// # Commands
//
//   - init: Initialize a new Clean Architecture project
//   - feature: Generate a complete feature (entity + usecase + repository + handler)
//   - entity: Generate domain entities with validation
//   - usecase: Generate use case layer with interfaces
//   - repository: Generate repository layer with database implementation
//   - handler: Generate handler layer for different protocols
//   - messages: Generate error messages and response structures
//   - di: Generate dependency injection container
//   - interfaces: Generate interfaces for all layers
//
// # Architecture Layers
//
//   - Domain: Core business entities and rules
//   - UseCase: Application business logic and workflows
//   - Repository: Data access and persistence
//   - Handler: External interfaces (HTTP, gRPC, CLI)
//   - Infrastructure: External dependencies and configuration
//
// # Code Quality
//
// All generated code is guaranteed to:
//   - Compile without errors or warnings
//   - Pass go vet and golangci-lint checks
//   - Follow Go naming conventions and best practices
//   - Include proper error handling and validation
//   - Be properly formatted with gofmt
//   - Have comprehensive test coverage
//
// For more information and documentation, visit:
// https://github.com/jorgefuertes/goca
package main
