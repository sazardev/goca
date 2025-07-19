/*
Package main provides the Goca CLI tool for generating Go Clean Architecture projects.

Goca (Go Clean Architecture) is a powerful command-line tool that generates
well-structured Go applications following Clean Architecture principles.
It creates complete projects with domain entities, use cases, repositories,
handlers, and proper dependency injection.

# Features

- Generate complete Clean Architecture projects
- Support for multiple databases (PostgreSQL, MySQL, MongoDB)
- Multiple handler types (HTTP, gRPC, CLI, Worker)
- Automatic dependency injection container generation
- Built-in validation and error handling
- Business rules and domain logic templates
- Message and response management

# Installation

	go install github.com/sazardev/goca@latest

# Quick Start

	# Initialize a new project
	goca init myproject --module github.com/user/myproject --database postgres

	# Generate a complete feature
	goca feature User --fields "name:string,email:string" --database postgres

	# Generate dependency injection
	goca di --features "User" --database postgres

# Project Structure

Goca generates projects following this structure:

	project/
	├── cmd/server/main.go           # Application entry point
	├── internal/
	│   ├── domain/                  # Business entities and rules
	│   ├── usecase/                 # Application business logic
	│   ├── repository/              # Data access layer
	│   ├── handler/                 # Interface adapters
	│   └── infrastructure/          # Framework and drivers
	└── pkg/                         # Shared utilities

# Commands

The main commands available are:

- init: Initialize a new Clean Architecture project
- feature: Generate complete feature with all layers
- entity: Generate domain entities with validation
- usecase: Generate use cases and business logic
- repository: Generate data access layer
- handler: Generate interface adapters (HTTP, gRPC, CLI)
- di: Generate dependency injection container
- interfaces: Generate interfaces for TDD
- messages: Generate error messages and responses

For detailed usage of each command, use:

	goca <command> --help

# Architecture Principles

Goca enforces Clean Architecture principles:

- Dependency Rule: Dependencies point inward toward the domain
- Interface Segregation: Small, focused interfaces
- Dependency Inversion: Depend on abstractions, not concretions
- Single Responsibility: Each layer has one reason to change

# Examples

	# E-commerce project
	goca init ecommerce --module github.com/user/ecommerce --database postgres --auth
	goca feature Product --fields "name:string,price:float64,stock:int"
	goca feature Order --fields "user_id:int,total:float64,status:string"
	goca di --features "Product,Order" --database postgres

	# Microservice
	goca init user-service --module github.com/user/user-service --api grpc
	goca feature User --fields "name:string,email:string" --database postgres
	goca handler User --type grpc --validation

For more examples and detailed guides, visit:
https://github.com/sazardev/goca/blob/main/README.md
*/
package main
