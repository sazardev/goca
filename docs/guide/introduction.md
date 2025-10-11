# What is Goca?

**Goca** (Go Clean Architecture) is a powerful CLI code generator that helps you build production-ready Go applications following **Clean Architecture** principles designed by Uncle Bob (Robert C. Martin).

## The Problem

Building Go applications with proper architecture is time-consuming:

- ❌ Writing repetitive boilerplate code
- ❌ Maintaining consistent structure across features
- ❌ Ensuring proper layer separation
- ❌ Setting up dependency injection
- ❌ Configuring routing and handlers
- ❌ Fighting architectural drift over time

## The Solution

Goca automates all of this while teaching you Clean Architecture:

```bash
# One command generates all layers properly structured
goca feature Product --fields "name:string,price:float64"
```

This creates:
- ✅ Domain entities with business validations
- ✅ Use cases with clear DTOs
- ✅ Repository interfaces and implementations
- ✅ HTTP handlers with proper routing
- ✅ Dependency injection automatically configured

## Core Philosophy

### 1. Clean Architecture by Default

Every line of code follows Uncle Bob's principles:

- **Dependency Rule**: Dependencies point inward toward domain
- **Layer Separation**: Clear boundaries between layers
- **Interface Segregation**: Small, focused contracts
- **Dependency Inversion**: Details depend on abstractions

### 2. Prevention Over Correction

Goca prevents common anti-patterns:

- 🚫 Fat Controllers - Business logic stays in use cases
- 🚫 God Objects - Each entity has single responsibility
- 🚫 Anemic Domain - Entities include business behavior
- 🚫 Direct Database Access - Always through repositories

### 3. Production-Ready Code

Generated code is not a starting point - it's production-ready:

- Error handling included
- Validation at proper layers
- Clear separation of concerns
- Testable by design
- Well-documented

## Key Features

### 🏗️ Complete Architecture Layers

```
┌─────────────────────────────────────┐
│         🟢 Handlers                 │  HTTP, gRPC, CLI
│  (Input/Output Adapters)            │
└─────────────────────────────────────┘
           ↓ depends on
┌─────────────────────────────────────┐
│      🔴 Use Cases (Application)     │  Business workflows
│  (DTOs, Services, Interfaces)       │
└─────────────────────────────────────┘
           ↓ depends on
┌─────────────────────────────────────┐
│    🔵 Repositories (Infrastructure) │  Data persistence
│  (Database implementations)         │
└─────────────────────────────────────┘
           ↓ depends on
┌─────────────────────────────────────┐
│        🟡 Domain (Entities)         │  Pure business logic
│  (No external dependencies)         │
└─────────────────────────────────────┘
```

### ⚡ Instant Feature Generation

Generate complete features in seconds:

```bash
# From zero to fully functional CRUD
goca feature Order --fields "customer:string,total:float64,status:string"
```

Creates 10+ files with:
- Domain entity
- CRUD use cases
- Repository interface + implementation
- HTTP REST endpoints
- Automatic integration

### 🎯 Multi-Protocol Support

Generate adapters for different protocols:

```bash
# HTTP REST API
goca handler Product --type http

# gRPC Service
goca handler Product --type grpc

# CLI Commands
goca handler Product --type cli

# Background Workers
goca handler Product --type worker

# SOAP Client
goca handler Product --type soap
```

All following the same clean architecture pattern!

### 🔄 Automatic Integration

New features are automatically integrated:

- Dependency injection containers updated
- Routes registered automatically
- Database connections configured
- No manual wiring needed

### 🧪 Test-Friendly

Generated code is designed for testing:

- Clear interfaces for mocking
- Dependency injection throughout
- Pure functions in domain
- Isolated layers

## Why Clean Architecture?

Clean Architecture provides:

### Maintainability
- Changes isolated to specific layers
- Clear boundaries prevent cascading effects
- Easy to understand and modify

### Testability
- Business logic independent of frameworks
- Easy to mock dependencies
- Fast unit tests without infrastructure

### Flexibility
- Swap implementations without touching business logic
- Add new delivery mechanisms easily
- Database agnostic domain layer

### Scalability
- Clear structure makes onboarding easy
- Consistent patterns across features
- Easy to add new features

## When to Use Goca?

### Perfect For:

- ✅ New Go projects requiring solid architecture
- ✅ Microservices with consistent structure
- ✅ REST APIs with multiple resources
- ✅ Projects that will grow over time
- ✅ Teams learning Clean Architecture
- ✅ MVPs that need to scale to production

### Maybe Not For:

- ❌ Simple scripts or one-off tools
- ❌ Extremely unique architectures
- ❌ Projects with existing different patterns

## Comparison

### Without Goca

```bash
# Manual process (hours of work)
1. Create domain entity
2. Write use case interfaces
3. Implement use case logic
4. Create DTOs for each operation
5. Write repository interface
6. Implement repository
7. Create HTTP handlers
8. Set up routing
9. Configure DI container
10. Wire everything together
11. Test and fix integration issues
```

### With Goca

```bash
# One command (seconds)
goca feature Product --fields "name:string,price:float64"

# Done! Everything wired and working
```

## Real-World Usage

Goca is used in production for:

- 🏢 Enterprise microservices
- 🛒 E-commerce platforms
- 📱 Mobile backend APIs
- 🔧 Internal tools and services
- 📊 Data processing pipelines

## Next Steps

Ready to start building clean Go applications?

- 📦 [Install Goca](/guide/installation)
- 🚀 [Quick Start Guide](/getting-started)
- 📖 [Learn Clean Architecture](/guide/clean-architecture)
- 🎓 [Complete Tutorial](/tutorials/complete-tutorial)

## Philosophy

> "The goal of software architecture is to minimize the human resources required to build and maintain the required system."
>
> — Robert C. Martin (Uncle Bob)

Goca embodies this philosophy by automating the tedious parts while maintaining architectural excellence.
