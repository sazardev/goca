#!/usr/bin/env bash
# Builds one ~8s vertical video per Goca CLI command (terminal typing + card).
# Processes commands one at a time so PNG frame cache never exceeds ~1 video's
# worth of disk space.
# Usage: ./render-command-videos.sh [fps]
set -euo pipefail
cd "$(dirname "$0")"

FPS="${1:-30}"
SVG_ROOT="video/frames-svg"
PNG_ROOT="video/frames-png"
OUT_DIR="video/out/commands"

echo "==> Generating all per-command frames"
node scripts/generate-command-videos.mjs

mkdir -p "$OUT_DIR"

for cmd_dir in "$SVG_ROOT"/*/; do
  name="$(basename "$cmd_dir")"
  png_dir="$PNG_ROOT/$name"
  echo "==> [$name] rasterizing"
  mkdir -p "$png_dir"
  rm -f "$png_dir"/*.png
  ls "$cmd_dir" | xargs -P "$(nproc)" -I{} sh -c \
    'rsvg-convert -z 2 --background-color=none -o "'"$png_dir"'/${1%.svg}.png" "'"$cmd_dir"'/$1"' _ {}

  echo "==> [$name] encoding"
  ffmpeg -y -loglevel error -framerate "$FPS" -i "$png_dir/frame_%05d.png" \
    -vf "scale=1080:1920:flags=lanczos" \
    -c:v libx264 -profile:v high -pix_fmt yuv420p -crf 18 -preset medium \
    -movflags +faststart \
    "$OUT_DIR/$name.mp4"

  rm -rf "$png_dir"
  echo "==> [$name] done -> $OUT_DIR/$name.mp4"
done

echo "==> All command videos built in $OUT_DIR"
