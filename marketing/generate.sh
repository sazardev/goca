#!/usr/bin/env bash
# Regenerates every marketing asset:
#   1. Rebuilds the per-command SVGs from scripts/generate-commands.mjs
#   2. Renders every SVG (hand-authored + generated) to a 2x PNG
# Requires: node, rsvg-convert (librsvg), Liberation Sans/Mono fonts.
set -euo pipefail
cd "$(dirname "$0")"

node scripts/generate-commands.mjs

render_dir() {
  local svg_dir="$1" png_dir="$2"
  mkdir -p "$png_dir"
  for f in "$svg_dir"/*.svg; do
    name="$(basename "${f%.svg}")"
    rsvg-convert -z 2 --background-color=none -o "${png_dir}/${name}.png" "$f"
    echo "OK: ${png_dir}/${name}.png"
  done
}

render_dir svg png
render_dir svg/commands png/commands
