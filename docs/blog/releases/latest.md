---
layout: doc
title: Latest Changes
titleTemplate: Releases | Goca Blog
description: Rolling log of changes since the last named release, appended automatically by CI on every release.
---

# Latest Changes

Tracks changes since [v1.22.1 — v1.25.15](/blog/releases/v1-25-15). Each release appends its `CHANGELOG.md` entry here automatically — see [scripts/log-release.mjs](https://github.com/sazardev/goca/blob/master/docs/scripts/log-release.mjs). When this page accumulates enough for a proper write-up, its content gets consolidated into a new named post under [Releases](/blog/releases/) and this page resets.

<!-- release-log:start -->

## v1.25.16 — 2026-07-12

### Fixed
- **handler**: an entity name starting with "W" or "R" (Widget, Wallet, Report, ...) broke every generated HTTP handler method — the receiver variable, derived from the entity's first letter, collided with the fixed `w http.ResponseWriter`/`r *http.Request` parameter names ("w redeclared in this block"). Added `handlerReceiverVar`, which falls back to `h` for those two letters.


<!-- release-log:end -->
