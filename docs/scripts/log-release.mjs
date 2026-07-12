#!/usr/bin/env node
// Runs on every release (see .github/workflows/auto-release.yml): stamps the
// root CHANGELOG.md's "[Unreleased]" section with the new version + date,
// and appends that same content to docs/blog/releases/latest.md so the docs
// site's release blog stays current without a human writing a post for
// every tag. Exits without changes if there's nothing under [Unreleased].
//
// Usage: node docs/scripts/log-release.mjs <version> <date:YYYY-MM-DD>

import { readFileSync, writeFileSync } from 'fs'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const repoRoot = resolve(__dirname, '..', '..')

const [, , version, date] = process.argv
if (!version || !date) {
  console.error('Usage: node log-release.mjs <version> <date:YYYY-MM-DD>')
  process.exit(1)
}

const changelogPath = resolve(repoRoot, 'CHANGELOG.md')
const changelog = readFileSync(changelogPath, 'utf-8')

const unreleasedHeader = '## [Unreleased]'
const startIdx = changelog.indexOf(unreleasedHeader)
if (startIdx === -1) {
  console.log('No "## [Unreleased]" section found — nothing to log.')
  process.exit(0)
}

const afterHeader = startIdx + unreleasedHeader.length
const nextHeaderIdx = changelog.indexOf('\n## [', afterHeader)
const sectionEnd = nextHeaderIdx === -1 ? changelog.length : nextHeaderIdx
const body = changelog.slice(afterHeader, sectionEnd).trim()

if (!body) {
  console.log('"[Unreleased]" section is empty — nothing to log.')
  process.exit(0)
}

// 1. Stamp CHANGELOG.md: insert a dated version header right after
//    [Unreleased], leaving that section empty (but present) above it.
const stamped =
  changelog.slice(0, afterHeader) +
  `\n\n## [${version}] - ${date}\n\n${body}\n` +
  changelog.slice(sectionEnd)
writeFileSync(changelogPath, stamped)
console.log(`Stamped CHANGELOG.md: [${version}] - ${date}`)

// 2. Append the same content to docs/blog/releases/latest.md, between the
//    release-log markers (newest entry first).
const blogPath = resolve(__dirname, '..', 'blog', 'releases', 'latest.md')
const blog = readFileSync(blogPath, 'utf-8')

const startMarker = '<!-- release-log:start -->'
const endMarker = '<!-- release-log:end -->'
const mStart = blog.indexOf(startMarker)
const mEnd = blog.indexOf(endMarker)
if (mStart === -1 || mEnd === -1 || mEnd < mStart) {
  console.error(`Could not find release-log markers in ${blogPath}`)
  process.exit(1)
}

const existingLog = blog.slice(mStart + startMarker.length, mEnd).trim()
const newEntry = `## v${version} — ${date}\n\n${body}\n`
const updatedLog = existingLog ? `${newEntry}\n---\n\n${existingLog}` : newEntry

const updatedBlog =
  blog.slice(0, mStart + startMarker.length) +
  '\n\n' + updatedLog + '\n\n' +
  blog.slice(mEnd)
writeFileSync(blogPath, updatedBlog)
console.log(`Appended v${version} to ${blogPath}`)
