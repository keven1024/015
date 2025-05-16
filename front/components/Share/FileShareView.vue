<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import { isBoolean } from 'lodash-es';
import { LucideCheck, LucideX } from 'lucide-vue-next';
import { useQueryClient } from '@tanstack/vue-query';
dayjs.extend(duration)
dayjs.extend(relativeTime)

const props = defineProps<{
    data: any
}>()

const queryClient = useQueryClient()

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
    queryClient.invalidateQueries({ queryKey: ['share', id] })
}

const expireSeconds = computed(() => {
    return dayjs(props?.data?.expire_at * 10e2).unix() - dayjs().unix()
})

const { remaining, start} = useCountdown(expireSeconds.value)

onMounted(() => {
    start()
})

const fileShareInfo = computed(() => {
    return [
        { label: '需要密码', value: props?.data?.has_password ?? false },
        { label: '过期时间', value: dayjs.duration(remaining.value, 'seconds').format(`D天 HH:mm:ss`) },
        { label: '剩余下载次数', value: props?.data?.download_nums ?? 0 },
    ]
})
</script>

<template>
    <div class="flex flex-col gap-5 items-center">
        <h1 class="text-xl font-bold">下载文件</h1>
        <FilePreviewView :value="props?.data" />
        <div class="flex flex-col gap-2 md:flex-row w-full">
            <div class="flex flex-row md:flex-col md:gap-1 justify-between items-center md:flex-1" v-for="item in fileShareInfo">
                <div class="text-xs opacity-75">{{ item?.label }}</div>
                <component v-if="isBoolean(item?.value)" :is="item?.value ? LucideCheck : LucideX" class="size-6" />
                <div v-else class="md:text-xl">{{ item?.value }}</div>
            </div>
        </div>
        <div class="w-full">
            <AsyncButton @click="handleDownload" class="w-full">下载</AsyncButton>
        </div>
    </div>
</template>