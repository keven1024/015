import tailwindcss from '@tailwindcss/vite'
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
        'nuxt-lucide-icons',
        'shadcn-nuxt',
        '@vee-validate/nuxt',
        '@pinia/nuxt',
        '@nuxt/image',
        '@nuxtjs/i18n',
        'vue3-pixi-nuxt',
    ],
    // serwist: {},
    i18n: {
        strategy: 'no_prefix',
        defaultLocale: 'en',
        locales: [
            { code: 'zh-CN', name: '中文(简体)', file: 'zh-CN.json' },
            { code: 'en', name: 'English', file: 'en.json' },
        ],
    },
    vite: {
        transformMixedEsModules: true,
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
