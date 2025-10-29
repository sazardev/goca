# OG Image System for Goca Blog

## Overview

The Goca blog uses an automated Open Graph (OG) image generation system that creates custom social media preview images for every article. These images are automatically generated during the build process.

## Features

- **Dynamic Generation**: Each article gets a unique OG image based on its title and metadata
- **Brand Consistency**: Uses Goca's brand colors and design language
- **SEO Optimized**: 1200x630px images optimized for all social platforms
- **Zero Config**: Works automatically for all blog articles

## How It Works

### 1. Build Process

When you run `npm run build`, the system:

1. Scans all article markdown files in `blog/articles/`
2. Extracts title, description, and tags from frontmatter
3. Generates a PNG image for each article using Satori
4. Saves images to `public/og-images/`
5. Creates a mapping file for VitePress

### 2. Image Generation

The script (`scripts/generate-og-images.js`) uses:

- **Satori**: Converts React-like JSX to SVG
- **Sharp**: Converts SVG to optimized PNG
- **Gray Matter**: Parses markdown frontmatter

### 3. VitePress Integration

The `config.mts` file uses `transformHead` to:

- Detect blog article pages
- Inject custom OG image meta tags
- Override default social media tags

## Design Specification

### Colors

```typescript
Primary: #00ADD8      // Go/Goca cyan
Background: #0f172a   // Dark blue-gray
Text: #ffffff         // White
Muted: #94a3b8       // Slate gray
```

### Layout

```
┌─────────────────────────────────────┐
│ [G] Goca                           │  Header: Logo + Brand
│     Clean Architecture Blog        │
│                                     │
│  Article Title Here                │  Main: Large title
│  Brief description of the article  │        + Description
│                                     │
│  [Badge] [Badge] [Badge]  URL      │  Footer: Tags + Domain
└─────────────────────────────────────┘
```

### Dimensions

- **Size**: 1200x630px (recommended by all platforms)
- **Format**: PNG with transparent background support
- **Title**: 56-72px (responsive based on length)
- **Description**: 24px, truncated at 120 chars

## Adding Tags to Articles

To customize the badges shown in OG images, add tags to your article frontmatter:

```yaml
---
title: Your Article Title
description: Article description
tags:
  - Clean Architecture
  - Use Cases
  - Go
---
```

If no tags are provided, the system will auto-detect badges based on title keywords.

## Testing OG Images

### Social Media Debuggers

Test your generated OG images with these tools:

1. **Facebook**: https://developers.facebook.com/tools/debug/
2. **Twitter**: https://cards-dev.twitter.com/validator
3. **LinkedIn**: https://www.linkedin.com/post-inspector/

### Local Preview

After building, check generated images at:
```
docs/public/og-images/
```

View the mapping file at:
```
docs/.vitepress/og-images-map.json
```

## Manual Generation

To regenerate OG images without building the entire site:

```bash
npm run build:og-images
```

## Troubleshooting

### Images Not Updating

1. Delete `public/og-images/` folder
2. Run `npm run build:og-images`
3. Clear browser cache

### Missing Images

- Check that article has frontmatter with `title`
- Verify file is not `index.md` (these are skipped)
- Check console for generation errors

### Custom Styling

Edit `scripts/generate-og-images.js`:

- Modify `COLORS` object for different colors
- Update `generateOGImageTemplate()` for layout changes
- Adjust font sizes in the template

## Examples

### Generated Image for "Understanding Domain Entities"

**URL**: `/og-images/understanding-domain-entities.png`

**Meta Tags Injected**:
```html
<meta property="og:image" content="https://sazardev.github.io/goca/og-images/understanding-domain-entities.png" />
<meta property="og:image:width" content="1200" />
<meta property="og:image:height" content="630" />
<meta property="og:title" content="Understanding Domain Entities in Clean Architecture | Goca Blog" />
<meta name="twitter:card" content="summary_large_image" />
```

## Architecture

```
docs/
├── scripts/
│   └── generate-og-images.js       # Generator script
├── .vitepress/
│   ├── config.mts                  # VitePress config with transformHead
│   └── og-images-map.json         # Generated mapping (gitignored)
├── public/
│   └── og-images/                 # Generated PNG files (gitignored)
└── blog/
    └── articles/
        ├── understanding-domain-entities.md
        └── mastering-use-cases.md
```

## CI/CD Integration

The OG image generation runs automatically during GitHub Pages deployment:

```yaml
# .github/workflows/deploy-docs.yml
- run: npm run build  # Runs build:og-images first
```

## Best Practices

1. **Keep titles concise**: Long titles (>60 chars) will use smaller font
2. **Add descriptive tags**: Max 3 badges shown
3. **Write clear descriptions**: Truncated at 120 characters
4. **Test before merging**: Use social media debuggers
5. **Check mobile previews**: OG images should be readable on small screens

## Credits

- **Satori**: https://github.com/vercel/satori
- **Sharp**: https://sharp.pixelplumbing.com/
- **VitePress**: https://vitepress.dev/

---

**Note**: OG images are generated at build time and gitignored. GitHub Pages will generate them during deployment.
