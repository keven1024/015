<script setup lang="ts">
import FileUpload from '~/components/FileUpload.vue'
import { cx } from 'class-variance-authority'


const props = defineProps<{
    name: string
    rules?: string
}>()
const { value, setValue } = useField<File>(props?.name, props?.rules)

</script>

<template>
    <FileUpload @onChange="(file) => {
        setValue(file)
    }" v-slot="{ isOverDropZone }">
        <div :class="cx('bg-white/50 rounded-md p-2 w-full h-40 flex flex-col items-center justify-center border border-dashed border-black/20 cursor-pointer text-gray-500 gap-3',
            isOverDropZone && '!bg-green-100/50 '
        )">
            <template v-if="!!value">
                <FilePreviewView :value="value" />
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