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
      <a href="/goca/blog/releases/v1-17-1">v1.17.1</a>
      <span class="badge badge-latest">Latest</span>
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
