<script setup lang="ts">
import dayjs from 'dayjs'
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import duration from 'dayjs/plugin/duration'
import relativeTime from 'dayjs/plugin/relativeTime'
import { isBoolean } from 'lodash-es'
import { LucideCheck, LucideX } from '@lucide/vue'
import { cx } from 'class-variance-authority'
import { toast } from 'vue-sonner'
import MarkdownRender from '@/components/MarkdownRender.vue'
import showDrawer from '~/lib/showDrawer'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'

dayjs.extend(duration)
dayjs.extend(relativeTime)

const { t } = useI18n()
const props = defineProps<{
    data: any
}>()

const { getShareToken } = useMyAppShare()

const expireSeconds = computed(() => {
    return dayjs(props?.data?.expire_at * 10e2).unix() - dayjs().unix()
})

const { remaining, start } = useCountdown(expireSeconds.value)

onMounted(() => {
    start()
})

const fileShareInfo = computed(() => {
    return [
        { label: t('page.shareView.textShare.needPassword'), value: props?.data?.has_password ?? false },
        {
            label: t('page.shareView.textShare.expireTime'),
            value: dayjs.duration(remaining.value, 'seconds').format(t('page.shareView.textShare.durationFormat')),
        },
        { label: t('page.shareView.textShare.remainingViews'), value: props?.data?.download_nums ?? 0 },
    ]
})
const previewText = ref<string | null>(null)

const handlePreview = async () => {
    try {
        let token = null
        if (props?.data?.has_password) {
            token = await showDrawer({
                render: ({ ...rest }) => h(PasswallShareDrawer, { ...rest, share_id: props?.data?.id }),
            })
        } else {
            token = await getShareToken(props?.data?.id)
        }
        const r = await $fetch<{
            code: number
            data: {
                data: string
            }
        }>(`/api/download?token=${token}`)
        previewText.value = r?.data?.data
    } catch (error: any) {
        toast.error(error?.data?.message || error?.message || error)
    }
}
</script>
<template>
    <div :class="cx('flex flex-col max-h-full', !!previewText ? 'gap-3' : 'gap-16 items-center')">
        <div :class="cx('flex flex-row w-full', !!previewText ? 'justify-between' : 'justify-center')">
            <h1 class="text-xl">{{ t('page.shareView.textShare.title') }}</h1>
            <CopyButton v-if="!!previewText" :value="previewText as string" />
        </div>
        <template v-if="!previewText">
            <div class="flex flex-col gap-2 md:flex-row w-full">
                <div class="flex flex-row md:flex-col md:gap-1 justify-between items-center md:flex-1" v-for="item in fileShareInfo">
                    <div class="text-xs opacity-75">{{ item?.label }}</div>
                    <component v-if="isBoolean(item?.value)" :is="item?.value ? LucideCheck : LucideX" class="size-6" />
                    <div v-else class="md:text-xl">{{ item?.value }}</div>
                </div>
            </div>
            <div class="w-full">
                <AsyncButton @click="handlePreview" class="w-full">{{ t('page.shareView.textShare.viewBtn') }}</AsyncButton>
            </div>
        </template>
        <template v-else>
            <MarkdownRender :markdown="previewText" class="rounded-md bg-white/70 p-3 w-full max-w-full min-h-80 overflow-y-auto" />
        </template>
    </div>
</template>
