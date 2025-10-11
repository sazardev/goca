# goca version

Display version and build information.

## Syntax

```bash
goca version
```

## Description

Shows detailed version information about your Goca installation.

## Example Output

```
Goca v2.0.0
Build: 2025-10-11T10:00:00Z
Go Version: go1.24.5
OS/Arch: linux/amd64
Commit: abc123def
```

## Information Displayed

- **Version**: Current Goca version
- **Build**: Build timestamp
- **Go Version**: Go compiler version used
- **OS/Arch**: Operating system and architecture
- **Commit**: Git commit hash (if available)

## Usage

Check your current version:

```bash
goca version
```

Verify you have the latest version:

```bash
# Current version
goca version

# Latest available
# Check: https://github.com/sazardev/goca/releases/latest
```

## Updating

If you installed via `go install`:

```bash
go install github.com/sazardev/goca@latest
goca version
```

## See Also

- [Installation Guide](/guide/installation) - How to install Goca
- [GitHub Releases](https://github.com/sazardev/goca/releases) - All versions
