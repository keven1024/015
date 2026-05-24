<script setup lang="ts">
import { isHeic, heicTo } from 'heic-to'

const props = defineProps<{
    file: File
}>()

const { state: imageUrl } = useAsyncState(async () => {
    let blob: Blob = props?.file
    if (await isHeic(props?.file)) {
        blob = await heicTo({
            blob: props?.file,
            type: 'image/jpeg',
            quality: 1,
        })
    }
    return URL.createObjectURL(blob)
}, null)

onUnmounted(() => {
    if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value)
    }
})
</script>

<template>
    <img v-if="!!imageUrl" :src="imageUrl" />
</template>
