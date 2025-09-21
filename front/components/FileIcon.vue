<script setup lang="ts">
import { cx } from 'class-variance-authority'
export type filePreview = {
    type: string
    name: string
    size: number
}

import { LucideFileAudio, LucideFileVideo, LucideFile, LucideFileCode, LucideFileArchive, LucideFileText } from 'lucide-vue-next'
const props = defineProps<{
    file: File | filePreview
    class?: string
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

const fileIcon = computed(() => {
    const [baseType, type] = props?.file?.type?.split('/')
    if (baseType === 'video') {
        return LucideFileVideo
    }
    if (baseType === 'audio') {
        return LucideFileAudio
    }
    if (baseType === 'text' || ['json', 'ld+json', 'html']?.includes(type)) {
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
        ].includes(type)
    ) {
        return LucideFileText
    }
    if (['zip', 'vnd.rar', 'x-tar', 'gz', 'bz2', 'x-7z-compressed'].includes(type)) {
        return LucideFileArchive
    }
    return LucideFile
})
</script>

<template>
    <div v-if="!!imageUrl" class="flex max-w-30 max-h-20">
        <div class="object-contain m-auto h-full">
            <img :src="imageUrl" class="w-full h-full border border-black/20 rounded" />
        </div>
    </div>
    <div v-if="!imageUrl" :class="cx('flex justify-center items-center rounded-xl bg-white/80 size-16', props?.class)">
        <component :is="fileIcon" class="size-[62.5%]" />
    </div>
</template>
