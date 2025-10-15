# Roadmap

This document outlines the planned features and improvements for Goca.

## Vision

Goca aims to be the leading Go Clean Architecture code generator, providing developers with:
- Production-ready code generation
- Best practice enforcement
- Comprehensive testing support
- Multi-protocol support
- Enterprise-grade features

## Current Status

**Version**: 1.13.x  
**Status**: Production Ready  
**Focus**: Stability, documentation, and community growth

## Release Versioning

Goca follows [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes and minor improvements

## Planned Releases

### Version 1.14.0 (Q1 2026)

**Focus**: Enhanced Testing Support

**Planned Features:**
- Unit test generation for entities
- Integration test scaffolding
- Mock generation for interfaces
- Test coverage reporting
- Benchmark test templates

**Improvements:**
- Enhanced error messages
- Better dry-run output
- Performance optimizations

### Version 1.15.0 (Q2 2026)

**Focus**: Advanced Database Support

**Planned Features:**
- MongoDB repository templates
- PostgreSQL-specific optimizations
- Redis caching layer generation
- Database migration helpers
- Multi-database support

**Improvements:**
- Enhanced SQL query generation
- Better transaction handling
- Connection pool configuration

### Version 2.0.0 (Q3 2026)

**Focus**: Major Architecture Enhancement

**Planned Features:**
- Plugin system for custom generators
- Custom template engine
- GraphQL support
- Event sourcing patterns
- CQRS advanced patterns

**Breaking Changes:**
- Configuration file format update
- Command structure reorganization
- Template system overhaul

**Migration Guide:** Will be provided with detailed instructions

## Feature Backlog

### High Priority

**CLI Improvements:**
- Interactive mode for feature generation
- Project scaffolding wizard
- Configuration file generator
- Better error recovery

**Code Generation:**
- WebSocket handler support
- Message queue integration
- Observability instrumentation
- Security middleware generation

**Testing:**
- E2E test generation
- Contract testing support
- Load test templates
- Test data generators

### Medium Priority

**Documentation:**
- Video tutorials
- Interactive documentation
- Architecture decision records
- Migration guides from other frameworks

**Developer Experience:**
- VS Code extension
- IntelliJ plugin
- GitHub Actions templates
- Docker development environments

**Quality Assurance:**
- Static analysis integration
- Security scanning
- Dependency vulnerability checks
- Performance profiling

### Low Priority (Future Consideration)

**Advanced Features:**
- Multi-language support (TypeScript, Python)
- Kubernetes deployment generation
- Service mesh integration
- Infrastructure as code generation

**Community Features:**
- Template marketplace
- Community plugins
- Shared configurations
- Example project gallery

## Research and Exploration

**Topics Under Investigation:**

1. **AI-Assisted Generation**
   - Natural language to code generation
   - Intelligent field type detection
   - Automatic relationship mapping

2. **Advanced Patterns**
   - Domain-Driven Design patterns
   - Hexagonal architecture variants
   - Microservices patterns

3. **Performance Optimization**
   - Faster code generation
   - Parallel processing
   - Memory optimization

4. **Developer Tooling**
   - Live reload during development
   - Code refactoring tools
   - Migration helpers

## Community Requests

**Top Requested Features:**

Based on community feedback and GitHub issues:

1. Better support for existing projects (migration tools)
2. More database options (MongoDB, Redis)
3. GraphQL API generation
4. WebSocket support
5. Event-driven architecture templates

**Contributing:**

Want to influence the roadmap? 
- Create feature requests on GitHub
- Participate in discussions
- Vote on existing proposals
- Submit pull requests

## Completed Milestones

### Version 1.13.x (October 2024)
- Project template system
- Predefined configurations
- Bug fixes for GORM integration
- Improved validation handling

### Version 1.11.0 (January 2024)
- Safety features (dry-run, backups)
- Dependency management
- File conflict detection
- Name conflict detection

### Version 1.10.0 (December 2023)
- Complete feature generation
- Automatic integration
- HTTP handler support
- Repository pattern implementation

### Version 1.0.0 (November 2023)
- Initial stable release
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
