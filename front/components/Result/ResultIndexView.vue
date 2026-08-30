<script lang="ts" setup>
import TextTranslateResult from '@/components/Result/TextTranslateResult.vue'
import ImageCompressResult from '@/components/Result/ImageCompressResult.vue'
import ImageConvertResult from '@/components/Result/ImageConvertResult.vue'
import type { filehandleData, handleComponent, handleKey, texthandleData } from './types'

const props = defineProps<{
    data: filehandleData | texthandleData
}>()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

type ResultData = filehandleData | texthandleData
type ResultHandle = (data: ResultData) => void | Promise<void>
type ResultItem = { key: handleKey; component: handleComponent; handle?: never } | { key: handleKey; component?: never; handle: ResultHandle }

const router = useRouter()
const { createFileShare, createTextShare } = useMyAppShare()

const handleShare: ResultHandle = async (data) => {
    const { config } = data || {}
    let response: Awaited<ReturnType<typeof createFileShare>>
    if ('files' in data) {
        const { files } = data || {}
        response = await createFileShare({
            files: files.map(({ id, file }) => ({ id, file_name: file.name })),
            config: config as any,
        })
    } else {
        const { text } = data || {}
        response = await createTextShare({ text, config })
    }
    const id = response.data?.id
    if (!id) throw new Error('Failed to create share')

    await router.push(`/s/manage/${id}`)
}

const handleList: ResultItem[] = [
    { handle: handleShare, key: 'file-share' },
    { handle: handleShare, key: 'text-share' },
    { component: TextTranslateResult, key: 'text-translate' },
    { component: ImageCompressResult, key: 'file-image-compress' },
    { component: ImageConvertResult, key: 'file-image-convert' },
]

const activeHandle = computed(() => {
    return handleList.find((item) => item.key === props?.data?.handle_type)
})

onMounted(() => {
    activeHandle.value?.handle?.(props.data)
})
</script>
<template>
    <component
        v-if="activeHandle?.component && 'files' in data"
        :is="activeHandle.component"
        :data="data"
        @change="(key: string) => emit('change', key)"
    />
    <component
        v-if="activeHandle?.component && 'text' in data"
        :is="activeHandle.component"
        :data="data"
        @change="(key: string) => emit('change', key)"
    />
</template>
