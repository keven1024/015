import tailwindcss from "@tailwindcss/vite";
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  css: ["~/assets/css/main.css"],
  modules: [
    "@vueuse/nuxt", // '@serwist/nuxt',
    "motion-v/nuxt",
    "nuxt-lucide-icons",
    "shadcn-nuxt",
    "@vee-validate/nuxt",
    "@pinia/nuxt",
    "@nuxt/image",
  ],
  // serwist: {},
  vite: {
    plugins: [tailwindcss()],
  },
  nitro: {
    routeRules: {
      "/api/**": { proxy: "http://127.0.0.1:1323/**" },
    },
  },
  devServer: {
    port: 5000,
  },
});
