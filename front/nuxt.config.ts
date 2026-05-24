import tailwindcss from '@tailwindcss/vite'
import { defineNuxtConfig } from 'nuxt/config'
import getApiBaseUrl from './lib/getApiBaseUrl'
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: { enabled: true },
    css: ['@/assets/css/main.css'],
    modules: [
        // '@serwist/nuxt',
        '@vueuse/nuxt',
        'motion-v/nuxt',
        'shadcn-nuxt',
        '@vee-validate/nuxt',
        '@pinia/nuxt',
        '@nuxt/image',
        '@nuxtjs/i18n',
        'vue3-pixi-nuxt',
        'nuxt-lucide-icons',
    ],
    // serwist: {},
    i18n: {
        strategy: 'no_prefix',
        defaultLocale: 'en',
        locales: [
            { code: 'zh-CN', name: '中文(简体)', file: 'zh-CN.json' },
            { code: 'en', name: 'English', file: 'en.json' },
            { code: 'ja', name: '日本語', file: 'ja.json' },
            { code: 'ko', name: '한국어', file: 'ko.json' },
            { code: 'fr', name: 'Français', file: 'fr.json' },
            { code: 'de', name: 'Deutsch', file: 'de.json' },
            { code: 'zh-TW', name: '中文(繁體)', file: 'zh-TW.json' },
        ],
    },
    vite: {
        plugins: [tailwindcss()],
        optimizeDeps: {
            include: ['eventemitter3'],
        },
    },
    nitro: {
        routeRules: {
            '/api/**': {
                proxy: `${getApiBaseUrl()}/**`,
            },
        },
    },
    devServer: {
        port: 5000,
    },
})
