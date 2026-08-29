<script setup lang="ts">
import getFileSize from '@/lib/getFileSize'

const props = defineProps<{
    files: {
        id?: string
        file_name: string
        size?: number
        mime_type?: string
    }[]
}>()
</script>

<template>
    <div class="grid gap-2">
        <div
            v-for="file in props.files"
            :key="file.id || file.file_name"
            class="group flex min-w-0 items-center gap-2.5 rounded-lg bg-black/5 px-2.5 py-2 transition-colors hover:bg-black/10 dark:bg-white/[0.06] dark:hover:bg-white/[0.1]"
        >
            <FileIcon :file="{ name: file.file_name, type: file.mime_type || '', size: file.size || 0 }" size="sm" class="shrink-0" />
            <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium leading-tight">{{ file.file_name }}</p>
                <p class="mt-1 truncate text-xs text-muted-foreground">
                    {{ getFileSize(file.size ?? 0) }}
                </p>
            </div>
            <slot :file="file" />
        </div>
    </div>
</template>
