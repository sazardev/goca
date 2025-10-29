# Creating Blog Articles with Auto-Generated OG Images

## Quick Start

### 1. Create Your Article

Create a new markdown file in `blog/articles/`:

```markdown
---
layout: doc
title: Your Article Title
titleTemplate: Articles | Goca Blog
description: A concise description that will appear in social media previews
tags:
  - Clean Architecture
  - Domain
  - Go
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

<!-- Your article content here -->
```

### 2. Add to Index

Update `blog/articles/index.md`:

```html
<div class="article-item">
  <h3><a href="/goca/blog/articles/your-article">Your Article Title</a></h3>
  <p>Article description...</p>
  <div class="article-meta">
    <Badge icon="lucide:calendar">October 29, 2025</Badge>
    <Badge icon="lucide:tag">Clean Architecture</Badge>
  </div>
</div>
```

### 3. Update Blog Home

Add to recent articles in `blog/index.md`:

```markdown
- **[Your Article Title](/goca/blog/articles/your-article)** - October 29, 2025  
  Brief description of your article
```

### 4. Build & Test

```bash
# Generate OG images
npm run build:og-images

# Preview locally
npm run dev

# Full build
npm run build
```

## OG Image Auto-Generation

**What You Get:**
- ✅ Custom 1200x630px image for each article
- ✅ Goca branding with logo and colors
- ✅ Article title prominently displayed
- ✅ Up to 3 badge/tags shown
- ✅ Automatic social media meta tags

**What You Need:**
- ✅ `title` in frontmatter (required)
- ✅ `description` in frontmatter (recommended)
- ✅ `tags` array (optional, auto-detected if missing)

## Testing Social Media Previews

After deploying, test your article with:

1. **Facebook Sharing Debugger**  
   https://developers.facebook.com/tools/debug/
   
2. **Twitter Card Validator**  
   https://cards-dev.twitter.com/validator
   
3. **LinkedIn Post Inspector**  
   https://www.linkedin.com/post-inspector/

## Frontmatter Best Practices

```yaml
---
layout: doc                          # Always use 'doc'
title: Your Title Here               # Keep under 60 characters
titleTemplate: Articles | Goca Blog  # Standard template
description: A clear, SEO-friendly description that explains what readers will learn. Keep under 160 characters for best results.
tags:                                # Up to 5 tags (first 3 shown in OG image)
  - Clean Architecture
  - Domain Entities
  - Go
  - DDD
  - Goca
---
```

## Common Issues

### OG Image Not Generated?

**Check:**
- Article file is not named `index.md`
- Title exists in frontmatter
- File is in `blog/articles/` directory
- No syntax errors in frontmatter

### Wrong Preview on Social Media?

**Fix:**
1. Clear social media cache using debugger tools
2. Add `?v=2` to URL when sharing (cache buster)
3. Wait 24 hours for cache to expire naturally

### Custom Image Styling?

Edit `scripts/generate-og-images.js`:
- Change colors in `COLORS` object
- Modify layout in `generateOGImageTemplate()`
- Adjust font sizes

## Example Article Structure

```markdown
---
layout: doc
title: Mastering Repository Pattern in Go
titleTemplate: Articles | Goca Blog
description: Learn how to implement the Repository pattern in Go following Clean Architecture principles with practical examples from Goca
tags:
  - Repository Pattern
  - Data Access
  - Clean Architecture
---

<script setup>
import Badge from '../../.vitepress/theme/components/Badge.vue'
</script>

<Badge icon="lucide:calendar">October 29, 2025</Badge>
<Badge icon="lucide:user">sazardev</Badge>
<Badge icon="lucide:clock">15 min read</Badge>

## Introduction

Your engaging introduction...

## Main Content

### Section 1
...

### Section 2
...

## Conclusion

Your conclusions and takeaways...
```

## Resources

- **Full Documentation**: [OG-IMAGE-SYSTEM.md](../../OG-IMAGE-SYSTEM.md)
- **VitePress Docs**: https://vitepress.dev/
- **Social Media Image Specs**: https://www.designyourway.net/blog/social-media-image-sizes/

---

Happy writing! 🚀
