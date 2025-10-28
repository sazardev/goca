/**
 * Route Validation Script for Goca Blog
 * 
 * This script validates all internal links in the blog section to prevent 404 errors.
 * Run this before committing changes to ensure all routes are valid.
 * 
 * Usage: node validate-blog-routes.js
 */

const fs = require('fs');
const path = require('path');

// Configuration
const BLOG_DIR = path.join(__dirname, 'blog');
const BASE_PATH = '/goca';

// Colors for terminal output
const colors = {
    reset: '\x1b[0m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
};

// Track all markdown files and their links
const allFiles = new Set();
const allLinks = new Map(); // file -> links[]
const brokenLinks = [];

/**
 * Recursively find all markdown files
 */
function findMarkdownFiles(dir, fileList = []) {
    const files = fs.readdirSync(dir);

    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);

        if (stat.isDirectory()) {
            findMarkdownFiles(filePath, fileList);
        } else if (file.endsWith('.md')) {
            const relativePath = path.relative(BLOG_DIR, filePath)
                .replace(/\\/g, '/')
                .replace(/\.md$/, '');
            allFiles.add(relativePath);
            fileList.push(filePath);
        }
    });

    return fileList;
}

/**
 * Extract all markdown links from content
 */
function extractLinks(content) {
    const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
    const links = [];
    let match;

    while ((match = linkRegex.exec(content)) !== null) {
        links.push({
            text: match[1],
            url: match[2]
        });
    }

    return links;
}

/**
 * Check if a link is valid
 */
function isValidLink(link, sourceFile) {
    // Skip external links
    if (link.startsWith('http://') || link.startsWith('https://')) {
        return true;
    }

    // Skip anchor links
    if (link.startsWith('#')) {
        return true;
    }

    // Normalize the link
    let normalizedLink = link;

    // Remove base path if present
    if (normalizedLink.startsWith(BASE_PATH)) {
        normalizedLink = normalizedLink.substring(BASE_PATH.length);
    }

    // Remove leading slash
    if (normalizedLink.startsWith('/')) {
        normalizedLink = normalizedLink.substring(1);
    }

    // Remove blog prefix
    if (normalizedLink.startsWith('blog/')) {
        normalizedLink = normalizedLink.substring(5);
    }

    // Remove trailing slash
    if (normalizedLink.endsWith('/')) {
        normalizedLink = normalizedLink.slice(0, -1);
    }

    // Handle index files
    if (normalizedLink === '' || normalizedLink === 'index') {
        normalizedLink = 'index';
    }

    // Check if file exists
    return allFiles.has(normalizedLink) || allFiles.has(normalizedLink + '/index');
}

/**
 * Validate all links in a file
 */
function validateFile(filePath) {
    const content = fs.readFileSync(filePath, 'utf-8');
    const links = extractLinks(content);
    const relativePath = path.relative(BLOG_DIR, filePath);

    allLinks.set(relativePath, links);

    links.forEach(link => {
        if (!isValidLink(link.url, relativePath)) {
            brokenLinks.push({
                file: relativePath,
                link: link.url,
                text: link.text
            });
        }
    });
}

/**
 * Main validation function
 */
function validateBlogRoutes() {
    console.log(`${colors.blue}Goca Blog Route Validator${colors.reset}`);
    console.log('='.repeat(50));

    // Find all markdown files
    console.log(`\n${colors.yellow}Scanning blog directory...${colors.reset}`);
    const markdownFiles = findMarkdownFiles(BLOG_DIR);
    console.log(`Found ${markdownFiles.length} markdown files`);

    // Validate each file
    console.log(`\n${colors.yellow}Validating links...${colors.reset}`);
    markdownFiles.forEach(validateFile);

    // Report results
    console.log('\n' + '='.repeat(50));

    if (brokenLinks.length === 0) {
        console.log(`${colors.green}✓ All links are valid!${colors.reset}`);
        console.log(`\nValidated ${allLinks.size} files`);

        let totalLinks = 0;
        allLinks.forEach(links => totalLinks += links.length);
        console.log(`Checked ${totalLinks} links`);

        return true;
    } else {
        console.log(`${colors.red}✗ Found ${brokenLinks.length} broken link(s):${colors.reset}\n`);

        brokenLinks.forEach((broken, index) => {
            console.log(`${index + 1}. ${colors.yellow}${broken.file}${colors.reset}`);
            console.log(`   Text: "${broken.text}"`);
            console.log(`   Link: ${colors.red}${broken.link}${colors.reset}`);
            console.log('');
        });

        console.log(`${colors.red}Please fix these broken links before deploying.${colors.reset}`);
        return false;
    }
}

// Run validation
const isValid = validateBlogRoutes();
process.exit(isValid ? 0 : 1);
