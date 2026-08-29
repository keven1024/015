<script setup lang="ts">
import { h, type VNode } from 'vue'
import { Skeleton } from '@/components/ui/skeleton'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useQuery } from '@tanstack/vue-query'
import { useShare } from '@vueuse/core'
import useMyAppConfig from '@/composables/useMyAppConfig'
import CopyButton from '~/components/CopyButton.vue'
import QrCoreDrawer from '@/components/Drawer/QrCoreDrawer.vue'
import showDrawer from '@/lib/showDrawer'
import { LucideExternalLink, LucideQrCode, LucideShare } from '@lucide/vue'
import FileShareManage from '@/components/Share/Manage/FileShareManage.vue'
import TextShareManage from '@/components/Share/Manage/TextShareManage.vue'

type ShareFile = { id?: string; file_name: string; size?: number; mime_type?: string }
type ShareInfoItem = {
    label: string
    type: 'countdown' | 'string' | 'bool'
    value: string | number | boolean
    handle?: VNode
}
type ShareData = {
    id?: string
    type?: 'file' | 'text'
    download_nums?: number
    expire_at?: number
    has_password?: boolean
    has_notify?: boolean
    files?: ShareFile[]
    text?: string
    pickup_code?: string
}

const route = useRoute()
const { t } = useI18n()
const appConfig = useMyAppConfig()
const id = computed(() => String(route.params.id || ''))

const { data, isLoading, error, refetch, isFetching } = useQuery({
    queryKey: computed(() => ['share', id.value]),
    queryFn: async () => {
        const response = await $fetch<{ code: number; data: ShareData }>(`/api/share/${id.value}`)
        return response?.data
    },
    retry: false,
})

const shareUrl = computed(() => `${appConfig.value?.site_url || (import.meta.client ? window.location.origin : '')}/s/${id.value}`)
const isFileShare = computed(() => data.value?.type === 'file')
const manageComponentMap = {
    file: FileShareManage,
    text: TextShareManage,
}
const { share, isSupported: isShareSupported } = useShare()
const handleSystemShare = () => share({ title: t('page.shareView.title'), url: shareUrl.value })
const handleShowQrCode = () => {
    showDrawer({
        render: ({ ...rest }) => h(QrCoreDrawer, { ...rest, data: shareUrl.value }),
    })
}

const shareInfo = computed(() => {
    const { pickup_code } = data.value || {}
    return [
        {
            label: isFileShare.value ? t('page.shareManage.downloads') : t('page.shareManage.views'),
            type: 'string' as const,
            value: data.value?.download_nums ?? 0,
        },
        {
            label: t('page.shareManage.expiresAt'),
            type: 'countdown' as const,
            value: data.value?.expire_at ?? 0,
        },
        {
            label: t('page.shareView.needPassword'),
            type: 'bool' as const,
            value: data.value?.has_password ?? false,
        },
        {
            label: t('page.shareManage.notifyEnabled'),
            type: 'bool' as const,
            value: data.value?.has_notify ?? false,
        },
        ...(pickup_code
            ? [
                  {
                      label: t('page.shareManage.pickupCode'),
                      type: 'string' as const,
                      value: pickup_code,
                      handle: h(CopyButton, {
                          value: pickup_code,
                          class: 'size-6 shrink-0 rounded-full',
                      }),
                  },
              ]
            : []),
    ] satisfies ShareInfoItem[]
})
</script>

<template>
    <BaseCard class="flex flex-col gap-3 my-5" :title="t('page.shareManage.title')" :showBackButton="true">
        <div v-if="isLoading" class="space-y-8">
            <div class="space-y-3">
                <Skeleton class="h-4 w-24 rounded-md" />
                <Skeleton class="h-8 w-56 rounded-md" />
            </div>
            <div class="grid gap-3 sm:grid-cols-2">
                <Skeleton v-for="i in 2" :key="i" class="h-28 rounded-xl" />
            </div>
            <Skeleton class="h-20 rounded-xl" />
            <Skeleton class="h-64 rounded-xl" />
        </div>

        <ShareError v-else-if="error || !data" :title="t('page.shareManage.loadFailed')" />

        <div v-else class="flex flex-col gap-3">
            <h2 class="text-sm font-medium">{{ t('page.shareManage.basicInfo') }}</h2>
            <ShareInfoCards :items="shareInfo" />
            <h2 class="text-sm font-medium">{{ t('page.shareManage.shareUrl') }}</h2>
            <div class="flex gap-2">
                <Input :model-value="shareUrl" readonly class="min-w-0 bg-background/60 font-mono text-xs" />
                <Button variant="outline" size="icon" class="shrink-0" as-child>
                    <a
                        :href="shareUrl"
                        target="_blank"
                        rel="noopener noreferrer"
                        :aria-label="t('page.shareManage.openShare')"
                        :title="t('page.shareManage.openShare')"
                    >
                        <LucideExternalLink class="size-4" />
                    </a>
                </Button>
                <Button
                    v-if="isShareSupported"
                    variant="outline"
                    size="icon"
                    class="shrink-0"
                    :aria-label="t('common.systemShare')"
                    :title="t('common.systemShare')"
                    @click="handleSystemShare"
                >
                    <LucideShare class="size-4" />
                </Button>
                <CopyButton class="shrink-0" :value="shareUrl" />
                <Button
                    variant="outline"
                    size="icon"
                    class="shrink-0"
                    :aria-label="t('page.result.qrCode.title')"
                    :title="t('page.result.qrCode.title')"
                    @click="handleShowQrCode"
                >
                    <LucideQrCode class="size-4" />
                </Button>
            </div>
            <component :is="manageComponentMap[data.type as keyof typeof manageComponentMap] || 'div'" :data="data" />
        </div>
    </BaseCard>
</template>
