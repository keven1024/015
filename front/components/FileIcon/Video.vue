<script setup lang="ts">
import getVideoFileThumbnail from '@/lib/getVideoFileThumbnail'

const props = defineProps<{
    file: File
}>()

const { state: thumbnailUrl } = useAsyncState(async () => {
    if (props.file.type.startsWith('video/')) {
        return await getVideoFileThumbnail(props.file)
    }
    return null
}, null)

onUnmounted(() => {
    if (!!thumbnailUrl.value) {
        URL.revokeObjectURL(thumbnailUrl.value)
    }
})
</script>

<template>
    <div v-if="thumbnailUrl" class="relative grayscale-50 overflow-hidden">
        <img :src="thumbnailUrl" class="object-contain block max-w-full max-h-full" />
        <LucidePlay class="size-[40%] absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white" />
    </div>
</template>
