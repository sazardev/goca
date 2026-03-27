# Roadmap

This document outlines the planned features and improvements for Goca.

## Vision

Goca aims to be the leading Go Clean Architecture code generator, providing developers with:
- Production-ready code generation
- Best practice enforcement
- Comprehensive testing support
- Multi-protocol support
- Enterprise-grade features
- AI-powered developer experience

## Current Status

**Version**: 1.19.0  
**Go**: 1.25.1+  
**Status**: Production Ready  
**Focus**: Cache, CI, middleware generation, MCP integration, and developer experience

## Release Versioning

Goca follows [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes and minor improvements

## Planned Releases

### Version 1.19.0

**Focus**: Cache, CI & Middleware

**Planned Features:**
- Plugin system for custom generators
- GraphQL handler generation
- Event sourcing patterns
- WebSocket handler support

**Improvements:**
- Observability instrumentation (OpenTelemetry)
- E2E test generation
- Contract testing support

### Version 2.0.0

**Focus**: Major Architecture Enhancement

**Planned Features:**
- Plugin system for custom generators
- Custom template engine
- GraphQL handler generation
- Event sourcing patterns
- CQRS advanced patterns
- WebSocket handler support
- Message queue integration

**Breaking Changes:**
- Configuration file format update
- Template system overhaul

**Migration Guide:** Will be provided with detailed instructions

## Feature Backlog

### High Priority

**Code Generation:**
- WebSocket handler support
- Message queue integration (RabbitMQ, Kafka)
- Observability instrumentation (OpenTelemetry)

**Testing:**
- E2E test generation
- Contract testing support
- Load test templates

### Medium Priority

**Documentation:**
- Video tutorials
- Architecture decision records
- Migration guides from other frameworks

**Developer Experience:**
- Kubernetes deployment generation
- Benchmark test templates

**Quality Assurance:**
- Static analysis integration
- Security scanning
- Dependency vulnerability checks

### Low Priority (Future Consideration)

**Advanced Features:**
- Multi-language support (TypeScript, Python)
- Service mesh integration
- Infrastructure as code generation

**Community Features:**
- Template marketplace
- Community plugins
- Shared configurations
- Example project gallery

## Completed Milestones

### Version 1.19.0 (March 2026)
- **Redis Cache Layer** — `--cache` flag on `feature`, `repository`, `di` generates `Cached<Entity>Repository` decorator with Redis caching (FindByID, FindAll cached; writes invalidate)
- **CI Pipeline Generation** — `goca ci` generates GitHub Actions workflows (test, build, deploy) with auto-detected Go version and database service containers
- **Middleware Generation** — `goca middleware <name>` generates 7 composable middleware types (CORS, logging, auth, rate-limit, recovery, request-id, timeout)
- Handler auto-detection of middleware package for import-based usage
- `--middleware-types` flag on `goca feature` for combined scaffold
- DI container cache wiring: `NewContainer(db, redisClient)` with decorator pattern
- 51 new tests across 4 test files
- MCP tools for CI and middleware commands

### Version 1.18.x (March 2026)
- **MCP Server** — `goca mcp-server` exposes 13 tools + 2 resources for AI assistants (GitHub Copilot, Claude Desktop, Cursor, Zed)
- Interactive project initialization wizard (`goca init` with TUI via charmbracelet/huh)
- `goca doctor` command for project health checks
- `goca upgrade` command for config/metadata upgrades
- `goca template` command for template management
- Multi-handler support: HTTP, gRPC, CLI, Worker, SOAP
- DynamoDB repository generation
- Elasticsearch repository generation
- PostgreSQL JSON/JSONB repository generation
- SQL Server repository generation
- Enhanced charmbracelet/lipgloss TUI styling

### Version 1.17.x (January–February 2026)
- Mock generation for all interfaces (`goca mocks`)
- Integration test scaffolding (`goca test-integration`)
- Entity test generation with table-driven tests
- File naming convention support (snake_case, kebab-case)
- Auto-migration generation

### Version 1.14.x (October 2025)
- Unit test generation for entities
- Integration test improvements
- Enhanced error messages
- Better dry-run output

### Version 1.13.x (October 2024)
- Project template system (minimal, rest-api, microservice, monolith, enterprise)
- Predefined configurations via `.goca.yaml`
- Bug fixes for GORM integration
- Improved validation handling

### Version 1.11.0 (January 2024)
- Safety features (dry-run, backups, force mode)
- Dependency management (`go mod tidy` integration)
- File conflict detection
- Name conflict detection

### Version 1.10.0 (December 2023)
- Complete feature generation (`goca feature`)
- Automatic DI + routing integration (`goca integrate`)
- HTTP handler support with gorilla/mux
- Repository pattern implementation (PostgreSQL, MySQL, SQLite)
- MongoDB repository templates

### Version 1.0.0 (November 2023)
- Initial stable release
- Entity, UseCase, Repository, Handler generation
- Clean Architecture directory structure
- Go module creation and dependency management
- Core CLI commands
- Clean Architecture scaffolding
- Basic code generation

## Maintenance Windows

### Long-Term Support

- **Current Release (1.13.x)**: Full support
- **Previous Release (1.10.x)**: Critical fixes only
- **Older Releases**: Community support

### End of Life Policy

- Minor versions: Supported for 6 months after next release
- Major versions: Supported for 12 months after next major release
- Security patches: Provided for supported versions only

## Breaking Changes Policy

**Major Versions Only:**
- Breaking changes only in major versions
- Deprecation warnings in prior minor versions
- Migration guides provided
- Minimum 6-month deprecation period

**Communication:**
- Announced in release notes
- Documented in migration guides
- Discussed in community forums
- Early preview releases available

## Contributing to the Roadmap

### How to Contribute

1. **Feature Requests**: Open an issue with detailed use case
2. **Discussions**: Participate in design discussions
3. **Voting**: React to issues to show interest
4. **Implementation**: Submit pull requests for approved features

### Feature Proposal Process

1. **Proposal**: Create feature request issue
2. **Discussion**: Community feedback (2-4 weeks)
3. **Review**: Maintainer evaluation
4. **Planning**: Schedule for release
5. **Implementation**: Development and testing
6. **Release**: Include in appropriate version

### Roadmap Updates

This roadmap is reviewed and updated:
- Quarterly for near-term plans
- Bi-annually for long-term vision
- Ad-hoc for significant changes

**Last Updated**: October 15, 2025  
**Next Review**: January 15, 2026

## Questions and Feedback

Have thoughts on the roadmap?

- **GitHub Discussions**: Share ideas and feedback
- **Feature Requests**: Propose specific features
- **Email**: sazardev@gmail.com for detailed discussions

## Commitment

While we strive to deliver on this roadmap, timelines and features may change based on:
- Community feedback and priorities
- Resource availability
- Technical constraints
- Ecosystem changes

We're committed to transparency and will communicate changes promptly.

---

Thank you for your interest in Goca's future!
