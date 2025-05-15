<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
const props = defineProps<{
    data: any
}>()

const handleDownload = async () => {
    const { id } = props?.data || {}
    const data = await $fetch<{
        code: number
        data: {
            token?: string
        }
    }>(`/api/download`, {
        method: 'POST',
        body: {
            share_id: id
        }
    })
    const { token } = data?.data || {}
    if (!token) {
        return
    }
    (window as any)?.open(`/api/download?token=${token}`)
}

</script>

<template>
    <div class="flex flex-col gap-5 items-center">
        <h1>下载文件</h1>
        <FilePreviewView :value="props?.data" />
        <div>
            <AsyncButton @click="handleDownload">下载</AsyncButton>
        </div>
    </div>
</template>