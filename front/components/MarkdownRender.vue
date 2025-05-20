<script setup lang="ts">
import markdownit from 'markdown-it'
import { cx } from 'class-variance-authority'
const props = defineProps<{
    markdown: string
    class?: string
}>()
const renderHtml = computed(() => {
    const md = markdownit()
    return md.render(props?.markdown || '')
})
</script>
<template>
    <div :class="cx('prose', props?.class)" v-html="renderHtml" />
</template>