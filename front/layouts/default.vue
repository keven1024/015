<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
const { locale } = useI18n()
await useSeo({ locale: locale.value })
const appConfig = useMyAppConfig()
const bgUrl = computed(() => appConfig.value?.site_bg_url)
</script>
<template>
    <div class="h-screen w-screen">
        <GlobalDrawer />
        <GlobalDayjs />
        <Toaster position="top-center" richColors closeButton />
        <img class="w-full h-full object-cover absolute inset-0 -z-[1] bg-gradient-to-bl from-primary/40 to-primary" :src="bgUrl" />
        <div class="h-full w-full flex flex-col items-center lg:p-10 p-5 overflow-y-auto">
            <Navbar />
            <slot />
        </div>
    </div>
</template>
