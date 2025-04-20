<script setup lang="ts">
import FileUpload from '~/components/FileUpload.vue'
import { cx } from 'class-variance-authority'
import { filesize } from 'filesize'

const props = defineProps<{
    name: string
    rules?: string
}>()
const { value, setValue } = useField<File>(props?.name, props?.rules)
watch(value, (v) => {
    console.log('value', v)
})

const fileInfo = computed(() => {
    const [, name, ext] = value?.value?.name?.match(/^(.+)\.(.+)$/) || []
    return { name, ext }
})

</script>

<template>
    <FileUpload @onChange="(file) => {
        setValue(file)
    }" v-slot="{ isOverDropZone }">
        <div :class="cx('bg-white/50 rounded-md p-2 w-full h-40 flex flex-col items-center justify-center border border-dashed border-black/20 cursor-pointer text-gray-500 gap-3',
            isOverDropZone && '!bg-green-100/50 '
        )">
            <template v-if="!!value">
                <FileIcon :file="value" />
                <div class="flex flex-col gap-0.5 items-center">
                    <div class="flex max-w-30 w-full">
                        <div class="truncate">{{ fileInfo?.name }}</div>
                        <div>{{ `.${fileInfo?.ext}` }}</div>
                    </div>
                    <div class="text-xs opacity-50">{{ filesize(value?.size) }}</div>
                </div>
            </template>
            <template v-else>
                <LucideUpload class="size-10" />
                <div class="text-sm">
                    拖拽文件 或 点击上传
                </div>
            </template>

        </div>
    </FileUpload>
</template>