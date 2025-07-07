<script setup lang="ts">
import dayjs from 'dayjs'
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import duration from 'dayjs/plugin/duration'
import relativeTime from 'dayjs/plugin/relativeTime'
import { isBoolean } from 'lodash-es'
import { LucideCheck, LucideX } from 'lucide-vue-next'
import { cx } from 'class-variance-authority'
import { toast } from 'vue-sonner'
import MarkdownRender from '@/components/MarkdownRender.vue'
import { Button } from '@/components/ui/button'
import { LucideCopy } from 'lucide-vue-next'
import { useClipboard } from '@vueuse/core'
import showDrawer from '~/lib/showDrawer'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'

dayjs.extend(duration)
dayjs.extend(relativeTime)

const props = defineProps<{
    data: any
}>()

const { getShareToken } = useMyAppShare()

const expireSeconds = computed(() => {
    return dayjs(props?.data?.expire_at * 10e2).unix() - dayjs().unix()
})

const { remaining, start } = useCountdown(expireSeconds.value)

const { copy } = useClipboard()

onMounted(() => {
    start()
})

const fileShareInfo = computed(() => {
    return [
        { label: '需要密码', value: props?.data?.has_password ?? false },
        {
            label: '过期时间',
            value: dayjs.duration(remaining.value, 'seconds').format(`D天 HH:mm:ss`),
        },
        { label: '剩余浏览次数', value: props?.data?.download_nums ?? 0 },
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
            <h1 class="text-xl">查看文本</h1>
            <Button
                v-if="!!previewText"
                variant="outline"
                size="icon"
                @click="
                    () => {
                        copy(previewText as string)
                        toast.success('复制成功')
                    }
                "
            >
                <LucideCopy />
            </Button>
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
                <AsyncButton @click="handlePreview" class="w-full">浏览</AsyncButton>
            </div>
        </template>
        <template v-else>
            <MarkdownRender :markdown="previewText" class="rounded-md bg-white/70 p-3 w-full max-w-full min-h-80 overflow-y-auto" />
        </template>
    </div>
</template>
