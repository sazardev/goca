/**
 * OG Image Generator for Goca Blog Articles
 * 
 * Automatically generates Open Graph images for all blog articles using Sharp (Canvas-based).
 * Images include: Goca logo, article title, metadata badges, and brand colors.
 * 
 * Usage: npm run build:og-images
 */

import fs from 'fs-extra';
import path from 'path';
import { fileURLToPath } from 'url';
import { glob } from 'glob';
import matter from 'gray-matter';
import sharp from 'sharp';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Constants
const DOCS_ROOT = path.resolve(__dirname, '..');
const ARTICLES_DIR = path.join(DOCS_ROOT, 'blog', 'articles');
const OG_IMAGES_DIR = path.join(DOCS_ROOT, 'public', 'og-images');
const LOGO_PATH = path.join(DOCS_ROOT, 'public', 'logo.svg');

// Goca Brand Colors (matching VitePress theme)
const COLORS = {
    primary: '#00ADD8',      // Go/Goca brand cyan
    primaryDark: '#0099C0',  // Darker cyan
    secondary: '#00D8FF',    // Light cyan accent
    background: '#0f172a',   // Dark blue-gray
    text: '#ffffff',
    textMuted: '#94a3b8',
};

/**
 * Generate SVG template for OG image
 */
function generateSVGTemplate(title, description, badges = []) {
    const titleFontSize = title.length > 50 ? 56 : 72;
    const badgesHTML = badges.slice(0, 3).map(badge => `
        <g>
            <rect rx="8" fill="${COLORS.background}" fill-opacity="0.3" stroke="${COLORS.primary}" stroke-width="2"/>
            <text fill="${COLORS.primary}" font-size="18" font-weight="600" text-anchor="middle">${badge}</text>
        </g>
    `).join('');

    return `
<svg width="1200" height="630" xmlns="http://www.w3.org/2000/svg">
    <defs>
        <linearGradient id="bgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:${COLORS.background};stop-opacity:1" />
            <stop offset="100%" style="stop-color:#1e293b;stop-opacity:1" />
        </linearGradient>
        <radialGradient id="glowTop" cx="80%" cy="20%" r="40%">
            <stop offset="0%" style="stop-color:${COLORS.primary};stop-opacity:0.15" />
            <stop offset="100%" style="stop-color:${COLORS.primary};stop-opacity:0" />
        </radialGradient>
        <radialGradient id="glowBottom" cx="20%" cy="80%" r="40%">
            <stop offset="0%" style="stop-color:${COLORS.primary};stop-opacity:0.1" />
            <stop offset="100%" style="stop-color:${COLORS.primary};stop-opacity:0" />
        </radialGradient>
    </defs>
    
    <!-- Background -->
    <rect width="1200" height="630" fill="url(#bgGradient)"/>
    <rect width="1200" height="630" fill="url(#glowTop)"/>
    <rect width="1200" height="630" fill="url(#glowBottom)"/>
    
    <!-- Logo + Brand -->
    <g transform="translate(60, 60)">
        <rect width="70" height="70" rx="16" fill="${COLORS.primary}"/>
        <text x="35" y="53" font-family="Arial, sans-serif" font-size="48" font-weight="bold" fill="white" text-anchor="middle">G</text>
        <text x="90" y="40" font-family="Arial, sans-serif" font-size="32" font-weight="bold" fill="${COLORS.text}">Goca</text>
        <text x="90" y="68" font-family="Arial, sans-serif" font-size="18" fill="${COLORS.textMuted}">Clean Architecture Blog</text>
    </g>
    
    <!-- Title -->
    <text x="60" y="250" font-family="Arial, sans-serif" font-size="${titleFontSize}" font-weight="bold" fill="${COLORS.text}" style="line-height:1.1">
        ${wrapText(title, titleFontSize > 60 ? 16 : 20)}
    </text>
    
    <!-- Description -->
    ${description ? `
    <text x="60" y="400" font-family="Arial, sans-serif" font-size="24" fill="${COLORS.textMuted}" style="max-width:1000px">
        ${truncateText(description, 100)}
    </text>
    ` : ''}
    
    <!-- Badges + URL -->
    <g transform="translate(60, 540)">
        ${badges.slice(0, 3).map((badge, i) => `
            <rect x="${i * 140}" y="0" width="130" height="40" rx="8" fill="rgba(0, 173, 216, 0.2)" stroke="${COLORS.primary}" stroke-width="2"/>
            <text x="${i * 140 + 65}" y="26" font-family="Arial, sans-serif" font-size="16" font-weight="600" fill="${COLORS.primary}" text-anchor="middle">${badge}</text>
        `).join('')}
        
        <text x="1080" y="26" font-family="Arial, sans-serif" font-size="20" font-weight="500" fill="${COLORS.textMuted}" text-anchor="end">sazardev.github.io/goca</text>
    </g>
</svg>
    `.trim();
}

