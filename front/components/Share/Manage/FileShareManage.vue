<script setup lang="ts">
type ShareFile = {
    id?: string
    file_name: string
    size?: number
    mime_type?: string
}

defineProps<{
    data: {
        files?: ShareFile[]
    }
}>()

const { t } = useI18n()
</script>

<template>
    <div v-if="data.files?.length" class="flex flex-col gap-3">
        <h2 class="text-sm font-medium">{{ t('page.shareManage.files') }}</h2>
        <ShareFileInfoList :files="data.files" />
    </div>
</template>
