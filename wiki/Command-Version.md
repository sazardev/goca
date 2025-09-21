# goca version Command

The `goca version` command provides detailed information about the installed version of Goca, including build metadata and compatibility.

## 📋 Syntax

```bash
goca version [flags]
```

## 🎯 Purpose

Shows complete information about the current Goca installation:

- 🏷️ **Version number** - Current semantic version
- 📅 **Build date** - When this version was compiled
- 🔧 **Go version** - Go version used for compilation
- 📦 **Build information** - Additional build metadata

## 🚩 Available Flags

| Flag             | Type   | Required | Default Value | Description                  |
| ---------------- | ------ | -------- | ------------- | ---------------------------- |
| `--short` / `-s` | `bool` | ❌ No     | `false`       | Show only the version number |

## 📖 Usage Examples

### Complete Information
```bash
goca version
```

**Output:**
```
Goca v1.0.6
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### Short Version
```bash
goca version --short
# or
goca version -s
```

**Output:**
```
1.0.6
```

## 🔍 Detailed Information

### Version Number
- Follows **Semantic Versioning (SemVer)** format
- Format: `MAJOR.MINOR.PATCH`
- Example: `1.0.5` means:
  - **Major (1)**: Incompatible API changes
  - **Minor (0)**: New backward-compatible functionality
  - **Patch (5)**: Backward-compatible bug fixes

### Build Date
- **ISO 8601** format: `YYYY-MM-DDTHH:MM:SSZ`
- Always in **UTC**
- Indicates when the specific binary was compiled

### Go Version
- Shows the exact Go version used
- Important for **compatibility** and **debugging**
- Format: `go1.XX.Y`

## 🛠️ Use Cases

### Verify Installation
```bash
# Check that Goca is installed correctly
goca version
```

### Automation Scripts
```bash
#!/bin/bash

# Get only the version for scripts
VERSION=$(goca version --short)
echo "Using Goca v$VERSION"

# Check minimum required version
REQUIRED="1.0.0"
if [[ "$(printf '%s\n' "$REQUIRED" "$VERSION" | sort -V | head -n1)" != "$REQUIRED" ]]; then
    echo "Error: Goca v$REQUIRED or higher is required"
    exit 1
fi
```

### CI/CD Integration
```yaml
# GitHub Actions
- name: Check Goca Version
  run: |
    goca version
    GOCA_VERSION=$(goca version --short)
    echo "GOCA_VERSION=$GOCA_VERSION" >> $GITHUB_ENV
```

### Debugging
```bash
# Complete information for bug reports
goca version > goca-version.txt
echo "System: $(uname -a)" >> goca-version.txt
echo "Go installed: $(go version)" >> goca-version.txt
```

## 📊 Version Analysis

### Development Versions
```bash
# Development versions may include suffixes
goca version
# Output: Goca v1.1.0-dev
```

### Release Candidate Versions
```bash
# Release candidate versions
goca version
# Output: Goca v1.1.0-rc.1
```

### Stable Versions
```bash
# Final versions without suffixes
goca version
# Output: Goca v1.0.5
```

## 🔄 Compatibility

### Go Compatibility
| Goca Version | Minimum Go | Recommended Go | Notes          |
| ------------ | ---------- | -------------- | -------------- |
| v1.0.x       | Go 1.21    | Go 1.24+       | Stable version |
| v1.1.x       | Go 1.22    | Go 1.24+       | Next version   |

### Feature Compatibility
```bash
# Check if your version supports a feature
goca version

# Compare with feature documentation:
# v1.0.0: Basic functionalities
# v1.0.1: Bug fixes
# v1.0.5: gRPC and validation improvements
```

## 🚀 Updates

### Check for Updates
```bash
# Current version
CURRENT=$(goca version --short)
echo "Current version: v$CURRENT"

# Check latest version on GitHub (requires curl/jq)
LATEST=$(curl -s https://api.github.com/repos/sazardev/goca/releases/latest | jq -r .tag_name)
echo "Latest version: $LATEST"

if [ "v$CURRENT" != "$LATEST" ]; then
    echo "Update available!"
    echo "Run: go install github.com/sazardev/goca@latest"
fi
```

### Update to Latest Version
```bash
# Update using go install
go install github.com/sazardev/goca@latest

# Verify update
goca version
```

### Install Specific Version
```bash
# Install specific version
go install github.com/sazardev/goca@v1.0.6

# Verify installed version
goca version
```

## 🔍 Detailed Build Information

### Build Variables
The `version` command shows information that is compiled at build time:

```go
// Defined in cmd/version.go
var (
	Version   = "1.0.6"                    // Software version
	BuildTime = "2025-07-19T15:00:00Z"     // Compilation timestamp
	GoVersion = runtime.Version()          // Go runtime version
)
```

### Build Tags and Flags
```bash
# Build information (if available)
goca version --verbose  # (if implemented in future versions)
```

## 📝 Output Format

### Normal Format
```
Goca v1.0.6
Build: 2025-07-19T15:00:00Z
Go Version: go1.24.5
```

### Short Format
```
1.0.6
```

### JSON Format (Future)
```bash
# Possible future implementation
goca version --json
```

```json
{
  "version": "1.0.6",
  "build_time": "2025-07-19T15:00:00Z",
  "go_version": "go1.24.5",
  "git_commit": "abc123def",
  "build_user": "github-actions"
}
```

## 🐛 Troubleshooting

### Command Not Found
```bash
# Error: command not found
which goca          # Linux/macOS
where goca          # Windows

# Check PATH
echo $PATH          # Linux/macOS
echo $env:PATH      # PowerShell
```

### Old Version
```bash
# Check multiple installations
which -a goca       # Linux/macOS

# Clear Go cache
go clean -modcache

# Reinstall
go install github.com/sazardev/goca@latest
```

### Inconsistent Information
```bash
# Check integrity
goca version

# Compare with project's go.mod file
cat go.mod | grep goca

# Check on GitHub
curl -s https://api.github.com/repos/sazardev/goca/releases/latest
```

## 📞 Support and Reports

### Include in Bug Reports
Always include the output of `goca version` in bug reports:

```bash
# Information for reports
echo "=== GOCA VERSION INFO ===" > bug-report.txt
goca version >> bug-report.txt
echo "=== SYSTEM INFO ===" >> bug-report.txt
uname -a >> bug-report.txt
go version >> bug-report.txt
```

### Useful Links
- 🐛 **Issues**: [GitHub Issues](https://github.com/sazardev/goca/issues)
- 📋 **Releases**: [GitHub Releases](https://github.com/sazardev/goca/releases)
- 📖 **Changelog**: [CHANGELOG.md](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)

## 🔄 Version History

### Important Versions

#### v1.0.6 (Current)
- ✅ Critical bugs fixed in entity, interfaces and di
- ✅ --features flag added to di command
- ✅ --fields flag added to entity command  
- ✅ Flag conflict fixed in interfaces

#### v1.0.5 (Previous)
- ✅ gRPC generation improvements
- ✅ Enhanced validations
- ✅ Bug fixes

#### v1.0.0 (Initial Release)
- 🎉 Initial release
- ✅ Basic Clean Architecture functionalities
- ✅ Multi-database support
- ✅ HTTP and gRPC handlers

### Upcoming Versions
- 🔮 **v1.1.0**: Microservices support
- 🔮 **v1.2.0**: Customizable templates
- 🔮 **v2.0.0**: Rewrite with performance improvements

---

**← [goca messages Command](Command-Messages) | [Getting Started](Getting-Started) →**
