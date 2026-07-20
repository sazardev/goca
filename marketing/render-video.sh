#!/usr/bin/env bash
# Builds the "Goca in ~23 seconds" vertical explainer video:
#   1. Generates animated SVG frames (scripts/generate-video-frames.mjs)
#   2. Rasterizes each frame to PNG in parallel (rsvg-convert)
#   3. Encodes to H.264 MP4 at 1080x1920 (ffmpeg)
# Usage: ./render-video.sh [fps]
set -euo pipefail
cd "$(dirname "$0")"

FPS="${1:-30}"
SVG_DIR="video/frames-svg"
PNG_DIR="video/frames-png"
OUT="video/out/goca-explainer.mp4"

echo "==> Generating frames"
node scripts/generate-video-frames.mjs

echo "==> Clearing old PNG frames"
mkdir -p "$PNG_DIR"
rm -f "$PNG_DIR"/*.png

echo "==> Rasterizing SVG -> PNG (parallel, x$(nproc))"
ls "$SVG_DIR" | xargs -P "$(nproc)" -I{} sh -c \
  'rsvg-convert -z 2 --background-color=none -o "'"$PNG_DIR"'/${1%.svg}.png" "'"$SVG_DIR"'/$1"' _ {}

echo "==> Encoding video"
mkdir -p video/out
ffmpeg -y -framerate "$FPS" -i "$PNG_DIR/frame_%05d.png" \
  -vf "scale=1080:1920:flags=lanczos" \
  -c:v libx264 -profile:v high -pix_fmt yuv420p -crf 18 -preset medium \
  -movflags +faststart \
  "$OUT"

echo "==> Done: $OUT"
