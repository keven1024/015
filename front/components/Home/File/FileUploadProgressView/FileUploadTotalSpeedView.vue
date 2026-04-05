<script setup lang="ts">
import dayjs from 'dayjs'
import { motion } from 'motion-v'
import getFileSize from '~/lib/getFileSize'
import showDrawer from '~/lib/showDrawer'
import FileUploadSpeedInfoView from './FileUploadSpeedInfoView.vue'
import { clamp } from 'lodash-es'

const { t } = useI18n()

const props = defineProps<{
    speedChartData: Record<string, { fileId: string; index: number; value: number }[]>
}>()
const speedChartList = ref<{ timestamp: number; value: number }[]>([])
useIntervalFn(() => {
    speedChartList.value.push({
        timestamp: dayjs().unix() - 1,
        value: props.speedChartData[dayjs().unix() - 1]?.reduce((acc, curr) => acc + curr.value, 0) ?? 0,
    })
    if (speedChartList.value.length > 60) {
        speedChartList.value.shift()
    }
}, 1000)

const handleShowSpeedInfo = () => {
    showDrawer({
        render: ({ hide }) => h(FileUploadSpeedInfoView, { hide }),
    })
}
</script>
<template>
    <div class="rounded-xl p-3 bg-white/80 flex flex-col gap-2 col-span-4 md:col-span-3 h-32 md:h-auto">
        <div class="flex flex-col gap-1">
            <div @click="handleShowSpeedInfo" class="flex flex-row gap-1 items-center text-xs opacity-70">
                {{ t('page.progress.file.totalUploadProgress') }}
                <LucideInfo class="size-3" />
            </div>
            <div class="text-2xl font-bold">
                {{ `${getFileSize(speedChartData[dayjs().unix() - 1]?.reduce((acc, curr) => acc + curr.value, 0) ?? 0) ?? 0}/s` }}
            </div>
        </div>
        <div class="flex-1 relative overflow-hidden flex flex-row gap-0.5 justify-end items-end">
            <motion.div
                class="w-2 shrink-0 bg-primary relative"
                :style="{
                    height: `${clamp((i.value / Math.max(...(speedChartList?.map((r) => r.value) || [1]))) * 100, 1, 100)}%`,
                }"
                :layoutId="String(i.timestamp)"
                v-for="i in speedChartList"
                :key="i.timestamp"
                :initial="{ x: 10, opacity: 0 }"
                :animate="{ x: 0, opacity: 1 }"
                :transition="{ duration: 1 }"
            >
            </motion.div>
            <!-- <BarChart class="h-full" :data="data" index="time" :categories="['value']" :showTooltip="false"
            :showLegend="false" :showXAxis="false" :showYAxis="false" :showGrid="false" :groupMaxWidth="10" /> -->
        </div>
    </div>
</template>
