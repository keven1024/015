<script setup lang="ts">
import { cx } from 'class-variance-authority'
import { useQuery } from '@tanstack/vue-query'
import { Skeleton } from '@/components/ui/skeleton'
import dayjs from 'dayjs'
import { filesize } from 'filesize'
import { times } from 'lodash-es'
import type { ChartConfig } from '@/components/ui/chart'
import { VisArea, VisAxis, VisLine, VisXYContainer } from '@unovis/vue'
import { ChartContainer, ChartTooltip, ChartCrosshair, ChartLegendContent, componentToString, ChartTooltipContent } from '@/components/ui/chart'

interface StatChartData {
    file_size: number
    file_num: number
    share_num: number
    download_num: number
    date: Date
}

interface QueueChartData {
    processed: number
    failed: number
    date: Date
}

type ChartDataItem = StatChartData | QueueChartData

type AreaChartConfig = {
    data: ChartDataItem[]
    index: string
    config: ChartConfig
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
            label: t('page.about.file'),
            value: 'storage',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.file_num, 0)
                : 0,
        },
        {
            label: t('page.about.share'),
            value: 'share',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.share_num, 0)
                : 0,
        },
        {
            label: t('page.about.download'),
            value: 'download',
            total: data.value?.chart?.storage
                ? Object.values(data.value.chart.storage).reduce((acc: number, curr: StatChartData) => acc + curr.download_num, 0)
                : 0,
        },
        {
            label: t('page.about.task'),
            value: 'queue',
            total: data.value?.chart?.queue
                ? Object.values(data.value.chart.queue).reduce((acc: number, curr: QueueChartData) => acc + curr.processed + curr.failed, 0)
                : 0,
        },
    ]
})

const currentChartTab = ref<'storage' | 'queue' | 'share' | 'download'>('storage')
const currentChartData = computed((): AreaChartConfig => {
    const { storage, queue } = data.value?.chart || {}
    if (currentChartTab.value === 'queue') {
        const queueData = times(30, (i) => {
            return {
                date: dayjs().subtract(i, 'day').toDate(),
                processed: queue?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.processed || 0,
                failed: queue?.[dayjs().subtract(i, 'day').format('YYYY-MM-DD')]?.failed || 0,
            }
        })
        return {
            data: queueData,
            index: 'date' as const,
            config: {
                processed: { color: '#4ade80', label: t('page.about.processed') },
                failed: { color: '#f87171', label: t('page.about.failed') },
            },
        }
    }
    const storageData = times(30, (i) => {
        const base = { date: dayjs().subtract(i, 'day').toDate() }
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

    if (currentChartTab.value === 'share') {
        return {
            data: storageData as ChartDataItem[],
            index: 'date' as const,
            config: {
                share_num: { color: '#ea580c', label: t('page.about.share') },
            },
        }
    }
    if (currentChartTab.value === 'download') {
        return {
            data: storageData as ChartDataItem[],
            index: 'date' as const,
            config: {
                download_num: { color: '#a3e635', label: t('page.about.download') },
            },
        }
    }
    return {
        data: storageData as ChartDataItem[],
        index: 'date' as const,
        config: {
            file_size: { color: '#38bdf8', label: t('page.about.fileSize') },
            file_num: { color: '#a78bfa', label: t('page.about.fileNum') },
        },
    }
})
</script>

<template>
    <div class="font-semibold">{{ t('page.about.analysis') }}</div>
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
            <ChartContainer :config="currentChartData.config" class="h-64 w-full p-5" :cursor="false">
                <VisXYContainer :data="currentChartData.data" :x-domain="[dayjs().toDate(), dayjs().subtract(29, 'day').toDate()]">
                    <VisArea
                        :key="currentChartTab"
                        :x="(d: ChartDataItem) => d.date"
                        :y="Object.keys(currentChartData.config).map((key) => (d: ChartDataItem) => d?.[key as keyof ChartDataItem])"
                        :color="Object.values(currentChartData.config).map((c) => c.color)"
                        :opacity="0.6"
                    />
                    <VisLine
                        :key="currentChartTab"
                        :x="(d: ChartDataItem) => d.date"
                        :y="Object.keys(currentChartData.config).map((key) => (d: ChartDataItem) => d?.[key as keyof ChartDataItem])"
                        :color="Object.values(currentChartData.config).map((c) => c.color)"
                        :line-width="1"
                    />
                    <VisAxis
                        :key="currentChartTab"
                        type="x"
                        :tick-line="false"
                        :domain-line="false"
                        :grid-line="false"
                        :num-ticks="6"
                        :tick-format="
                            (d: Date) => {
                                return dayjs(d).format('MMM')
                            }
                        "
                        :tick-values="currentChartData.data.map((d) => d.date)"
                    />
                    <ChartTooltip />
                    <ChartCrosshair
                        :key="currentChartTab"
                        :template="
                            componentToString(currentChartData.config, ChartTooltipContent, {
                                class: 'w-[14rem]',
                                labelFormatter: (d) => {
                                    return dayjs(d).format('MM-DD')
                                },
                                valueFormatter: (value, key) => {
                                    if (key === 'file_size' && typeof value === 'number') {
                                        return filesize(value)
                                    }
                                    return String(value)
                                },
                            })
                        "
                        :color="(d: any, i: number) => Object.values(currentChartData.config).map((c) => c.color as string)[i]"
                    />
                </VisXYContainer>
                <ChartLegendContent />
            </ChartContainer>
        </div>
    </template>
</template>
