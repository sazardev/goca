# Goca Documentation

This directory contains the VitePress documentation for Goca.

##  Quick Start

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Install Dependencies

```bash
cd docs
npm install
```

### Run Development Server

```bash
npm run dev
```

The documentation will be available at `http://localhost:5173`

### Build for Production

```bash
npm run build
```

The static files will be generated in `.vitepress/dist/`

### Preview Production Build

```bash
npm run preview
```

## 📁 Structure

```
docs/
├── .vitepress/
│   ├── config.mts          # VitePress configuration
│   └── theme/
│       ├── index.ts        # Theme customization
│       └── style.css       # Custom styles
├── guide/
│   ├── introduction.md     # What is Goca?
│   ├── installation.md     # Installation guide
│   ├── clean-architecture.md
│   ├── project-structure.md
│   └── best-practices.md
├── commands/
│   ├── init.md
│   ├── feature.md
│   ├── entity.md
│   └── ...
├── tutorials/
│   ├── complete-tutorial.md
│   └── ...
├── index.md                # Homepage
├── getting-started.md      # Quick start guide
└── package.json
```

## 🎨 Customization

### Theme Colors

Edit `.vitepress/theme/style.css` to change colors:

```css
:root {
  --vp-c-brand-1: #00ADD8;  /* Go cyan */
  --vp-c-brand-2: #00758F;
  --vp-c-brand-3: #00A29C;
}
```

### Navigation

Edit `.vitepress/config.mts` to modify navigation and sidebar:

```ts
export default defineConfig({
  themeConfig: {
    nav: [...],
    sidebar: {...}
  }
})
```

## 📝 Writing Documentation

### Markdown Extensions

VitePress supports many markdown extensions:

#### Code Groups

```md
::: code-group
```bash [npm]
npm install goca
```

```bash [yarn]
yarn add goca
```
:::
```

#### Custom Containers

```md
::: tip
This is a tip
:::

::: warning
This is a warning
:::

::: danger
This is a danger message
:::

::: details Click to see more
Hidden content
:::
```

#### Code Highlighting

```md
```go{1,3-5}
package main // highlighted

func main() {
    // these lines
    // are highlighted
}
```
```

##  Deployment

The documentation is automatically deployed to GitHub Pages when you push to the `master` branch.

The deployment is handled by GitHub Actions (`.github/workflows/deploy-docs.yml`).

### Manual Deployment

If you need to deploy manually:

```bash
# Build
npm run build

# The dist folder is ready to be deployed
ls .vitepress/dist
```

##  Resources

- [VitePress Documentation](https://vitepress.dev/)
- [Markdown Extensions](https://vitepress.dev/guide/markdown)
- [Theme Configuration](https://vitepress.dev/reference/default-theme-config)

## 🤝 Contributing

When adding new documentation:

1. Create the markdown file in the appropriate directory
2. Add it to the sidebar in `.vitepress/config.mts`
3. Test locally with `npm run dev`
4. Submit a pull request

##  Support

- GitHub Issues: https://github.com/sazardev/goca/issues
- Discussions: https://github.com/sazardev/goca/discussions
