<script setup lang="ts">
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input'
import { useClipboard } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { useQuery } from '@tanstack/vue-query';

const props = defineProps<{
    data: { text: string, config: any, handle_type: string }
}>()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

const { data } = useQuery({
    queryKey: ['create-share', props?.data?.text],
    queryFn: async () => {
        const { config, text } = props?.data || {}
        const data = await $fetch<{
            code: number
            data: {
                id?: string
            }
        }>(`/api/share`, {
            method: 'POST',
            body: {
                type: 'text',
                config,
                data: text,
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
        <div class="flex flex-row justify-between">
            <h2 class="text-lg">分享成功</h2>
            <div class="flex flex-row gap-2 basis-1/2">
                <Button variant="outline" class="bg-white/70" size="icon" @click="() => {
                    emit('change', 'input')
                }">
                    <LucideHome />
                </Button>
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
        </div>
        <div class="prose rounded-md bg-white/70 p-3 w-full max-w-full" v-html="props?.data?.text" />
    </div>

</template>