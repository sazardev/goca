import { defineConfig } from 'vitepress'

export default defineConfig({
    title: 'GOCA CLI',
    description: 'Clean Architecture Code Generator for Go',
    base: '/goca/',
    ignoreDeadLinks: true,

    themeConfig: {
        logo: '/logo.svg', nav: [
            { text: 'Home', link: '/' },
            { text: 'Getting Started', link: '/getting-started' },
            { text: 'Commands', link: '/commands/' },
            { text: 'Guide', link: '/guide/' },
            { text: 'v2.0.0', link: '/releases/v2.0.0' }
        ],

        sidebar: {
            '/': [
                {
                    text: 'Introduction',
                    items: [
                        { text: 'What is GOCA?', link: '/' },
                        { text: 'Getting Started', link: '/getting-started' },
                        { text: 'Installation', link: '/installation' }
                    ]
                },
                {
                    text: 'Commands',
                    items: [
                        { text: 'Overview', link: '/commands/' },
                        { text: 'init', link: '/commands/init' },
                        { text: 'feature', link: '/commands/feature' },
                        { text: 'entity', link: '/commands/entity' },
                        { text: 'handler', link: '/commands/handler' },
                        { text: 'repository', link: '/commands/repository' },
                        { text: 'usecase', link: '/commands/usecase' }
                    ]
                },
                {
                    text: 'Configuration',
                    items: [
                        { text: 'YAML Configuration', link: '/configuration/yaml' },
                        { text: 'Advanced Config', link: '/configuration/advanced' },
                        { text: 'Migration Guide', link: '/configuration/migration' }
                    ]
                },
                {
                    text: 'Releases',
                    items: [
                        { text: 'v2.0.0 (Latest)', link: '/releases/v2.0.0' },
                        { text: 'v1.0.1', link: '/releases/v1.0.1' },
                        { text: 'Changelog', link: '/releases/changelog' }
                    ]
                }
            ]
        },

        socialLinks: [
            { icon: 'github', link: 'https://github.com/sazardev/goca' }
        ],

        footer: {
            message: 'Released under the MIT License.',
            copyright: 'Copyright © 2024-2025 SazarDev'
        },

        search: {
            provider: 'local'
        }
    },

    markdown: {
        lineNumbers: true
    }
})
