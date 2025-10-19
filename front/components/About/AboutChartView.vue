<script setup lang="ts">
import { CurveType } from '@unovis/ts'
import { AreaChart } from '@/components/ui/chart-area'
import { cx } from 'class-variance-authority'
import { useQuery } from '@tanstack/vue-query'
import { Skeleton } from '@/components/ui/skeleton'
import AboutChartTooltip from '@/components/AboutChartTooltip.vue'
import dayjs from 'dayjs'
import { times } from 'lodash-es'

interface StatChartData {
    file_size: number
    file_num: number
    share_num: number
    download_num: number
    date: string
}

interface QueueChartData {
    processed: number
    failed: number
    date: string
}

type ChartDataItem = StatChartData | QueueChartData

type ChartConfig = {
    data: ChartDataItem[]
    index: string
    categories: string[]
    colors: string[]
}

const { data, isLoading } = useQuery({
    queryKey: ['stat'],
    queryFn: async () => {
        const response = await $fetch<{
            data: {
                chart: {
                    storage: Record<string, StatChartData>
                    queue: Record<string, QueueChartData>
                }
            }
        }>('/api/stat')
        return response.data
    },
})

const { t } = useI18n()

const chartTabs = computed(() => {
    return [
        {
            label: t('about.file'),
            value: 'storage',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.file_num, 0)
                : 0,
        },
        {
            label: t('about.share'),
            value: 'share',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.share_num, 0)
                : 0,
        },
        {
            label: t('about.download'),
            value: 'download',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.download_num, 0)
                : 0,
        },
        {
            label: t('about.task'),
            value: 'queue',
            total: data.value?.chart?.queue
                ? Object.values(data.value.chart.queue).reduce((acc: number, curr: QueueChartData) => acc + curr.processed + curr.failed, 0)
                : 0,
        },
    ]
})

const currentChartTab = ref<'storage' | 'queue' | 'share' | 'download'>('storage')
const currentChartData = computed((): ChartConfig => {
    const { storage, queue } = data.value?.chart || {}
    if (currentChartTab.value === 'queue') {
        const queueData = times(30, (i) => {
            return {
                date: dayjs().subtract(i, 'day').format('YYYY-MM-DD'),
                processed: queue?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.processed || 0,
                failed: queue?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.failed || 0,
            }
        })
        return {
            data: queueData,
            index: 'date' as const,
            categories: ['processed', 'failed'] as const,
            colors: ['#4ade80', '#f87171'],
        }
    }
    const storageData = times(30, (i) => {
        const base = { date: dayjs().subtract(i, 'day').format('YYYY-MM-DD') }
        if (currentChartTab.value === 'share') {
            return {
                ...base,
                share_num: storage?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.share_num || 0,
            }
        }
        if (currentChartTab.value === 'download') {
            return {
                ...base,
                download_num: storage?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.download_num || 0,
            }
        }
        return {
            ...base,
            file_size: storage?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.file_size || 0,
            file_num: storage?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.file_num || 0,
        }
    })

    let categories = ['file_size', 'file_num']
    if (currentChartTab.value === 'share') {
        categories = ['share_num']
    }
    if (currentChartTab.value === 'download') {
        categories = ['download_num']
    }
    let colors = ['#38bdf8', '#a78bfa']
    if (currentChartTab.value === 'share') {
        colors = ['#ea580c']
    }
    if (currentChartTab.value === 'download') {
        colors = ['#a3e635']
    }
    return {
        data: storageData as ChartDataItem[],
        index: 'date' as const,
        categories,
        colors,
    }
})
</script>

<template>
    <div class="font-semibold">{{ t('about.analysis') }}</div>
    <template v-if="isLoading">
        <div class="flex flex-row gap-2">
            <Skeleton class="w-full h-96 rounded-xl" />
        </div>
    </template>
    <template v-else>
        <div class="flex flex-col gap-2 bg-white/50 w-full rounded-xl py-5">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-2 px-5">
                <div
                    :class="cx('rounded-md min-w-30 flex flex-col px-3 py-1.5 cursor-pointer', currentChartTab === tab.value && 'bg-black/10')"
                    v-for="tab in chartTabs"
                    :key="tab.value"
                    @click="
                        () => {
                            currentChartTab = tab.value as 'storage' | 'queue'
                        }
                    "
                >
                    <div class="opacity-75 text-xs">{{ tab.label }}</div>
                    <div class="text-lg font-semibold">{{ tab.total }}</div>
                </div>
            </div>
            <AreaChart
                v-if="currentChartData"
                class="h-64 w-full"
                :key="currentChartTab"
                :index="currentChartData.index"
                :data="currentChartData.data"
                :categories="currentChartData.categories"
                :show-grid-line="false"
                :show-legend="false"
                :show-y-axis="true"
                :show-x-axis="true"
                :colors="currentChartData.colors"
                :custom-tooltip="AboutChartTooltip"
                :curve-type="CurveType.CatmullRom"
            />
        </div>
    </template>
</template>
