# Goca Marketing Assets

SVG sources in `svg/`, rendered PNGs (2x) in `png/`. Regenerate with:

```bash
./generate.sh
```

Requires `rsvg-convert` (`librsvg`) and the Liberation Sans/Mono fonts (Arial/Courier-metric-compatible, used instead of proprietary fonts).

## Brand tokens used

- Primary: `#00ADD8` &middot; Accent: `#00D9FF` &middot; Secondary: `#00758F`
- Dark background gradient: `#0a0a0a` &rarr; `#1a1a2e`
- Logo mark: blue circle with 4 white bars at decreasing opacity (0.9/0.7/0.5/0.3), representing the Clean Architecture layers
- Icons: stroke-based (2px, round caps), matching `docs/public/icons/*.svg`
- Fonts: Liberation Sans (headings/body), Liberation Mono (code/terminal)

## Assets

Dark is the default brand mode. Files ending in `-light` are the same design flipped
for use on white/light backgrounds (blog posts, light-mode docs, printed material).
The terminal demo (06) is dark-only &mdash; terminal screenshots read as dark regardless
of surrounding theme.

| File | Size | Use |
|---|---|---|
| `01-social-banner[-light]` | 1200&times;630 | OG/Twitter/LinkedIn card |
| `02-card-layers[-light]` | 1080&times;1080 | Feature post &mdash; layered architecture |
| `03-card-speed[-light]` | 1080&times;1080 | Feature post &mdash; generation speed |
| `04-card-antipatterns[-light]` | 1080&times;1080 | Feature post &mdash; anti-pattern prevention |
| `05-card-testing[-light]` | 1080&times;1080 | Feature post &mdash; generated tests |
| `06-terminal-demo` | 1200&times;750 | Product demo screenshot |
| `07-before-after[-light]` | 1200&times;800 | Without/With Goca comparison |
| `08-readme-banner[-light]` | 1280&times;400 | Wide banner for README/GitHub |
| `09-story-banner` | 1080&times;1920 | Instagram/Facebook Story &mdash; brand intro + CTA |
| `10`&ndash;`13-story-card-*` | 1080&times;1920 | Instagram/Facebook Story, one per feature |

Story assets (09&ndash;13) leave extra top/bottom margin (~170&ndash;250px) clear of text,
matching IG/FB's own profile-header and reply-bar overlay zones.

## Per-command cards (`svg/commands/`, `png/commands/`)

One card per real Goca CLI command (all 21, verified against `cmd/*.go`), in
**3 formats &times; 2 themes = 6 files/command, 126 total**:

| Format | Size | Use |
|---|---|---|
| `square` | 1080&times;1080 | Feed posts (IG/X/LinkedIn) |
| `landscape` | 1200&times;630 | OG cards, blog headers, slides |
| `portrait` | 1080&times;1920 | Stories/Reels/Shorts covers |

Naming: `svg/commands/<command>-<format>-<theme>.svg`, e.g. `entity-square-dark.svg`.

These are **generated, not hand-authored** &mdash; edit `scripts/generate-commands.mjs`
(the `COMMANDS` array holds copy/tag/icon/example per command; `ICONS` holds the
stroke-icon library; layout math lives in `renderSquare/Landscape/Portrait`) and
rerun `./generate.sh`. Do not hand-edit files under `svg/commands/` &mdash; they're
overwritten on every run.

Commands covered: `init`, `feature`, `entity`, `usecase`, `repository`, `handler`,
`di`, `interfaces`, `messages`, `integrate`, `middleware`, `mocks`,
`test-integration`, `ci`, `doctor`, `analyze`, `config`, `template`, `upgrade`,
`mcp-server`, `version`.

## Video (`video/`)

A ~23s vertical explainer, ready for TikTok/Reels/Shorts (1080&times;1920, 30fps,
H.264, ~2MB): `video/out/goca-explainer.mp4`. Build it with:

```bash
./render-video.sh          # 30fps (default)
./render-video.sh 60       # or any other fps
```

This runs `scripts/generate-video-frames.mjs` (writes one animated SVG frame per
video frame to `video/frames-svg/`), rasterizes every frame to PNG in parallel
with `rsvg-convert` (`video/frames-png/`), then encodes with `ffmpeg`. Both frame
directories are build cache &mdash; gitignored, regenerated on every run, safe to
delete.

**Storyboard** (all timings in `generate-video-frames.mjs`, easily tweakable):

| Scene | Time | What happens |
|---|---|---|
| Intro | 0.0&ndash;2.8s | Logo mark scales in, wordmark + tagline fade up |
| Terminal | 2.8&ndash;10.0s | `goca feature Order ...` types out character-by-character with a blinking cursor, then a 5-line checklist reveals staggered, ending in a green "ready in 1.4s" |
| Layers | 10.0&ndash;14.5s | The 4 Clean Architecture layers build bottom-up as labeled bars (Domain &rarr; Use Case &rarr; Repository &rarr; Handler) |
| Value cards | 14.5&ndash;19.5s | 3 crossfading cards: anti-patterns blocked, tests generated, wired together |
| Outro | 19.5&ndash;23.0s | Logo + tagline, "Get Goca free" CTA, install command, URL |

Scenes crossfade into each other (0.35s) and the whole video fades from/to black
at the very start/end. To change copy, timing, or add scenes, edit the `SCENES`
array and the per-scene `render*` functions &mdash; they reuse the same brand
tokens/icons from `scripts/brand.mjs` as the static cards, so the video stays
visually consistent with the rest of this folder automatically.

### Per-command videos (`video/out/commands/`)

All 21 commands, same visual system, one ~7&ndash;8.6s vertical clip each (1080&times;1920,
30fps, H.264, ~0.5MB each, ~11MB total):

```bash
./render-command-videos.sh        # 30fps (default)
./render-command-videos.sh 60
```

Each clip is generated by `scripts/generate-command-videos.mjs` from the same
`commands-data.mjs` used by the static cards (one source of truth for copy), in
two scenes:

1. **Terminal** &mdash; the real `goca <command> ...` invocation types out (typing
   speed scales with command length, clamped 1.1&ndash;2.6s), then the command's tag
   and two-line explanation fade in as code comments, ending in a green
   "&#10003; Ready." held on screen for at least 0.5s before the cut &mdash; so short
   commands (e.g. `doctor`) get a snappier ~7.1s clip and long ones (e.g.
   `usecase`) stretch to ~8.6s rather than rushing or padding.
2. **Card** &mdash; icon + two-line headline (same copy as the static cards) with a
   command chip (truncated with `…` past 46 chars) and the URL footer.

`render-command-videos.sh` processes one command at a time (generate &rarr;
rasterize &rarr; encode &rarr; delete that command's PNG cache) so disk usage never
exceeds roughly one video's worth of frames, even though there are 21 of them.
