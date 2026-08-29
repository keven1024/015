<script setup lang="ts">
defineProps<{
    data: {
        text?: string
    }
}>()

const { t } = useI18n()
</script>

<template>
    <div v-if="data.text" class="flex flex-col gap-3">
        <h2 class="text-sm font-medium">{{ t('page.shareManage.textContent') }}</h2>
        <div class="max-h-[55vh] overflow-y-auto rounded-xl border border-black/5 bg-background/60 p-4 dark:border-white/10 sm:p-5">
            <MarkdownRender class="prose w-full max-w-full dark:prose-invert" :markdown="data.text" />
        </div>
    </div>
</template>
