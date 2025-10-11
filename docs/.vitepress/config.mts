import { defineConfig } from 'vitepress'

export default defineConfig({
    title: "Goca",
    description: "Go Clean Architecture Code Generator - Build production-ready Go applications following Clean Architecture principles",
    base: '/goca/',
    head: [
        ['link', { rel: 'icon', href: '/goca/favicon.ico' }],
        ['meta', { name: 'theme-color', content: '#00ADD8' }],
        ['meta', { name: 'og:type', content: 'website' }],
        ['meta', { name: 'og:locale', content: 'en' }],
        ['meta', { name: 'og:site_name', content: 'Goca Documentation' }],
    ],

    themeConfig: {
        logo: '/logo.svg',

        nav: [
            { text: 'Home', link: '/' },
            { text: 'Getting Started', link: '/getting-started' },
            { text: 'Guide', link: '/guide/introduction' },
            { text: 'Commands', link: '/commands/init' },
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
    }
})
