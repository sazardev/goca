# Scripts Directory

This directory contains utility scripts for managing the Goca project.

## Available Scripts

### `release.sh` (Linux/macOS)
Create a new release with automatic version management.

**Usage:**
```bash
./scripts/release.sh [version]
```

**Example:**
```bash
./scripts/release.sh 1.0.1
```

**What it does:**
1. Updates version in `cmd/version.go`
2. Updates `CHANGELOG.md` with release date
3. Runs tests to ensure everything works
4. Creates git commit and tag
5. Pushes to GitHub
6. Triggers automated release build

### `release.bat` (Windows)
Windows version of the release script.

**Usage:**
```cmd
scripts\release.bat [version]
```

**Example:**
```cmd
scripts\release.bat 1.0.1
```

## Alternative Release Methods

### Using Make (Linux/macOS)
```bash
# Check if ready for release
make pre-release-check

# Create release
make release VERSION=1.0.1
```

### Manual Release Process
If you prefer to do it manually:

1. **Update version:**
   ```bash
   # Edit cmd/version.go
   Version = "1.0.1"
   BuildTime = "2025-01-19T10:30:00Z"
   ```

2. **Update CHANGELOG.md:**
   ```bash
   # Add release date to version section
   ## [1.0.1] - 2025-01-19
   ```

3. **Test and commit:**
   ```bash
   go test ./...
   git add cmd/version.go CHANGELOG.md
   git commit -m "chore: bump version to 1.0.1"
   ```

4. **Create tag:**
   ```bash
   git tag -a v1.0.1 -m "Release v1.0.1"
   git push origin master
   git push origin v1.0.1
   ```

### GitHub Actions Release
You can also trigger a release manually from GitHub:

1. Go to `Actions` tab in your repository
2. Select `Release` workflow
3. Click `Run workflow`
4. Enter the version (e.g., `v1.0.1`)
5. Click `Run workflow`

## GitHub Actions Workflows

### Automatic Release Trigger
- **File:** `.github/workflows/release.yml`
- **Trigger:** When you push a tag starting with `v` (e.g., `v1.0.1`)
- **What it does:** Builds binaries for all platforms and creates GitHub release

### Auto Release on Version Change
- **File:** `.github/workflows/auto-release.yml`  
- **Trigger:** When you push to master and version in `cmd/version.go` changed
- **What it does:** Automatically creates tag and triggers release

### Testing
- **File:** `.github/workflows/test.yml`
- **Trigger:** On every push and pull request
- **What it does:** Runs tests and CLI functionality checks

## Release Checklist

Before creating a release:

- [ ] All tests pass: `go test ./...`
- [ ] CLI builds successfully: `go build .`
- [ ] CLI functionality works: `./goca version && ./goca help`
- [ ] Documentation is updated
- [ ] CHANGELOG.md has unreleased changes
- [ ] Version number follows semantic versioning

After release:

- [ ] Verify release appears on GitHub
- [ ] Test installation: `go install github.com/sazardev/goca@v1.0.1`
- [ ] Verify binaries work on different platforms
- [ ] Update any external documentation if needed
