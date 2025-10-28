#!/usr/bin/env node

/**
 * Simple link validator for Goca blog
 * Checks that all internal blog links include the /goca/ base path
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const BASE_PATH = '/goca';
const BLOG_DIR = path.join(__dirname, 'blog');

/**
 * Extract all links from markdown content (excluding code blocks)
 */
function extractLinks(content) {
    // Remove code blocks first
    const withoutCodeBlocks = content.replace(/```[\s\S]*?```/g, '');

    const links = new Set();

    // [text](url)
    const mdRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
    let match;
    while ((match = mdRegex.exec(withoutCodeBlocks)) !== null) {
        links.add(match[2]);
    }

    // href="url"
    const hrefRegex = /href=["']([^"']+)["']/g;
    while ((match = hrefRegex.exec(withoutCodeBlocks)) !== null) {
        links.add(match[1]);
    }

    return Array.from(links);
}/**
 * Check if link needs base path validation
 */
function needsValidation(link) {
    if (!link.startsWith('/')) return false;
    if (link.startsWith('http')) return false;
    if (link.includes('/blog/')) return true;
    return false;
}

/**
 * Scan directory recursively
 */
function scanDir(dir) {
    const files = [];
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
        const full = path.join(dir, entry.name);
        if (entry.isDirectory()) {
            files.push(...scanDir(full));
        } else if (entry.name.endsWith('.md')) {
            files.push(full);
        }
    }
    return files;
}

/**
 * Main validation
 */
console.log('🔍 Validating blog links...\n');

const files = scanDir(BLOG_DIR);
const issues = [];

for (const file of files) {
    const content = fs.readFileSync(file, 'utf-8');
    const links = extractLinks(content);
    const relPath = path.relative(__dirname, file).replace(/\\/g, '/');

    for (const link of links) {
        if (needsValidation(link) && !link.startsWith(BASE_PATH)) {
            issues.push({ file: relPath, link });
        }
    }
}

if (issues.length === 0) {
    console.log('✅ All blog links are valid!');
    console.log(`   Scanned ${files.length} files\n`);
    process.exit(0);
}

console.log(`❌ Found ${issues.length} invalid link(s):\n`);
for (const { file, link } of issues) {
    console.log(`   File: ${file}`);
    console.log(`   ❌ ${link}`);
    console.log(`   ✅ ${link.replace('/blog/', BASE_PATH + '/blog/')}\n`);
}
process.exit(1);
