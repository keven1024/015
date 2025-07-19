<!-- app.vue -->
<script setup lang="ts">
import { defineRule } from 'vee-validate'
import { all } from '@vee-validate/rules'
Object.entries(all).forEach(([name, rule]) => {
    defineRule(name, rule)
})
// @ts-check
// const { $serwist } = useNuxtApp();
// $serwist?.addEventListener("installed", () => {
//   console.log("Serwist installed!");
// });
// void $serwist?.register({ immediate: true });
</script>
<template>
    <NuxtLayout>
        <NuxtPage />
    </NuxtLayout>
</template>
