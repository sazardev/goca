# Support

Thank you for using Goca. This document provides information on how to get help and support.

## Documentation

Before seeking support, please check our comprehensive documentation:

- **[Official Documentation](https://sazardev.github.io/goca)** - Complete guides and tutorials
- **[Getting Started Guide](https://sazardev.github.io/goca/getting-started)** - Quick start in 5 minutes
- **[Commands Reference](https://sazardev.github.io/goca/commands/)** - Detailed command documentation
- **[Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial)** - Build a real application
- **[Clean Architecture Guide](https://sazardev.github.io/goca/guide/clean-architecture)** - Understand the principles

## Common Issues

### Installation Issues

**Problem: `goca: command not found`**
- Ensure the binary is in your PATH
- Verify installation with `which goca` (Unix) or `where goca` (Windows)
- Try reinstalling from [GitHub Releases](https://github.com/sazardev/goca/releases)

**Problem: Version shows as "dev"**
- You may have installed via `go install` instead of binary release
- Download the official binary from GitHub Releases for proper version info

### Generation Issues

**Problem: Files not being generated**
- Check that you're in the project root directory
- Verify the project was initialized with `goca init`
- Ensure you have write permissions in the directory

**Problem: Import errors in generated code**
- Run `go mod tidy` to update dependencies
- Verify your Go version is 1.21 or higher
- Check that the module path in `go.mod` is correct

**Problem: Duplicate entity names**
- Use `--force` flag to overwrite existing files
- Consider using `--backup` to preserve old files
- Check for case-insensitive name conflicts

## Getting Help

### Search Existing Resources

Before creating a new issue:

1. **Search Existing Issues**: Check if your problem has been reported
   - [Open Issues](https://github.com/sazardev/goca/issues)
   - [Closed Issues](https://github.com/sazardev/goca/issues?q=is%3Aissue+is%3Aclosed)

2. **Check Discussions**: Browse community discussions
   - [GitHub Discussions](https://github.com/sazardev/goca/discussions)

3. **Review Documentation**: Consult the docs site
   - [Documentation Website](https://sazardev.github.io/goca)

### Ask Questions

For general questions and discussions:

**GitHub Discussions** (Recommended)
- [Start a Discussion](https://github.com/sazardev/goca/discussions/new)
- Ask questions, share ideas, show your projects
- Get help from the community

**Topics:**
- General questions about usage
- Design and architecture discussions
- Feature ideas and suggestions
- Share your experience with Goca

### Report Bugs

If you've found a bug:

1. **Verify it's a bug**: Ensure it's not a configuration issue
2. **Check existing issues**: Search for duplicates
3. **Create a bug report**: Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.md)
4. **Provide details**: Include version, OS, steps to reproduce

[Report a Bug](https://github.com/sazardev/goca/issues/new?template=bug_report.md)

### Request Features

For feature requests and enhancements:

1. **Search existing requests**: Check if it's already been suggested
2. **Consider the scope**: Ensure it aligns with project goals
3. **Create a feature request**: Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.md)
4. **Explain the use case**: Help us understand the value

[Request a Feature](https://github.com/sazardev/goca/issues/new?template=feature_request.md)

### Security Issues

For security vulnerabilities:

**DO NOT create public issues for security problems.**

- Email: sazardev@gmail.com
- See: [Security Policy](SECURITY.md)
- We take security seriously and will respond promptly

## Response Times

**GitHub Issues and Discussions:**
- Initial response: Within 48-72 hours
- Resolution time: Depends on complexity and priority

**Security Issues:**
- Initial acknowledgment: Within 24-48 hours
- Critical vulnerabilities: Prioritized for immediate action

**Pull Requests:**
- Initial review: Within 3-5 days
- Merge time: Depends on complexity and test results

Please note: Goca is maintained by volunteers. Response times are estimates and may vary.

## Community Guidelines

When seeking support:

1. **Be Respectful**: Treat others with courtesy and respect
2. **Be Clear**: Provide clear, detailed information
3. **Be Patient**: Maintainers and community members are volunteers
4. **Search First**: Check existing resources before asking
5. **Follow Templates**: Use issue templates when reporting bugs or requesting features
6. **Stay On Topic**: Keep discussions relevant to the issue at hand
7. **Follow the Code of Conduct**: Read and follow our [Code of Conduct](CODE_OF_CONDUCT.md)

## Self-Help Resources

### Debug Mode

Enable verbose output for troubleshooting:

```bash
# Use --dry-run to preview actions
goca feature User --fields "name:string" --dry-run

# Check generated file conflicts
goca feature User --fields "name:string"
# (Will warn about existing files)
```

### Common Commands

```bash
# Check version
goca version

# List available templates
goca init myproject --list-templates

# Preview feature generation
goca feature User --fields "name:string" --dry-run

# Generate with backup
goca feature User --fields "name:string" --force --backup
```

### Logging Issues

When reporting issues, include:

1. **Goca version**: Output of `goca version`
2. **Go version**: Output of `go version`
3. **Operating system**: OS and architecture
4. **Full command**: Exact command you ran
5. **Error output**: Complete error message
6. **Project structure**: Relevant parts of your project
7. **Configuration**: Contents of `.goca.yaml` if applicable

## Contributing

Want to help improve Goca?

- **Report Bugs**: Help us identify and fix issues
- **Suggest Features**: Share your ideas for improvements
- **Improve Documentation**: Help make docs clearer
- **Submit Pull Requests**: Contribute code improvements

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## Commercial Support

Currently, Goca is a community-supported open source project. For inquiries about commercial support, training, or custom development, contact:

- Email: sazardev@gmail.com

## Stay Updated

Follow project updates:

- **GitHub Releases**: [Subscribe to releases](https://github.com/sazardev/goca/releases)
- **Star the Repository**: Get notified of important updates
- **Watch Discussions**: Stay informed about community conversations

## Useful Links

- [GitHub Repository](https://github.com/sazardev/goca)
- [Documentation Website](https://sazardev.github.io/goca)
- [Issue Tracker](https://github.com/sazardev/goca/issues)
- [Discussions](https://github.com/sazardev/goca/discussions)
- [Release Notes](https://github.com/sazardev/goca/releases)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)

## Thank You

Thank you for using Goca and being part of our community. Your feedback and contributions help make Goca better for everyone.
