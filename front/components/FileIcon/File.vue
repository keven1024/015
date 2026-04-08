<script setup lang="ts">
import { LucideFileAudio, LucideFileVideo, LucideFile, LucideFileCode, LucideFileArchive, LucideFileText } from 'lucide-vue-next'
import type { filePreview } from './Index.vue'

const props = defineProps<{
    file: File | filePreview
}>()
const fileIcon = computed(() => {
    const [baseType, type] = props?.file?.type?.split('/')
    // if (baseType === 'video') {
    //     return LucideFileVideo
    // }
    if (baseType === 'audio') {
        return LucideFileAudio
    }
    if (baseType === 'text' || ['json', 'ld+json', 'html']?.includes(type ?? '')) {
        return LucideFileCode
    }
    if (
        [
            'pdf',
            'msword',
            'vnd.openxmlformats-officedocument.wordprocessingml.document',
            'vnd.ms-excel',
            'vnd.openxmlformats-officedocument.spreadsheetml.sheet',
            'vnd.ms-powerpoint',
            'vnd.openxmlformats-officedocument.presentationml.presentation',
        ].includes(type ?? '')
    ) {
        return LucideFileText
    }
    if (['zip', 'vnd.rar', 'x-tar', 'gz', 'bz2', 'x-7z-compressed'].includes(type ?? '')) {
        return LucideFileArchive
    }
    return LucideFile
})
</script>

<template>
    <component :is="fileIcon" />
</template>
