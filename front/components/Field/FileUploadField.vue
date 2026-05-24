<script setup lang="ts">
import FileUpload from '@/components/FileUpload.vue'
import { cx } from 'class-variance-authority'
import type { RuleExpression } from 'vee-validate'
import Button from '../ui/button/Button.vue'
import { PlusIcon } from '@lucide/vue'
import { isEmpty } from 'lodash-es'

const props = defineProps<{
    name: string
    rules?: RuleExpression<File[]>
}>()
const { value, setValue } = useField<File[]>(props?.name, props?.rules)
const { t } = useI18n()

const filterOutSameFile = (files: File[] | undefined, targetFile: File[] | undefined) => {
    return files?.filter((file) => !targetFile?.some((r) => r?.name === file?.name && r?.type === file?.type && r?.size === file?.size)) || []
}

useEventListener(document, 'paste', (evt: ClipboardEvent) => {
    const { files } = evt.clipboardData || {}
    if (files?.length) {
        setValue([...filterOutSameFile(value?.value, Array.from(files)), ...Array.from(files)])
    }
})
</script>

<template>
    <FileUpload
        @onChange="
            (files) => {
                // 这里没hash，我们姑且认为name和size,type都一样的为同一个文件
                setValue([...filterOutSameFile(value, files), ...files])
            }
        "
        v-slot="{ isOverDropZone }"
    >
        <div
            :class="
                cx(
                    'bg-white/50 rounded-md p-2 w-full min-h-40 flex flex-col items-center justify-center border border-dashed border-black/20 cursor-pointer text-gray-500 gap-3 transition-all duration-300',
                    isOverDropZone && '!bg-green-100/50 '
                )
            "
        >
            <template v-if="!isEmpty(value)">
                <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-5 w-full">
                    <div v-for="item in value" :key="item?.name" class="flex flex-col items-center justify-center relative">
                        <div class="absolute top-0 right-0">
                            <Button
                                class="size-5 p-0 bg-red-500/20 hover:bg-red-500/60 text-red-500 hover:text-white"
                                @click="
                                    (e: any) => {
                                        e.stopPropagation()
                                        setValue(filterOutSameFile(value, [item]))
                                    }
                                "
                            >
                                <LucideX class="size-3" />
                            </Button>
                        </div>
                        <FilePreviewView :value="item" />
                    </div>
                    <div class="flex flex-col items-center justify-center opacity-50 gap-1">
                        <div class="size-16 flex justify-center items-center rounded-xl bg-white/80">
                            <PlusIcon class="size-7" />
                        </div>
                        <div class="mb-3">{{ t('page.upload.file.addMore') }}</div>
                    </div>
                </div>
            </template>
            <template v-else>
                <LucideUpload class="size-10" />
                <div class="text-sm select-none">
                    {{ t('page.upload.file.uploadFilePlaceholder') }}
                </div>
            </template>
        </div>
    </FileUpload>
</template>
