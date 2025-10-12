import { defineConfig } from 'vitepress'

export default defineConfig({
    title: "Goca - Go Clean Architecture Generator",
    titleTemplate: ":title | Goca Docs",
    description: "CLI code generator for Go that helps you build production-ready applications following Clean Architecture principles. Generate entities, use cases, repositories, and handlers in seconds.",
    base: '/goca/',
    lang: 'en-US',

    head: [
        // Favicons - using existing files
        ['link', { rel: 'icon', type: 'image/x-icon', href: '/goca/favicon.ico' }],
        ['link', { rel: 'icon', type: 'image/svg+xml', href: '/goca/logo.svg' }],
        ['link', { rel: 'apple-touch-icon', href: '/goca/logo.svg' }],
        ['link', { rel: 'manifest', href: '/goca/site.webmanifest' }],

        // Theme and Mobile
        ['meta', { name: 'theme-color', content: '#00ADD8' }],
        ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0, viewport-fit=cover' }],
        ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
        ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' }],

        // SEO Meta Tags
        ['meta', { name: 'author', content: 'sazardev' }],
        ['meta', { name: 'keywords', content: 'go, golang, clean architecture, code generator, cli tool, rest api, ddd, domain driven design, repository pattern, dependency injection, goca, hexagonal architecture, uncle bob, backend development, go framework' }],
        ['meta', { name: 'robots', content: 'index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1' }],
        ['meta', { name: 'googlebot', content: 'index, follow' }],
        ['meta', { name: 'language', content: 'English' }],
        ['meta', { name: 'revisit-after', content: '7 days' }],

        // Open Graph / Facebook
        ['meta', { property: 'og:type', content: 'website' }],
        ['meta', { property: 'og:url', content: 'https://sazardev.github.io/goca/' }],
        ['meta', { property: 'og:site_name', content: 'Goca Documentation' }],
        ['meta', { property: 'og:title', content: 'Goca - Go Clean Architecture Code Generator' }],
        ['meta', { property: 'og:description', content: 'Build production-ready Go applications following Clean Architecture principles. CLI tool that generates complete features with entities, use cases, repositories, and handlers in seconds.' }],
        ['meta', { property: 'og:image', content: 'https://sazardev.github.io/goca/og-image.png' }],
        ['meta', { property: 'og:image:width', content: '1200' }],
        ['meta', { property: 'og:image:height', content: '630' }],
        ['meta', { property: 'og:image:alt', content: 'Goca - Go Clean Architecture Code Generator' }],
        ['meta', { property: 'og:locale', content: 'en_US' }],

        // Twitter Card
        ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
        ['meta', { name: 'twitter:site', content: '@sazardev' }],
        ['meta', { name: 'twitter:creator', content: '@sazardev' }],
        ['meta', { name: 'twitter:url', content: 'https://sazardev.github.io/goca/' }],
        ['meta', { name: 'twitter:title', content: 'Goca - Go Clean Architecture Code Generator' }],
        ['meta', { name: 'twitter:description', content: 'Build production-ready Go applications following Clean Architecture. Generate complete features with entities, use cases, repositories, and handlers in seconds.' }],
        ['meta', { name: 'twitter:image', content: 'https://sazardev.github.io/goca/og-image.png' }],
        ['meta', { name: 'twitter:image:alt', content: 'Goca - Go Clean Architecture Code Generator' }],

        // Additional SEO
        ['link', { rel: 'canonical', href: 'https://sazardev.github.io/goca/' }],
        ['meta', { name: 'application-name', content: 'Goca' }],
        ['meta', { name: 'apple-mobile-web-app-title', content: 'Goca' }],
        ['meta', { name: 'format-detection', content: 'telephone=no' }],

        // Schema.org structured data
        ['script', { type: 'application/ld+json' }, JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'SoftwareApplication',
            'name': 'Goca',
            'applicationCategory': 'DeveloperApplication',
            'operatingSystem': 'Linux, macOS, Windows',
            'offers': {
                '@type': 'Offer',
                'price': '0',
                'priceCurrency': 'USD'
            },
            'description': 'CLI code generator for Go that helps you build production-ready applications following Clean Architecture principles',
            'url': 'https://sazardev.github.io/goca/',
            'softwareVersion': '2.0.0',
            'author': {
                '@type': 'Person',
                'name': 'sazardev',
                'url': 'https://github.com/sazardev'
            },
            'programmingLanguage': 'Go',
            'downloadUrl': 'https://github.com/sazardev/goca/releases',
            'codeRepository': 'https://github.com/sazardev/goca',
            'license': 'https://opensource.org/licenses/MIT'
        })],
    ],
    themeConfig: {
        logo: '/logo.svg',

        nav: [
            { text: 'Home', link: '/' },
            { text: 'Getting Started', link: '/getting-started' },
            { text: 'Guide', link: '/guide/introduction' },
            { text: 'Commands', link: '/commands/init' },
            { text: 'Features', link: '/features/safety-and-dependencies' },
            { text: 'GitHub', link: 'https://github.com/sazardev/goca' }
        ],

        sidebar: {
            '/guide/': [
                {
                    text: 'Introduction',
                    items: [
                        { text: 'What is Goca?', link: '/guide/introduction' },
                        { text: 'Installation', link: '/guide/installation' },
                        { text: 'Quick Start', link: '/getting-started' },
                    ]
                },
                {
                    text: 'Core Concepts',
                    items: [
                        { text: 'Clean Architecture', link: '/guide/clean-architecture' },
                        { text: 'Configuration', link: '/guide/configuration' },
                        { text: 'Project Structure', link: '/guide/project-structure' },
                        { text: 'Best Practices', link: '/guide/best-practices' },
                    ]
                },
                {
                    text: 'Tutorials',
                    items: [
                        { text: 'Complete Tutorial', link: '/tutorials/complete-tutorial' },
                        { text: 'Building a REST API', link: '/tutorials/rest-api' },
                        { text: 'Adding Features', link: '/tutorials/adding-features' },
                    ]
                }
            ],
            '/commands/': [
                {
                    text: 'Commands Reference',
                    items: [
                        { text: 'Overview', link: '/commands/' },
                        { text: 'goca init', link: '/commands/init' },
                        { text: 'goca feature', link: '/commands/feature' },
                        { text: 'goca entity', link: '/commands/entity' },
                        { text: 'goca usecase', link: '/commands/usecase' },
                        { text: 'goca repository', link: '/commands/repository' },
                        { text: 'goca handler', link: '/commands/handler' },
                        { text: 'goca di', link: '/commands/di' },
                        { text: 'goca integrate', link: '/commands/integrate' },
                        { text: 'goca interfaces', link: '/commands/interfaces' },
                        { text: 'goca messages', link: '/commands/messages' },
                        { text: 'goca version', link: '/commands/version' },
                    ]
                }
            ],
            '/features/': [
                {
                    text: 'Features',
                    items: [
                        { text: 'Safety & Dependencies', link: '/features/safety-and-dependencies' },
                    ]
                }
            ]
        },

        socialLinks: [
            { icon: 'github', link: 'https://github.com/sazardev/goca' }
        ],

        footer: {
            message: 'Released under the MIT License.',
            copyright: 'Copyright © 2025 sazardev'
        },

        search: {
            provider: 'local'
        },

        editLink: {
            pattern: 'https://github.com/sazardev/goca/edit/master/docs/:path',
            text: 'Edit this page on GitHub'
        }
    },
    vite: {
        server: {
            port: 3567,
            strictPort: true,
            open: true,

        }
    }
})
