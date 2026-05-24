<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
const { locale } = useI18n()
await useSeo({ locale: locale.value })
const appConfig = useMyAppConfig()
const bgUrl = computed(() => appConfig.value?.site_bg_url)
const enableBg = computed(() => appConfig.value?.site_enable_bg ?? true)
</script>
<template>
    <div class="h-screen w-screen">
        <GlobalDrawer />
        <GlobalDayjs />
        <Toaster position="top-center" richColors closeButton />
        <p class="absolute inset-0 z-[-1] bg-linear-to-bl from-primary/40 to-primary">
            <img v-if="enableBg" class="w-full h-full block object-cover" :src="bgUrl" />
        </p>
        <div class="h-full w-full flex flex-col items-center lg:p-10 p-5 overflow-y-auto">
            <Navbar />
            <slot />
        </div>
    </div>
</template>
