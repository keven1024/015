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
import { cx } from 'class-variance-authority'
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
        const data = await createFileShare({
            files: files?.map((item) => {
                const { id, file } = item || {}
                return { id, name: file.name }
            }),
            config: config as any,
        })
        return data?.map((item) => item?.data)
    },
})
const selectedFile = ref<string | undefined>()
const selectedFileShare = computed(() => {
    return data?.value?.find((item) => item?.id === selectedFile.value)
})
watchEffect(() => {
    if (data?.value && data?.value?.length === 1 && !!data?.value?.[0]?.id) {
        selectedFile.value = data.value[0].id
    }
})

const appConfig = useMyAppConfig()
const getShareUrl = (id: string) => {
    return `${appConfig?.value?.site_url}/s/${id}`
}

const { share, isSupported: isShareSupported } = useShare()

const handleShare = async (id: string, fileName?: string) => {
    await share({
        title: fileName || 'File Share',
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
            <div v-if="data?.length === 1" class="flex flex-col h-30 items-center">
                <FilePreviewView :value="props?.data?.files?.[0]?.file as File" />
            </div>
            <div v-else class="flex flex-col gap-2 w-full p-5 bg-white/20 backdrop-blur-xl rounded-md">
                <div class="text-sm font-semibold">{{ t('page.result.file.fileList') }}</div>
                <div
                    v-for="file in data"
                    :class="
                        cx(
                            'flex flex-row justify-between items-center gap-1 rounded-md p-2 border border-black/10 w-full cursor-pointer',
                            selectedFile === file?.id && 'bg-primary text-white'
                        )
                    "
                    @click="selectedFile = file?.id"
                >
                    <div class="flex flex-row items-center gap-2 flex-1 min-w-0">
                        <FileIcon
                            :file="props?.data?.files?.[data?.findIndex((i) => i?.id === file?.id) as number]?.file as File"
                            size="sm"
                            :class="cx('shrink-0', selectedFile === file?.id && 'bg-white/50!')"
                        />
                        <div class="text-sm flex-1 truncate">{{ file?.file_name }}</div>
                    </div>
                    <div class="flex flex-row items-center gap-2 shrink-0">
                        <Button
                            v-if="isShareSupported"
                            variant="outline"
                            :class="cx('bg-white/70', selectedFile === file?.id && '!bg-white/30 border-none hover:text-white/80')"
                            size="icon"
                            @click.stop="handleShare(file?.id as string, file?.file_name)"
                        >
                            <LucideShare class="size-1/2" />
                        </Button>
                        <CopyButton
                            :class="cx('bg-white/70', selectedFile === file?.id && '!bg-white/30 border-none hover:text-white/80')"
                            :value="getShareUrl(file?.id as string)"
                            @click.stop
                        />
                        <Button
                            variant="outline"
                            :class="cx('bg-white/70', selectedFile === file?.id && '!bg-white/30 border-none hover:text-white/80')"
                            size="icon"
                            @click.stop="handleShowQrCode(file?.id as string)"
                        >
                            <LucideQrCode class="size-1/2" />
                        </Button>
                    </div>
                </div>
            </div>
            <div v-if="!!selectedFileShare" class="flex flex-col md:flex-row gap-5 rounded-md p-5 bg-white/20 backdrop-blur-xl w-full">
                <div class="flex flex-col gap-2 flex-1">
                    <div class="text-sm font-semibold">{{ t('page.result.file.info') }}</div>
                    <div class="grid grid-cols-2 gap-2">
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">{{ t('page.result.file.downloadNums') }}</div>
                            <div class="text-3xl font-light">{{ selectedFileShare?.download_nums }}</div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">{{ t('page.result.file.expireTime') }}</div>
                            <div class="text-md font-light">
                                {{ dayjs((selectedFileShare?.expire_at || 0) * 1000).format('YYYY-MM-DD HH:mm:ss') }}
                            </div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1" v-if="selectedFileShare?.pickup_code">
                            <div class="flex flex-row justify-between w-full items-center">
                                <div class="text-xs font-semibold">{{ t('page.result.file.pickupCode') }}</div>
                                <CopyButton class="bg-white/70 p-0 size-6" :value="selectedFileShare?.pickup_code as string" />
                            </div>
                            <div class="flex flex-row gap-2">
                                <div v-for="s in selectedFileShare?.pickup_code" class="text-2xl font-light">
                                    {{ s }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex flex-col gap-5 flex-1">
                    <div class="text-sm font-semibold">{{ t('page.result.file.link') }}</div>
                    <div class="flex flex-row gap-2">
                        <Input :model-value="getShareUrl(selectedFileShare?.id as string)" class="bg-white/70" readonly />
                        <Button
                            v-if="isShareSupported"
                            variant="outline"
                            class="bg-white/70"
                            size="icon"
                            @click="
                                handleShare(
                                    selectedFileShare?.id as string,
                                    props?.data?.files?.[data?.findIndex((item) => item?.id === selectedFileShare?.id) as number]?.file?.name
                                )
                            "
                        >
                            <LucideShare class="size-1/2" />
                        </Button>
                        <CopyButton class="bg-white/70" :value="getShareUrl(selectedFileShare?.id as string)" />

                        <Button variant="outline" class="bg-white/70" size="icon" @click="handleShowQrCode(selectedFileShare?.id as string)">
                            <LucideQrCode class="size-1/2" />
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    </BaseCard>
</template>
