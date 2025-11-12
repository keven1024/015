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
    <div
        :class="
            cx(
                'prose prose-sm [&>*]:outline-none prose-p:my-1 prose-headings:my-2 prose-pre:mb-0 prose-blockquote:border-black/50 selection:bg-primary/20 break-all',
                props?.class
            )
        "
        v-html="renderHtml"
    />
</template>
