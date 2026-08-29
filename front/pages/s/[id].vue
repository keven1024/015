<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
import dayjs from 'dayjs'
import FileShareView from '@/components/Share/FileShareView.vue'
import TextShareView from '@/components/Share/TextShareView.vue'
import { useQuery } from '@tanstack/vue-query'

type ShareData = {
    id?: string
    expire_at?: number
    type?: 'file' | 'text'
    has_password?: boolean
    download_nums?: number
}

const route = useRoute()
const { t } = useI18n()
const id = computed(() => route.params.id)

const { data, isLoading } = useQuery({
    queryKey: ['share', id.value],
    queryFn: async () => {
        const data = await $fetch<{ code: number; data: ShareData }>(`/api/share/${id.value}`)
        return data?.data
    },
    retry: false,
})

const isExpired = computed(() => {
    const { expire_at } = data.value || {}
    return !data || !expire_at || dayjs(expire_at * 10e2).isBefore(dayjs())
})

const componentMap = {
    file: FileShareView,
    text: TextShareView,
}

const shareInfo = computed(() => {
    const type = data.value?.type
    if (!type) {
        return []
    }

    const key = `page.shareView.${type}Share`
    return [
        { label: t('page.shareView.needPassword'), type: 'bool' as const, value: data.value?.has_password ?? false },
        { label: t('page.shareView.expireTime'), type: 'countdown' as const, value: data.value?.expire_at ?? 0 },
        {
            label: t(`${key}.${type === 'file' ? 'remainingDownloads' : 'remainingViews'}`),
            type: 'string' as const,
            value: data.value?.download_nums ?? 0,
        },
    ]
})
</script>

<template>
    <BaseCard class="my-5 overflow-hidden">
        <div v-if="isLoading" class="flex flex-col items-center gap-5">
            <Skeleton class="h-6 w-32 rounded-md" />
            <Skeleton class="size-16 rounded-xl" />
            <Skeleton class="h-5 w-28 rounded-md" />
            <div class="flex w-full flex-col gap-2 md:flex-row">
                <div
                    v-for="i in 3"
                    :key="i"
                    class="flex min-h-11 flex-1 items-center justify-between gap-1 rounded-xl bg-black/5 px-3 py-2 md:flex-col md:items-start md:justify-between"
                >
                    <Skeleton class="h-3 w-16 rounded-md bg-black/10" />
                    <Skeleton v-if="i === 1" class="size-7 rounded-full bg-white/50" />
                    <Skeleton v-else class="h-6 w-16 rounded-md bg-black/10" />
                </div>
            </div>
            <Skeleton class="h-10 w-full rounded-md" />
        </div>
        <template v-else>
            <ShareError v-if="isExpired || !data" :title="t('page.shareView.linkExpired')" />
            <div v-else class="flex flex-col items-center gap-3">
                <h1 class="text-xl">{{ t('page.shareView.title') }}</h1>
                <ShareInfoCards :items="shareInfo" />
                <component :is="componentMap[data?.type as keyof typeof componentMap] || 'div'" :data="data" />
            </div>
        </template>
    </BaseCard>
</template>
