<!-- PROJECT SHIELDS -->

[![Contributors](https://img.shields.io/github/contributors/sazardev/goca.svg?style=flat)](https://github.com/sazardev/goca/graphs/contributors)
[![Forks](https://img.shields.io/github/forks/sazardev/goca.svg?style=flat)](https://github.com/sazardev/goca/network/members)
[![Stargazers](https://img.shields.io/github/stars/sazardev/goca.svg?style=flat)](https://github.com/sazardev/goca/stargazers)
[![Issues](https://img.shields.io/github/issues/sazardev/goca.svg?style=flat)](https://github.com/sazardev/goca/issues)
[![License](https://img.shields.io/github/license/sazardev/goca.svg?style=flat)](https://github.com/sazardev/goca/blob/master/LICENSE)

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h1>Goca</h1>
  <h3>Go Clean Architecture Code Generator</h3>
  <p>
    Production-ready code generation for Clean Architecture projects
    <br />
    <a href="https://sazardev.github.io/goca"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://sazardev.github.io/goca/tutorials/complete-tutorial">View Tutorial</a>
    ·
    <a href="https://github.com/sazardev/goca/issues">Report Bug</a>
    ·
    <a href="https://github.com/sazardev/goca/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about">About The Project</a></li>
    <li><a href="#key-features">Key Features</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#documentation">Documentation</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About

Goca is a professional CLI tool that generates production-ready Go applications following Clean Architecture principles. It eliminates boilerplate code and enforces best practices, allowing developers to focus on business logic.

### Built With

- [Go](https://golang.org/) - Programming language
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## Key Features

- **Layer-Based Generation**: Generates domain, use case, handler, and repository layers
- **Clean Architecture Enforcement**: Automatically maintains proper dependency direction
- **Multi-Protocol Support**: HTTP, gRPC, and CLI handlers
- **Safety Features**: Dry-run mode, file conflict detection, automatic backups
- **Dependency Management**: Automatic go.mod updates and version verification
- **Project Templates**: Predefined configurations for common architectures
- **Production Ready**: Battle-tested code generation with comprehensive testing

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

Download the latest release from [GitHub Releases](https://github.com/sazardev/goca/releases):

**Linux:**
```bash
wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64
chmod +x goca-linux-amd64
sudo mv goca-linux-amd64 /usr/local/bin/goca
```

**macOS:**
```bash
# Intel Macs
wget https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64
chmod +x goca-darwin-amd64
sudo mv goca-darwin-amd64 /usr/local/bin/goca
```

**Windows:**
```powershell
# Download goca-windows-amd64.exe from releases
# Rename to goca.exe and add to PATH
```

## Usage

### Initialize Project

```bash
goca init myproject --module github.com/username/myproject
cd myproject
```

### Generate Complete Feature

```bash
goca feature Employee --fields "name:string,email:string,role:string"
```

### Generate Individual Components

```bash
# Generate entity only
goca entity User --fields "name:string,email:string"

# Generate use case only
goca usecase User

# Generate handler only
goca handler User --type http
```

## Documentation

- **[Official Documentation](https://sazardev.github.io/goca)** - Complete guides
- **[Getting Started](https://sazardev.github.io/goca/getting-started)** - Quick start guide
- **[Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial)** - Step-by-step tutorial
- **[Commands Reference](https://sazardev.github.io/goca/commands/)** - All commands
- **[Clean Architecture Guide](https://sazardev.github.io/goca/guide/clean-architecture)** - Architecture principles

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

### How to Contribute

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## Roadmap

See the [Roadmap](ROADMAP.md) for planned features and releases.

## Community

- [Discussions](https://github.com/sazardev/goca/discussions) - Ask questions and share ideas
- [Issues](https://github.com/sazardev/goca/issues) - Report bugs and request features

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## Contact

Project Maintainer - [@sazardev](https://github.com/sazardev) - sazardev@gmail.com

Project Link: [https://github.com/sazardev/goca](https://github.com/sazardev/goca)

## Acknowledgments

- Clean Architecture by Robert C. Martin
- The Go community for excellent tools and libraries
- All contributors who help improve Goca

## Support

If you find Goca useful, please consider:
- Starring the repository
- Sharing with others
- Contributing to the project
- Reporting bugs and suggesting features

For detailed support information, see [SUPPORT.md](SUPPORT.md).

---

<div align="center">
  Made with dedication for the Go community
  <br />
  <a href="https://github.com/sazardev/goca">GitHub</a>
  ·
  <a href="https://sazardev.github.io/goca">Documentation</a>
  ·
  <a href="https://github.com/sazardev/goca/releases">Releases</a>
</div>
