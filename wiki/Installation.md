# Installation

This page will guide you through the different methods to install Goca on your system.

## 📋 Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - For cloning repositories and version control
- **Terminal/PowerShell** - To run commands

## 🚀 Installation Methods

### 1. Installation with go install (Recommended)

This is the fastest method and will always give you the latest stable version:

```bash
go install github.com/sazardev/goca@latest
```

**Verify installation:**
```bash
goca version
```

**Expected output:**
```
Goca v1.0.5
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### 2. Binary Downloads

Download the pre-compiled binary for your operating system from [GitHub Releases](https://github.com/sazardev/goca/releases).

#### For Windows:
```powershell
# Download latest version
Invoke-WebRequest -Uri "https://github.com/sazardev/goca/releases/latest/download/goca-windows-amd64.exe" -OutFile "goca.exe"

# Move to a location in PATH
Move-Item goca.exe C:\Windows\System32\goca.exe
```

#### For Linux:
```bash
# Download latest version
wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64

# Make executable and move to PATH
chmod +x goca-linux-amd64
sudo mv goca-linux-amd64 /usr/local/bin/goca
```

#### For macOS (Intel):
```bash
# Download latest version
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64 -o goca

# Make executable and move to PATH
chmod +x goca
sudo mv goca /usr/local/bin/goca
```

#### For macOS (Apple Silicon):
```bash
# Download latest version
curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-arm64 -o goca

# Make executable and move to PATH
chmod +x goca
sudo mv goca /usr/local/bin/goca
```

### 3. Installation with Homebrew (macOS)

If you have Homebrew installed:

```bash
# Add the tap
brew tap sazardev/tools

# Install goca
brew install goca
```

### 4. Build from Source Code

For developers who want the latest development version:

```bash
# Clone the repository
git clone https://github.com/sazardev/goca.git
cd goca

# Build
go build -o goca

# Install globally (optional)
go install
```

## 🔧 PATH Configuration

If you manually installed the binary, make sure it's in your PATH:

### Windows:
1. Open "System Environment Variables"
2. Click "Environment Variables"
3. In "System Variables", find "Path" and click "Edit"
4. Click "New" and add the path where you saved `goca.exe`

### Linux/macOS:
Add this line to your `~/.bashrc`, `~/.zshrc` or `~/.profile`:

```bash
export PATH=$PATH:/path/where/you/saved/goca
```

Then reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc
```

## ✅ Installation Verification

Once installed, verify that everything works correctly:

```bash
# Check version
goca version

# Show help
goca help

# Test basic command
goca init test-project --module test
```

If you see the version information and help, the installation was successful! 🎉

## 🆙 Updates

### With go install:
```bash
go install github.com/sazardev/goca@latest
```

### With Homebrew:
```bash
brew upgrade goca
```

### With binaries:
Download the new version following the binary installation steps.

## 🐛 Troubleshooting

### Error: "goca: command not found"
- ✅ Verify that Goca is in your PATH
- ✅ Restart your terminal after installation
- ✅ On Windows, make sure to use PowerShell or CMD as administrator

### Error: "permission denied"
```bash
# Linux/macOS - Add execution permissions
chmod +x goca

# Windows - Run as administrator
```

### SSL certificates error (Windows)
```powershell
# Use TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
```

### Old Go version
Goca requires Go 1.21+. Update Go from [golang.org](https://golang.org/dl/).

## 🔄 Uninstallation

### If installed with go install:
```bash
# Find location
which goca  # Linux/macOS
where goca  # Windows

# Remove binary
rm $(which goca)  # Linux/macOS
del (where goca)  # Windows
```

### With Homebrew:
```bash
brew uninstall goca
brew untap sazardev/tools
```

## 📞 Support

If you have installation problems:

1. 🔍 Check [Known Issues](https://github.com/sazardev/goca/issues)
2. 💬 Ask in [GitHub Discussions](https://github.com/sazardev/goca/discussions)
3. 🐛 Report a new [Issue](https://github.com/sazardev/goca/issues/new)

---

**Next step: [Getting Started](Getting-Started) →**
