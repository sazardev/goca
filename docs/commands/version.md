---
layout: doc
title: goca version
titleTemplate: Commands | Goca
description: Display the current Goca version number and build metadata such as commit hash, build date, and Go version.
---

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
Git Commit: abc123def
```

There is no separate `OS/Arch` line, and the commit field is labeled `Git Commit`, not `Commit`.

`Build` and `Git Commit` are only populated when the binary is built with the release `-ldflags` (as the official release binaries are). A binary built via a plain `go install`/`go build` — no ldflags — shows `Build: unknown` and `Git Commit: unknown`; that's expected, not a bug.

## Information Displayed

- **Version**: Current Goca version
- **Build**: Build timestamp (`unknown` without release ldflags)
- **Go Version**: Go compiler version used
- **Git Commit**: Git commit hash (`unknown` without release ldflags)

## Usage

Check your current version:

```bash
goca version
```

Print just the bare version number (useful in scripts):

```bash
goca version --short   # or: goca version -s
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