/**
 * Wrap text to fit within width (simple line breaking)
 */
function wrapText(text, maxLength) {
    if (text.length <= maxLength) return text;

    const words = text.split(' ');
    const lines = [];
    let currentLine = '';

    for (const word of words) {
        if ((currentLine + word).length <= maxLength) {
            currentLine += (currentLine ? ' ' : '') + word;
        } else {
            lines.push(currentLine);
            currentLine = word;
        }
    }
    if (currentLine) lines.push(currentLine);

    return lines.slice(0, 3).map((line, i) => `
        <tspan x="60" dy="${i === 0 ? 0 : 1.1}em">${line}</tspan>
    `).join('');
}

/**
 * Truncate text to max length
 */
function truncateText(text, maxLength) {
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
}/**
 * Parse frontmatter from markdown file
 */
async function parseArticleFrontmatter(filePath) {
    const content = await fs.readFile(filePath, 'utf-8');
    const { data } = matter(content);
    return data;
}

/**
 * Generate OG image for a single article
 */
async function generateOGImage(articlePath) {
    try {
        const frontmatter = await parseArticleFrontmatter(articlePath);
        const title = frontmatter.title || 'Untitled Article';
        const description = frontmatter.description || '';

        // Extract badges from frontmatter (if available)
        const badges = [];
        if (frontmatter.tags) {
            badges.push(...frontmatter.tags.slice(0, 3));
        } else {
            // Default badges based on title keywords
            if (title.toLowerCase().includes('entity') || title.toLowerCase().includes('entities')) {
                badges.push('Domain', 'Entities', 'DDD');
            } else if (title.toLowerCase().includes('use case')) {
                badges.push('Use Cases', 'Services', 'Clean Architecture');
            } else {
                badges.push('Clean Architecture', 'Go', 'Goca');
            }
        }

        // Generate filename from article path
        const relativePath = path.relative(ARTICLES_DIR, articlePath);
        const filename = relativePath.replace(/\.md$/, '.png').replace(/\\/g, '-').replace(/\//g, '-');
        const outputPath = path.join(OG_IMAGES_DIR, filename);

        console.log(`📸 Generating OG image: ${filename}`);
        console.log(`   Title: ${title}`);
        console.log(`   Badges: ${badges.join(', ')}`);

        // Generate SVG
        const svg = generateSVGTemplate(title, description, badges);

        // Convert SVG to PNG using Sharp
        await sharp(Buffer.from(svg))
            .png()
            .toFile(outputPath);

        console.log(`✅ Generated: ${outputPath}\n`);

        return {
            articlePath: relativePath.replace(/\\/g, '/'),
            imagePath: `/og-images/${filename}`,
        };
    } catch (error) {
        console.error(`❌ Error generating OG image for ${articlePath}:`, error.message);
        console.error(error.stack);
        return null;
    }
}/**
 * Main function
 */
async function main() {
    console.log('🚀 Goca OG Image Generator\n');
    console.log('📁 Scanning articles directory...\n');

    // Ensure output directory exists
    await fs.ensureDir(OG_IMAGES_DIR);

    // Find all article markdown files
    const articleFiles = await glob('**/*.md', {
        cwd: ARTICLES_DIR,
        absolute: true,
        ignore: ['**/index.md'], // Skip index files
    });

    console.log(`📝 Found ${articleFiles.length} articles\n`);

    // Generate OG images for all articles
    const results = [];
    for (const articleFile of articleFiles) {
        const result = await generateOGImage(articleFile);
        if (result) {
            results.push(result);
        }
    }

    // Generate mapping file for VitePress config
    const mappingPath = path.join(DOCS_ROOT, '.vitepress', 'og-images-map.json');
    await fs.writeJson(mappingPath, results, { spaces: 2 });

    console.log(`\n✨ Done! Generated ${results.length} OG images`);
    console.log(`📄 Mapping saved to: og-images-map.json\n`);
}

// Run
main().catch(console.error);
