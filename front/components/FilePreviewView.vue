<script lang="ts" setup>
import getFileSize from '~/lib/getFileSize'
import type { filePreview } from './FileIcon.vue'
const props = defineProps<{
    value: File | filePreview
}>()
const fileInfo = computed(() => {
    const [, name, ext] = props?.value?.name?.match(/^(.+)\.(.+)$/) || []
    return { name, ext }
})
</script>

<template>
    <FileIcon :file="value" />
    <div class="flex flex-col gap-0.5 items-center">
        <div class="flex max-w-30 w-full">
            <div class="truncate">{{ fileInfo?.name }}</div>
            <div>{{ `.${fileInfo?.ext}` }}</div>
        </div>
        <div class="text-xs opacity-50">{{ getFileSize(value?.size ?? 0) }}</div>
    </div>
</template>
