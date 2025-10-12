# OG Image Setup

## Current Status

✅ **og-image.svg** - Vector version (2.1KB)
⚠️  **og-image.png** - Placeholder (needs proper conversion)

## Why Convert SVG to PNG?

Some social media platforms (Facebook, Twitter, LinkedIn) prefer PNG for better compatibility and consistent rendering across different devices and platforms.

## How to Convert SVG to PNG (1200x630px)

### Option 1: Online Converter (Easiest)
1. Go to https://cloudconvert.com/svg-to-png
2. Upload `docs/public/og-image.svg`
3. Set output size: **1200x630 pixels**
4. Download and replace `docs/public/og-image.png`

### Option 2: Using ImageMagick (Command Line)
```bash
cd docs/public
convert -background none -size 1200x630 og-image.svg og-image.png
```

### Option 3: Using Inkscape (Best Quality)
```bash
inkscape og-image.svg --export-filename=og-image.png --export-width=1200 --export-height=630
```

### Option 4: Using Node.js Sharp (Automated)
```bash
cd docs
npm install sharp
node -e "const sharp = require('sharp'); sharp('public/og-image.svg').resize(1200, 630).png().toFile('public/og-image.png');"
```

## Verification

After conversion, check:
1. File size should be ~50-200KB (PNG)
2. Dimensions: exactly 1200x630 pixels
3. Test on: https://www.opengraph.xyz/url/https%3A%2F%2Fsazardev.github.io%2Fgoca%2F

## SEO Implementation Checklist

✅ Favicon (using existing favicon.ico and logo.svg)
✅ Open Graph tags (Facebook, LinkedIn)
✅ Twitter Cards
✅ Canonical URLs
✅ Schema.org structured data (JSON-LD)
✅ robots.txt
✅ site.webmanifest (PWA support)
✅ Comprehensive meta keywords
✅ Mobile optimization
✅ OG image (SVG created, PNG pending conversion)

## Testing Your SEO

After deployment, test with:
- **Facebook**: https://developers.facebook.com/tools/debug/
- **Twitter**: https://cards-dev.twitter.com/validator
- **LinkedIn**: https://www.linkedin.com/post-inspector/
- **General**: https://www.opengraph.xyz/
- **Google**: Search Console (https://search.google.com/search-console)

## Notes

The current PNG is just a copy of the SVG (2.1KB) which will work but isn't optimal. For best results, convert to a proper PNG using one of the methods above before pushing to production.
