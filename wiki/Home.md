# Goca Wiki - Go Clean Architecture Code Generator

Welcome to the official **Goca** documentation! 🎉

Goca is a powerful CLI code generator for Go that helps you create projects following Uncle Bob's **Clean Architecture** principles. This wiki will guide you step by step to make the most of all Goca functionalities.

## 📚 Table of Contents

### 🚀 Quick Start
- [**Installation**](Installation) - How to install Goca on your system
- [**Getting Started**](Getting-Started) - Your first project with Goca
- [**Complete Tutorial**](Complete-Tutorial) - Step-by-step example of a real project

### 📖 Command Reference
- [**goca init**](Command-Init) - Initialize Clean Architecture projects
- [**goca feature**](Command-Feature) - Generate complete features
- [**goca entity**](Command-Entity) - Create domain entities
- [**goca usecase**](Command-UseCase) - Generate use cases
- [**goca repository**](Command-Repository) - Create repositories
- [**goca handler**](Command-Handler) - Generate input adapters
- [**goca di**](Command-DI) - Dependency injection
- [**goca integrate**](Command-Integrate) - Integrate existing features
- [**goca interfaces**](Command-Interfaces) - Generate interfaces for TDD
- [**goca messages**](Command-Messages) - Messages and constants
- [**goca ci**](Command-CI) - Generate CI/CD pipelines
- [**goca middleware**](Command-Middleware) - Generate HTTP middleware
- [**goca analyze**](Command-Analyze) - Deep project self-analysis (architecture, security, quality)
- [**goca version**](Command-Version) - Version information

### 🏗️ Architecture and Concepts
- [**Clean Architecture**](Clean-Architecture) - Principles and structure
- [**Project Structure**](Project-Structure) - Directory organization
- [**Database Support**](Database-Support) - 8 supported databases
- [**Implemented Patterns**](Design-Patterns) - Design patterns used
- [**Best Practices**](Best-Practices) - Recommendations and conventions

### 💡 Examples and Use Cases
- [**E-commerce Project**](Example-Ecommerce) - Complete e-commerce system
- [**REST API**](Example-REST-API) - RESTful API with multiple endpoints
- [**Microservice**](Example-Microservice) - Microservice with gRPC
- [**CLI Tool**](Example-CLI-Tool) - Command-line tool

### 🔧 Advanced
- [**Customization**](Customization) - Adapt templates to your needs
- [**CI/CD Integration**](CICD-Integration) - Automation and deployment
- [**Testing**](Testing-Guide) - Testing strategies with generated code
- [**Troubleshooting**](Troubleshooting) - Common problem solutions

### 🤝 Contributing
- [**Contributing Guide**](Contributing) - How to contribute to the project
- [**Development**](Development) - Set up development environment
- [**Roadmap**](Roadmap) - Future features

## 🎯 What is Clean Architecture?

Clean Architecture is an architectural pattern created by Robert C. Martin (Uncle Bob) that organizes code in concentric layers, where dependencies point towards the center of the system. This guarantees:

- ✅ **Framework independence** - Business code doesn't depend on external libraries
- ✅ **Testability** - Easy to test without external dependencies
- ✅ **UI independence** - The interface can change without affecting logic
- ✅ **Database independence** - Persistence is an implementation detail
- ✅ **External agent independence** - Business code doesn't know the outside world

## 🚀 Quick Start

### 1. Installation
```bash
go install github.com/sazardev/goca@latest
```

### 2. Create your first project
```bash
goca init my-project --module github.com/user/my-project
cd my-project
```

### 3. Generate your first feature
```bash
goca feature User --fields "name:string,email:string" --validation
```

### 4. Configure dependencies
```bash
goca di --features "User" --database postgres
```

You now have a complete project with Clean Architecture! 🎉

## 📈 Goca Philosophy

Goca doesn't just generate code, it **teaches** and **enforces** Clean Architecture best practices:

- **🟡 Domain** → Pure entities without external dependencies
- **🔴 Use Cases** → Application logic with clear DTOs
- **🟢 Adapters** → Interfaces that adapt input/output
- **🔵 Infrastructure** → Technology-specific implementations

## 🛡️ Quality Guarantees

- ✅ **Dependencies directed towards the center**
- ✅ **Clear interfaces between layers**
- ✅ **Separation of responsibilities**
- ✅ **Testable code by design**
- ✅ **Production-proven patterns**

## 🌟 Featured Characteristics

- **Multi-protocol**: HTTP REST, gRPC, CLI, Workers, SOAP
- **Multi-database**: PostgreSQL, MySQL, MongoDB
- **Dependency injection**: Manual and with Wire.dev
- **Validations**: In domain and DTOs
- **Testing**: Interfaces for TDD
- **Documentation**: Automatic Swagger

## 📞 Support and Community

- 🐛 **Issues**: [GitHub Issues](https://github.com/sazardev/goca/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/sazardev/goca/discussions)
- 📧 **Contact**: [sazardev@email.com](mailto:sazardev@email.com)

---

**Explore the documentation and start creating amazing projects with Clean Architecture!** 🚀
