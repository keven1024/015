<script lang="ts" setup>
import FileShareResult from '@/components/Result/FileShareResult.vue'
import TextShareResult from '@/components/Result/TextShareResult.vue'
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

const handleList: { component: handleComponent; key: handleKey }[] = [
    { component: FileShareResult, key: 'file-share' },
    { component: TextShareResult, key: 'text-share' },
    { component: TextTranslateResult, key: 'text-translate' },
    { component: ImageCompressResult, key: 'file-image-compress' },
    { component: ImageConvertResult, key: 'file-image-convert' },
]

const activeHandle = computed(() => {
    return handleList.find((item) => item.key === props?.data?.handle_type)
})
// vue这个ts蠢的没边了，本来想写component: FileShareResult | TextShareResult，结果不行
</script>
<template>
    <component v-if="'files' in data" :is="activeHandle?.component" :data="data" @change="(key: string) => emit('change', key)" />
    <component v-if="'text' in data" :is="activeHandle?.component" :data="data" @change="(key: string) => emit('change', key)" />
</template>
