#!/bin/bash

# Goca Release Script
# Usage: ./scripts/release.sh [version]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get version from argument or prompt
if [ -z "$1" ]; then
    echo -e "${YELLOW}Enter new version (e.g., 1.0.1):${NC}"
    read -r NEW_VERSION
else
    NEW_VERSION="$1"
fi

# Validate version format
if [[ ! "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format. Use semantic versioning (e.g., 1.0.1)"
    exit 1
fi

TAG_VERSION="v$NEW_VERSION"

print_status "Preparing release for version $TAG_VERSION"

# Check if we're on the main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "You're not on the master/main branch. Current branch: $CURRENT_BRANCH"
    echo "Do you want to continue? (y/N)"
    read -r CONTINUE
    if [ "$CONTINUE" != "y" ] && [ "$CONTINUE" != "Y" ]; then
        print_error "Release cancelled"
        exit 1
    fi
fi

# Check if tag already exists
if git tag -l | grep -q "^$TAG_VERSION$"; then
    print_error "Tag $TAG_VERSION already exists"
    exit 1
fi

# Update version in version.go
print_status "Updating version in cmd/version.go"
sed -i.bak "s/Version.*=.*/Version   = \"$NEW_VERSION\"/" cmd/version.go
sed -i.bak "s/BuildTime.*=.*/BuildTime = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"/" cmd/version.go
rm cmd/version.go.bak

# Update CHANGELOG.md
print_status "Updating CHANGELOG.md"
CURRENT_DATE=$(date +%Y-%m-%d)
sed -i.bak "s/## \[Unreleased\]/## [Unreleased]\n\n## [$NEW_VERSION] - $CURRENT_DATE/" CHANGELOG.md
rm CHANGELOG.md.bak

# Build and test
print_status "Running tests"
go test -v ./...

print_status "Building application"
go build -o goca .

# Test CLI
print_status "Testing CLI functionality"
./goca version

# Git operations
print_status "Committing changes"
git add cmd/version.go CHANGELOG.md
git commit -m "chore: bump version to $NEW_VERSION

- Updated version in cmd/version.go
- Updated CHANGELOG.md with release date
"

print_status "Creating tag $TAG_VERSION"
git tag -a "$TAG_VERSION" -m "Release $TAG_VERSION

See CHANGELOG.md for details.
"

print_status "Pushing changes and tag"
git push origin "$CURRENT_BRANCH"
git push origin "$TAG_VERSION"

print_success "Release $TAG_VERSION created successfully!"
print_status "GitHub Actions will now build and publish the release automatically."
print_status "Check the progress at: https://github.com/sazardev/goca/actions"

# Clean up
rm -f goca
rm -f goca.exe

print_success "✅ Release process completed!"
echo ""
echo "🎉 Your release is ready! Here's what happens next:"
echo "1. GitHub Actions will build binaries for all platforms"
echo "2. A new release will be published automatically"
echo "3. Users can install via: go install github.com/sazardev/goca@$TAG_VERSION"
