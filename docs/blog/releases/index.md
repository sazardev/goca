---
layout: doc
title: Releases
titleTemplate: Blog | Goca
---

# Release Notes

Track the evolution of Goca through detailed release notes. Each release includes new features, bug fixes, improvements, and migration guides.

---

<style scoped>
.release-list {
  margin-top: 2rem;
}

.release-item {
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  margin-bottom: 1.5rem;
  transition: all 0.3s ease;
}

.release-item:hover {
  border-color: var(--vp-c-brand);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.release-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 0.5rem;
}

.release-version {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--vp-c-brand);
}

.release-date {
  color: var(--vp-c-text-2);
  font-size: 0.9rem;
}

.release-description {
  color: var(--vp-c-text-2);
  margin-top: 0.5rem;
  line-height: 1.6;
}

.release-highlights {
  margin-top: 1rem;
}

.release-highlights ul {
  margin-top: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
  margin-left: 0.5rem;
}

.badge-latest {
  background: var(--vp-c-brand-soft);
  color: var(--vp-c-brand);
  border: 1px solid var(--vp-c-brand);
}

.badge-major {
  background: var(--vp-c-green-soft);
  color: var(--vp-c-green);
  border: 1px solid var(--vp-c-green);
}
</style>

<div class="release-list">

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-18-7">v1.18.7</a>
      <span class="badge badge-latest">Latest</span>
    </h2>
    <span class="release-date">March 24, 2026</span>
  </div>
  <p class="release-description">
    Universal Safety Flags - <code>--dry-run</code>, <code>--force</code>, and <code>--backup</code> now correctly wired through all 12 file-generating commands. SafetyManager fully connected to every sub-generator.
  </p>
  <div class="release-highlights">
    <strong>Key Changes:</strong>
    <ul>
      <li>Fixed <code>goca init --dry-run</code> failing with "unknown flag" error</li>
      <li><code>--dry-run</code>, <code>--force</code>, <code>--backup</code> registered on all 12 commands: entity, usecase, repository, handler, di, messages, interfaces, mocks, init, integrate, feature, test-integration</li>
      <li><code>writeFile()</code> and <code>writeGoFile()</code> now accept variadic <code>*SafetyManager</code> — backward-compatible</li>
      <li>SafetyManager forwarded from <code>feature</code> into all 8 sub-generators</li>
      <li>SafetyManager threaded through <code>integrate</code> DI and main.go generation paths</li>
      <li>15 files modified across <code>cmd/</code> package</li>
    </ul>
  </div>
</div>

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-18-2">v1.18.1 — v1.18.6</a>
    </h2>
    <span class="release-date">March 24–25, 2026</span>
  </div>
  <p class="release-description">
    Patch series: interactive wizard fix, version display fix, <code>goca doctor</code>, <code>goca upgrade</code>, global <code>--quiet</code>/<code>--verbose</code> flags, improved dry-run table, and CI fixes.
  </p>
  <div class="release-highlights">
    <strong>Key Changes:</strong>
    <ul>
      <li><a href="/goca/blog/releases/v1-18-1">v1.18.1</a>: <code>goca init</code> no-arg wizard fix; <code>goca version</code> reads correct installed version via <code>runtime/debug.ReadBuildInfo</code></li>
      <li><a href="/goca/blog/releases/v1-18-2">v1.18.2</a>: New <code>goca doctor</code> command (6 health checks, <code>--fix</code>); new <code>goca upgrade</code> command; <code>--quiet</code>/<code>--verbose</code> global flags; improved 3-column dry-run table</li>
      <li>v1.18.3: Resolved CI failures from v1.18.2</li>
      <li>v1.18.4: Removed conflicting <code>-v</code> shorthand from <code>--validation</code> flag</li>
      <li>v1.18.5: Restored missing <code>--business-rules</code> and <code>--dry-run</code> flags in <code>feature</code> command</li>
      <li>v1.18.6: YAML workflow formatting fix</li>
    </ul>
  </div>
</div>

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-18-0">v1.18.0</a>
    </h2>
    <span class="release-date">March 24, 2026</span>
  </div>
  <p class="release-description">
    CLI UX Overhaul - Centralized output rendering system, interactive initialization wizard, global --no-color and --no-interactive flags, unified English output, and multiple bug fixes.
  </p>
  <div class="release-highlights">
    <strong>Key Changes:</strong>
    <ul>
      <li>New UIRenderer system replaces all raw fmt.Printf calls across every command</li>
      <li>Interactive wizard for goca init using huh forms</li>
      <li>Global --no-color and --no-interactive flags on every command</li>
      <li>Table output for generated file summaries in entity, feature, and di commands</li>
      <li>Spinner animations for long-running operations</li>
      <li>goca handler now auto-injects go-playground/validator/v10 when validation is enabled</li>
      <li>Removed duplicate configCmd registration and spurious DEBUG: output</li>
    </ul>
  </div>
</div>

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-17-2">v1.17.2</a>
    </h2>
    <span class="release-date">February 1, 2026</span>
  </div>
  <p class="release-description">
    Database Defaults and MongoDB Fixes - Changed default database to SQLite and fixed MongoDB code generation issues.
  </p>
  <div class="release-highlights">
    <strong>Key Changes:</strong>
    <ul>
      <li>🎯 SQLite is now the default database (was PostgreSQL)</li>
      <li>✅ Fixed MongoDB code generation to use mongo-driver correctly</li>
      <li>✅ MongoDB projects no longer import GORM incorrectly</li>
      <li>🧪 Added comprehensive database initialization tests</li>
    </ul>
  </div>
</div>

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-17-1">v1.17.1</a>
    </h2>
    <span class="release-date">January 11, 2026</span>
  </div>
  <p class="release-description">
    Critical Bug Fix - SQLite and other database drivers now properly configured during project initialization. Resolves issue #31.
  </p>
  <div class="release-highlights">
    <strong>Bug Fixes:</strong>
    <ul>
      <li>✅ Fixed database driver configuration during <code>goca init</code></li>
      <li>✅ SQLite, MySQL, SQL Server, MongoDB, DynamoDB, and Elasticsearch now generate correct dependencies</li>
      <li>✅ Main.go now imports and uses correct database driver package</li>
      <li>✅ All 8 supported database types verified working</li>
    </ul>
  </div>
</div>

<div class="release-item">
  <div class="release-header">
    <h2 class="release-version">
      <a href="/goca/blog/releases/v1-14-1">v1.14.1</a>
    </h2>
    <span class="release-date">October 27, 2025</span>
  </div>
  <p class="release-description">
    Test Suite Improvements - Fixed Windows path handling, test working directory management, and module dependencies. Test success rate improved to 99.04%.
  </p>
  <div class="release-highlights">
    <strong>Key Improvements:</strong>
    <ul>
      <li>Fixed Windows path handling in BackupFile</li>
      <li>Improved test working directory management</li>
      <li>Updated test message validation</li>
      <li>Fixed module dependencies for testify</li>
    </ul>
  </div>
</div>

</div>

---

## Release Versioning

Goca follows [Semantic Versioning](https://semver.org/):

- **Major (X.0.0)**: Breaking changes
- **Minor (1.X.0)**: New features (backward compatible)
- **Patch (1.14.X)**: Bug fixes and minor improvements

## Stay Updated

- Watch the [GitHub repository](https://github.com/sazardev/goca) for release notifications
- View the complete [CHANGELOG](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)
- Subscribe to [GitHub Releases](https://github.com/sazardev/goca/releases)
