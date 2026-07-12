#!/usr/bin/env bash
# Cuts a new goca release: bumps the version, updates CHANGELOG.md, tags and
# pushes. CI (.github/workflows/release.yml) picks up the pushed vX.Y.Z tag
# and runs goreleaser.
#
# Usage:
#   scripts/release.sh patch|minor|major   # bump the last release tag
#   scripts/release.sh auto                # infer patch/minor/major from
#                                           # Conventional Commit prefixes
#                                           # since the last tag
#   scripts/release.sh 1.4.0               # use an explicit version
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

if [[ -n "$(git status --porcelain)" ]]; then
	echo "Error: working tree is not clean. Commit or stash changes before releasing." >&2
	exit 1
fi

mode="${1:-}"
if [[ -z "$mode" ]]; then
	echo "Usage: $0 <patch|minor|major|auto|X.Y.Z>" >&2
	exit 1
fi

last_tag="$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")"
current_version="${last_tag#v}"
IFS='.' read -r major minor patch <<<"$current_version"

case "$mode" in
patch | minor | major)
	bump="$mode"
	;;
auto)
	commits="$(git log "${last_tag}..HEAD" --pretty=format:'%s' 2>/dev/null || true)"
	if echo "$commits" | grep -qE '^[a-zA-Z]+(\([^)]*\))?!:|BREAKING CHANGE:'; then
		bump="major"
	elif echo "$commits" | grep -qE '^feat(\([^)]*\))?:'; then
		bump="minor"
	else
		bump="patch"
	fi
	echo "Auto-detected bump: $bump"
	;;
[0-9]*.[0-9]*.[0-9]*)
	bump="explicit"
	new_version="${mode#v}"
	;;
*)
	echo "Error: unrecognized release type or version '$mode'" >&2
	echo "Usage: $0 <patch|minor|major|auto|X.Y.Z>" >&2
	exit 1
	;;
esac

if [[ "$bump" != "explicit" ]]; then
	case "$bump" in
	patch) new_version="${major}.${minor}.$((patch + 1))" ;;
	minor) new_version="${major}.$((minor + 1)).0" ;;
	major) new_version="$((major + 1)).0.0" ;;
	esac
fi

new_tag="v${new_version}"

if git rev-parse "$new_tag" >/dev/null 2>&1; then
	echo "Error: tag $new_tag already exists" >&2
	exit 1
fi

echo "Releasing $last_tag -> $new_tag"

changelog="CHANGELOG.md"
if [[ -f "$changelog" ]] && grep -q '^## \[Unreleased\]' "$changelog"; then
	release_date="$(date +%Y-%m-%d)"
	# Insert a dated release header right after [Unreleased], leaving that
	# section empty (but present) for the next round of changes.
	awk -v date="$release_date" -v ver="$new_version" '
		/^## \[Unreleased\]/ && !done {
			print
			print ""
			print "## [" ver "] - " date
			done = 1
			next
		}
		{ print }
	' "$changelog" >"$changelog.tmp"
	mv "$changelog.tmp" "$changelog"
	git add "$changelog"
	git commit -m "chore: release $new_tag"
fi

git tag -a "$new_tag" -m "Release $new_tag"
git push origin HEAD
git push origin "$new_tag"

echo "Tagged and pushed $new_tag. CI will build and publish the release."
