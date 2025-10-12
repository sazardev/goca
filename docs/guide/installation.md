# Installation

This guide will walk you through installing Goca on your system using various methods.

## Prerequisites

Before installing Goca, ensure you have:

- **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
- **Git** - For version control and cloning repositories
- **Terminal or Command Prompt** - To run installation commands

::: tip Check Your Go Version
```bash
go version
```
You should see `go version go1.21` or higher.
:::

## Installation Methods

### Method 1: go install (Recommended)

This is the fastest and simplest method:

```bash
go install github.com/sazardev/goca@latest
```

**Verify the installation:**

```bash
goca version
```

**Expected output:**

```
Goca v2.0.0
Build: 2025-10-11T10:00:00Z
Go Version: go1.24.5
OS/Arch: linux/amd64
```

::: details Troubleshooting: Command Not Found
If you get `command not found`, ensure `$GOPATH/bin` is in your PATH:

**Linux/macOS:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Add this to your `~/.bashrc`, `~/.zshrc`, or `~/.profile` to make it permanent.

**Windows:**
Add `%USERPROFILE%\go\bin` to your system PATH environment variable.
:::

### Method 2: Binary Downloads

Download pre-compiled binaries directly from [GitHub Releases](https://github.com/sazardev/goca/releases).

::: code-group

```bash [Linux]
# Download the binary
wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64

# Make it executable
chmod +x goca-linux-amd64

# Move to PATH
sudo mv goca-linux-amd64 /usr/local/bin/goca

# Verify
goca version
```

```bash [macOS (Intel)]
# Download the binary
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64 -o goca

# Make it executable
chmod +x goca

# Move to PATH
sudo mv goca /usr/local/bin/goca

# Verify
goca version
```

```bash [macOS (Apple Silicon)]
# Download the binary
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-arm64 -o goca

# Make it executable
chmod +x goca

# Move to PATH
sudo mv goca /usr/local/bin/goca

# Verify
goca version
```

```powershell [Windows]
# Download the binary
Invoke-WebRequest -Uri "https://github.com/sazardev/goca/releases/latest/download/goca-windows-amd64.exe" -OutFile "goca.exe"

# Move to a directory in PATH (requires admin)
Move-Item goca.exe C:\Windows\System32\goca.exe

# Verify
goca version
```

:::

### Method 3: Homebrew (macOS)

If you're on macOS and use Homebrew:

```bash
# Add the Goca tap
brew tap sazardev/tools

# Install Goca
brew install goca

# Verify
goca version
```

::: tip Updating via Homebrew
```bash
brew upgrade goca
```
:::

### Method 4: Build from Source

For developers who want the latest development version or want to contribute:

```bash
# Clone the repository
git clone https://github.com/sazardev/goca.git
cd goca

# Build the binary
go build -o goca

# (Optional) Install globally
sudo mv goca /usr/local/bin/goca

# Verify
goca version
```

::: details Building for Different Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o goca-linux-amd64

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o goca-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o goca-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o goca-windows-amd64.exe
```
:::

## Verify Installation

After installation, run:

```bash
goca --help
```

You should see the help menu with all available commands:

```
Goca - Go Clean Architecture Code Generator

Usage:
  goca [command]

Available Commands:
  init        Initialize a new Clean Architecture project
  feature     Generate a complete feature with all layers
  entity      Generate a domain entity
  usecase     Generate use cases
  repository  Generate repositories
  handler     Generate handlers
  di          Generate dependency injection
  integrate   Integrate existing features
  version     Show version information

Flags:
  -h, --help      help for goca
  -v, --version   version for goca

Use "goca [command] --help" for more information about a command.
```

## Shell Completion (Optional)

Enable command auto-completion for your shell:

::: code-group

```bash [Bash]
# Generate completion script
goca completion bash > /etc/bash_completion.d/goca

# Or for current user only
goca completion bash > ~/.bash_completion
source ~/.bash_completion
```

```bash [Zsh]
# Generate completion script
goca completion zsh > "${fpath[1]}/_goca"

# Reload completions
autoload -U compinit && compinit
```

```bash [Fish]
# Generate completion script
goca completion fish > ~/.config/fish/completions/goca.fish
```

```powershell [PowerShell]
# Generate completion script
goca completion powershell | Out-String | Invoke-Expression

# To make permanent, add to profile
goca completion powershell >> $PROFILE
```

:::

## Update Goca

### If installed via go install:

```bash
go install github.com/sazardev/goca@latest
```

### If installed via Homebrew:

```bash
brew upgrade goca
```

### If installed from binary:

Download the latest binary and replace the existing one following the [Binary Downloads](#method-2-binary-downloads) steps.

## Uninstall Goca

### If installed via go install:

```bash
rm $(which goca)
```

### If installed via Homebrew:

```bash
brew uninstall goca
brew untap sazardev/tools
```

### If installed from binary:

```bash
# Linux/macOS
sudo rm /usr/local/bin/goca

# Windows (as Administrator)
del C:\Windows\System32\goca.exe
```

## Next Steps

Now that you have Goca installed, you're ready to start building!

-  [Quick Start Guide](/getting-started) - Create your first project
-  [Learn Clean Architecture](/guide/clean-architecture) - Understand the principles
-  [Complete Tutorial](/tutorials/complete-tutorial) - Build a real application

## Troubleshooting

### Permission Denied

If you get permission errors on Linux/macOS:

```bash
sudo chmod +x /usr/local/bin/goca
```

### Command Not Found After Installation

Make sure your `$PATH` includes Go's bin directory:

```bash
echo $PATH | grep -q "go/bin" && echo " Go bin in PATH" || echo "✗ Add Go bin to PATH"
```

### Version Mismatch

If `goca version` shows an old version:

```bash
# Clear Go cache
go clean -modcache

# Reinstall
go install github.com/sazardev/goca@latest
```

### Windows: goca is not recognized

Ensure you've added Go's bin directory to your system PATH:

1. Open System Properties → Environment Variables
2. Edit the `Path` variable
3. Add `%USERPROFILE%\go\bin`
4. Restart your terminal

## Need Help?

-  [GitHub Issues](https://github.com/sazardev/goca/issues) - Report bugs
-  [Discussions](https://github.com/sazardev/goca/discussions) - Ask questions
-  [Documentation](/) - Read the docs
