<script setup lang="ts">
import { Button } from '@/components/ui/button';
import FilePreviewView from '@/components/FilePreviewView.vue';
import { Input } from '@/components/ui/input'
import { useClipboard } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { useQuery } from '@tanstack/vue-query';
const props = defineProps<{
    data: { file: File, config: any, handle_type: string, file_id: string }
}>()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const { data } = useQuery({
    queryKey: ['create-share', props?.data?.file_id],
    queryFn: async () => {
        const { file_id, config, file } = props?.data || {}
        const { name } = file || {}
        const data = await $fetch<{
            code: number
            data: {
                id?: string
            }
        }>(`/api/share`, {
            method: 'POST',
            body: {
                type: 'file',
                config,
                data: file_id,
                file_name: name
            }
        })
        return data?.data
    }
})

const appConfig = useAppConfig()
const url = computed(() => {
    const { id } = data?.value || {}
    return `${appConfig?.value?.site_url}/s/${id}`
})

const { copy } = useClipboard()
</script>

<template>

    <div class="flex flex-col gap-3">
        <h2 class="text-lg">上传成功</h2>
        <div class="flex flex-col gap-3 items-center">
            <div class="flex flex-col h-30 items-center">
                <FilePreviewView :value="props?.data?.file" />
            </div>
            <div class="flex flex-row gap-2">
                <Input v-model="url" class="bg-white/70" />
                <Button variant="outline" class="bg-white/70" size="icon" @click="() => {
                    copy(url)
                    toast.success('复制成功')
                }">
                    <LucideCopy />
                </Button>

                <Button variant="outline" class="bg-white/70" size="icon">
                    <LucideQrCode />
                </Button>
            </div>
            <Button variant="ghost" class="hover:bg-white/50 w-40" @click="() => {
                emit('change', 'input')
            }">
                确定
            </Button>
        </div>
    </div>

</template>