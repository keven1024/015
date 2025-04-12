// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  app: {
    head: {
      htmlAttrs: {
        lang: "en",
        dir: "ltr",
      },
      title: "Home",
      titleTemplate: "%s - PWA App",
      link: [
        {
          rel: "manifest",
          href: "/manifest.json",
        },
        {
          rel: "shortcut icon",
          href: "/favicon.ico",
        },
        {
          rel: "apple-touch-icon",
          href: "/apple-touch-icon.png",
          sizes: "180x180",
        },
      ],
      meta: [
        {
          name: "application-name",
          content: "PWA App",
        },
        {
          name: "description",
          content: "Best PWA app in the world!",
        },
        {
          name: "apple-mobile-web-app-capable",
          content: "yes",
        },
        {
          name: "apple-mobile-web-app-status-bar-style",
          content: "default",
        },
        {
          name: "apple-mobile-web-app-title",
          content: "My Awesome PWA App",
        },
        {
          name: "format-detection",
          content: "telephone=no",
        },
        {
          property: "og:type",
          content: "website",
        },
        {
          property: "og:title",
          content: "My Awesome PWA App",
        },
        {
          property: "og:description",
          content: "Best PWA app in the world!",
        },
        {
          property: "og:site:name",
          content: "PWA App",
        },
        {
          name: "twitter:card",
          content: "website",
        },
        {
          name: "twitter:title",
          content: "My Awesome PWA App",
        },
        {
          name: "twitter:description",
          content: "Best PWA app in the world!",
        },
        {
          name: "theme-color",
          content: "#FFFFFF",
        },
      ],
    },
  },
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@vueuse/nuxt',
    '@serwist/nuxt',
    'motion-v/nuxt',
    "nuxt-lucide-icons",
  ],
  serwist: {},
})