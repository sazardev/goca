import { h, onMounted, watch, nextTick } from 'vue'
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import './style.css'
import Badge from './components/Badge.vue'
import FeatureCard from './components/FeatureCard.vue'
import mermaid from 'mermaid'

export default {
    extends: DefaultTheme,
    Layout: () => {
        return h(DefaultTheme.Layout, null, {})
    },
    enhanceApp({ app, router, siteData }) {
        // Register global components
        app.component('Badge', Badge)
        app.component('FeatureCard', FeatureCard)

        // Initialize Mermaid
        if (typeof window !== 'undefined') {
            onMounted(() => {
                mermaid.initialize({
                    startOnLoad: true,
                    theme: 'default'
                })
            })

            // Re-render diagrams on route change
            watch(
                () => router.route.path,
                () => nextTick(() => {
                    mermaid.contentLoaded()
                }),
                { immediate: true }
            )
        }
    }
} satisfies Theme

