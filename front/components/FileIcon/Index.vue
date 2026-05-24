<script setup lang="ts">
import { cx } from 'class-variance-authority'
import FileIcon from './File.vue'
import ImageIcon from './Image.vue'
import VideoIcon from './Video.vue'
import { fileTypeFromBuffer } from 'file-type'

export type filePreview = {
    type: string
    name: string
    size: number
}

const props = withDefaults(
    defineProps<{
        file: File | filePreview
        class?: string
        size?: 'sm' | 'md' | 'lg'
    }>(),
    {
        size: 'md',
    }
)
const isFile = computed(() => props?.file instanceof File)
const { state: fileType } = useAsyncState(async () => {
    if (!isFile.value) {
        return null
    }
    if (!!props?.file?.type) {
        return props?.file?.type
    }
    const { mime } = (await fileTypeFromBuffer(await (props?.file as File)?.arrayBuffer())) || {}
    return mime
}, null)

const isImage = computed(() => isFile.value && fileType.value?.startsWith('image/'))
const isVideo = computed(() => isFile.value && fileType.value?.startsWith('video/'))
</script>

<template>
    <div v-if="isImage || isVideo" :class="cx('flex overflow-hidden', size === 'sm' && 'max-w-20 max-h-16', size === 'md' && 'max-w-30 max-h-20')">
        <component
            :is="isImage ? ImageIcon : VideoIcon"
            :file="props?.file as File"
            class="block max-w-full max-h-full object-contain border border-black/20 rounded"
        />
    </div>
    <div
        v-else
        :class="
            cx(
                'flex justify-center items-center bg-white/80',
                size === 'sm' && 'size-7 rounded-md',
                size === 'md' && 'size-16 rounded-xl',
                props?.class
            )
        "
    >
        <component :is="FileIcon" :file="props?.file" class="size-[62.5%]" />
    </div>
</template>
