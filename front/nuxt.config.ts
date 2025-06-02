import tailwindcss from "@tailwindcss/vite";
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  css: ["~/assets/css/main.css"],
  modules: [
    // '@serwist/nuxt',
    "@vueuse/nuxt",
    "motion-v/nuxt",
    "nuxt-lucide-icons",
    "shadcn-nuxt",
    "@vee-validate/nuxt",
    "@pinia/nuxt",
    "@nuxt/image",
    "@nuxtjs/i18n",
  ],
  // serwist: {},
  i18n: {
    strategy: "no_prefix",
    defaultLocale: "en",
    locales: [
      { code: "zh-CN", name: "中文(简体)", file: "zh-CN.json" },
      { code: "en", name: "English", file: "en.json" },
    ],
  },
  vite: {
    plugins: [tailwindcss()],
  },
  nitro: {
    routeRules: {
      "/api/**": {
        proxy: process.env.API_BASE_URL || "http://127.0.0.1:1323/**",
      },
    },
  },
  devServer: {
    port: parseInt(process.env.PORT || "5000"),
    host: process.env.HOST || "0.0.0.0",
  },
});
