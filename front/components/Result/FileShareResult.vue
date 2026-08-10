<script setup lang="ts">
import { Button } from '@/components/ui/button'
import FilePreviewView from '@/components/FilePreviewView.vue'
import { Input } from '@/components/ui/input'
import { useShare } from '@vueuse/core'
import { useQuery } from '@tanstack/vue-query'
import useMyAppShare from '@/composables/useMyAppShare'
import useMyAppConfig from '@/composables/useMyAppConfig'
import dayjs from 'dayjs'
import showDrawer from '@/lib/showDrawer'
import QrCoreDrawer from '@/components/Drawer/QrCoreDrawer.vue'
import { h } from 'vue'
import type { handleFileComponentProps } from './types'

const props = defineProps<handleFileComponentProps>()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const { t } = useI18n()
const { createFileShare } = useMyAppShare()
const { data } = useQuery({
    queryKey: ['create-share', ...props?.data?.files?.map((item) => item.id)],
    staleTime: Infinity,
    queryFn: async () => {
        const { files, config } = props?.data || {}
        const res = await createFileShare({
            files: files?.map((item) => {
                const { id, file } = item || {}
                return { id, file_name: file.name }
            }),
            config: config as any,
        })
        return res?.data
    },
})

const appConfig = useMyAppConfig()
const getShareUrl = (id: string) => {
    return `${appConfig?.value?.site_url}/s/${id}`
}

const { share, isSupported: isShareSupported } = useShare()

const handleShare = async (id: string) => {
    await share({
        url: getShareUrl(id),
    })
}

const handleShowQrCode = (id: string) => {
    showDrawer({
        render: ({ ...rest }) =>
            h(QrCoreDrawer, {
                ...rest,
                data: getShareUrl(id),
            }),
    })
}
</script>

<template>
    <BaseCard class="flex flex-col gap-3" :title="t('page.result.file.title')" :showBackButton="true">
        <div class="flex flex-col gap-3 items-center">
            <div v-if="props?.data?.files?.length === 1" class="flex flex-col h-30 items-center">
                <FilePreviewView :value="props?.data?.files?.[0]?.file as File" />
            </div>
            <div v-else class="flex flex-col gap-2 w-full p-5 bg-white/20 backdrop-blur-xl rounded-md">
                <div class="text-sm font-semibold">{{ t('page.result.file.fileList') }}</div>
                <div
                    v-for="file in props?.data?.files"
                    :key="file?.id"
                    class="flex flex-row justify-between items-center gap-1 rounded-md p-2 border border-black/10 w-full"
                >
                    <div class="flex flex-row items-center gap-2 flex-1 min-w-0">
                        <FileIcon :file="file?.file as File" size="sm" class="shrink-0" />
                        <div class="text-sm flex-1 truncate">{{ file?.file?.name }}</div>
                    </div>
                </div>
            </div>
            <div v-if="!!data" class="flex flex-col md:flex-row gap-5 rounded-md p-5 bg-white/20 backdrop-blur-xl w-full">
                <div class="flex flex-col gap-2 flex-1">
                    <div class="text-sm font-semibold">{{ t('page.result.file.info') }}</div>
                    <div class="grid grid-cols-2 gap-2">
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">{{ t('page.result.file.downloadNums') }}</div>
                            <div class="text-3xl font-light">{{ data?.download_nums }}</div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">{{ t('page.result.file.expireTime') }}</div>
                            <div class="text-md font-light">
                                {{ dayjs((data?.expire_at || 0) * 1000).format('YYYY-MM-DD HH:mm:ss') }}
                            </div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1" v-if="data?.pickup_code">
                            <div class="flex flex-row justify-between w-full items-center">
                                <div class="text-xs font-semibold">{{ t('page.result.file.pickupCode') }}</div>
                                <CopyButton class="bg-white/70 p-0 size-6" :value="data?.pickup_code as string" />
                            </div>
                            <div class="flex flex-row gap-2">
                                <div v-for="s in data?.pickup_code" class="text-2xl font-light">
                                    {{ s }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex flex-col gap-5 flex-1">
                    <div class="text-sm font-semibold">{{ t('page.result.file.link') }}</div>
                    <div class="flex flex-row gap-2">
                        <Input :model-value="getShareUrl(data?.id as string)" class="bg-white/70" readonly />
                        <Button v-if="isShareSupported" variant="outline" class="bg-white/70" size="icon" @click="handleShare(data?.id as string)">
                            <LucideShare class="size-1/2" />
                        </Button>
                        <CopyButton class="bg-white/70" :value="getShareUrl(data?.id as string)" />

                        <Button variant="outline" class="bg-white/70" size="icon" @click="handleShowQrCode(data?.id as string)">
                            <LucideQrCode class="size-1/2" />
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    </BaseCard>
</template>
