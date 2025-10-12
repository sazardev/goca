import { h } from 'vue'
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import './style.css'
import Badge from './components/Badge.vue'
import FeatureCard from './components/FeatureCard.vue'

export default {
    extends: DefaultTheme,
    Layout: () => {
        return h(DefaultTheme.Layout, null, {})
    },
    enhanceApp({ app, router, siteData }) {
        // Register global components
        app.component('Badge', Badge)
        app.component('FeatureCard', FeatureCard)
    }
} satisfies Theme
