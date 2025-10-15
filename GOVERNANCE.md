# Governance

This document outlines the governance model for the Goca project.

## Overview

Goca is an open source project maintained by volunteers. This document describes how decisions are made and how the community works together.

## Roles and Responsibilities

### Project Creator and Lead Maintainer

**sazardev** serves as the project creator and lead maintainer with the following responsibilities:

- Setting the project vision and direction
- Making final decisions on major architectural changes
- Reviewing and approving pull requests
- Managing releases and versioning
- Maintaining project infrastructure
- Enforcing the Code of Conduct

### Core Contributors

Core contributors are individuals who have made significant, sustained contributions to the project. They:

- Have merge access to the repository
- Review pull requests
- Participate in architectural discussions
- Help triage issues
- Guide new contributors

**Becoming a Core Contributor:**

Core contributor status is granted to individuals who:
- Have made multiple substantial contributions over time
- Demonstrate deep understanding of the project
- Show commitment to the project's goals
- Act in accordance with the Code of Conduct

The lead maintainer nominates and approves core contributors.

### Contributors

Anyone can become a contributor by:
- Submitting pull requests
- Reporting bugs
- Suggesting features
- Improving documentation
- Helping others in discussions

All contributions are valued and recognized.

## Decision Making

### Minor Changes

Minor changes (bug fixes, documentation, small enhancements) can be:
- Proposed via pull request
- Reviewed by any core contributor
- Merged with approval from one core contributor

### Major Changes

Major changes (new features, breaking changes, architectural changes) require:
- Discussion in a GitHub issue or discussion
- Design proposal if significant
- Approval from the lead maintainer
- Review from at least one core contributor

### Controversial Decisions

For controversial decisions:
- Open a GitHub discussion
- Allow reasonable time for community input (minimum 7 days)
- Consider all perspectives
- Lead maintainer makes final decision if consensus isn't reached

## Communication

### Primary Channels

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: General questions, ideas, design discussions
- **Pull Requests**: Code changes, reviews
- **Email**: Security issues, private matters (sazardev@gmail.com)

### Community Standards

All communication must follow the [Code of Conduct](CODE_OF_CONDUCT.md):
- Be respectful and professional
- Focus on constructive feedback
- Welcome newcomers
- Assume good intentions

## Release Process

### Version Numbering

Goca follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version: Incompatible API changes
- **MINOR** version: New functionality (backwards compatible)
- **PATCH** version: Bug fixes (backwards compatible)

### Release Cycle

- **Patch releases**: As needed for bug fixes
- **Minor releases**: When significant features are ready
- **Major releases**: When breaking changes are necessary

### Release Checklist

1. Update CHANGELOG.md
2. Update version numbers
3. Run full test suite
4. Create GitHub release
5. Build and upload binaries
6. Update documentation
7. Announce release

## Contributing Process

### For New Contributors

1. Read [CONTRIBUTING.md](CONTRIBUTING.md)
2. Check existing issues and discussions
3. Start with small contributions
4. Ask questions if unsure
5. Follow coding standards
6. Be patient with review process

### For Pull Requests

1. Create feature branch
2. Make changes with tests
3. Update documentation
4. Submit pull request
5. Address review feedback
6. Wait for merge approval

### For Reviewers

1. Review code quality and design
2. Verify tests pass
3. Check documentation updates
4. Provide constructive feedback
5. Approve or request changes
6. Merge when ready

## Code of Conduct Enforcement

The [Code of Conduct](CODE_OF_CONDUCT.md) is enforced by:

1. **Warning**: First violation results in a private warning
2. **Temporary Ban**: Repeated violations result in temporary ban
3. **Permanent Ban**: Serious or continued violations result in permanent ban

Enforcement decisions are made by the lead maintainer. Appeals can be sent to sazardev@gmail.com.

## Conflict Resolution

### Process

1. **Discussion**: Try to resolve through respectful discussion
2. **Mediation**: Core contributors can mediate disputes
3. **Final Decision**: Lead maintainer makes final decision if needed

### Escalation

For serious issues:
- Contact lead maintainer directly: sazardev@gmail.com
- Provide clear documentation of the issue
- Allow time for investigation and response

## Intellectual Property

### Copyright

- Contributors retain copyright to their contributions
- All contributions are licensed under the project's MIT License
- By contributing, you agree to license your work under MIT

### Attribution

- Contributors are recognized in AUTHORS file
- Significant contributions acknowledged in release notes
- GitHub automatically tracks all contributions

## Project Assets

### Repository Access

- **Lead Maintainer**: Full admin access
- **Core Contributors**: Write access
- **Contributors**: Fork and pull request

### Infrastructure

The following are managed by the lead maintainer:
- GitHub repository and settings
- GitHub Pages documentation site
- Release binaries and artifacts
- Domain names (if applicable)

## Changes to Governance

This governance document may be updated by:
- Proposing changes via pull request
- Allowing community discussion period
- Approval by lead maintainer
- Announcement to community

## Succession Planning

In the event the lead maintainer is unable to continue:
- Core contributors will elect a new lead maintainer
- Decision requires consensus among core contributors
- Community will be notified of the transition

## Acknowledgments

This governance model is inspired by successful open source projects and adapted for Goca's needs.

## Questions

For questions about governance, contact:
- GitHub Discussions for public questions
- sazardev@gmail.com for private matters

---

Last Updated: October 15, 2025
