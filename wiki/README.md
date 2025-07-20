# Wiki Configuration

This directory contains the complete documentation for the Goca project.

## How to Use This Wiki

### For GitHub Users
1. Go to your repository on GitHub
2. Click on the "Wiki" tab
3. Click "Create the first page"
4. Copy the content from `Home.md` to create your home page
5. Continue adding pages by copying from the respective `.md` files

### For Local Documentation
You can serve this documentation locally using any markdown server:

```bash
# Using Python
cd wiki/
python -m http.server 8000

# Using Node.js (if you have markdown-it)
npx markdown-it-cli *.md

# Using Go (if you have a markdown server)
go run github.com/shurcooL/markdownfmt/cmd/markdownfmt *.md
```

## Wiki Structure

The wiki is organized as follows:

### Getting Started
- `Home.md` - Main wiki page with overview
- `Installation.md` - How to install Goca
- `Getting-Started.md` - Quick start guide
- `Complete-Tutorial.md` - Full e-commerce tutorial

### Command Reference
- `Command-Init.md` - goca init command
- `Command-Feature.md` - goca feature command  
- `Command-Version.md` - goca version command
- (Additional command pages can be added)

### Architecture & Concepts
- `Clean-Architecture.md` - Clean Architecture principles
- `Project-Structure.md` - Directory organization
- (Additional architecture pages can be added)

## Contributing to the Wiki

1. Edit the appropriate `.md` file in this directory
2. Test your changes locally
3. Commit and push to update the repository
4. Update the GitHub wiki by copying the updated content

## Markdown Guidelines

- Use proper heading hierarchy (H1 for page title, H2 for sections, etc.)
- Include code examples with appropriate language tags
- Use emoji for visual appeal (✅ ❌ 🎯 🚀 etc.)
- Include navigation links at the bottom of each page
- Keep examples practical and tested
