<script setup lang="ts">
import type { filePreview } from './Index.vue'

const props = defineProps<{
    file: File | filePreview
}>()

const imageUrl = computed(() => {
    if (props?.file?.type?.startsWith('image/') && props?.file instanceof File) {
        return URL.createObjectURL(props?.file)
    }
    return null
})

onUnmounted(() => {
    if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value)
    }
})
</script>

<template>
    <img v-if="!!imageUrl" :src="imageUrl" />
</template>
