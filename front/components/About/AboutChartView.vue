<script setup lang="ts">
import { cx } from 'class-variance-authority'
import { useQuery } from '@tanstack/vue-query'
import { Skeleton } from '@/components/ui/skeleton'
import dayjs from 'dayjs'
import { filesize } from 'filesize'
import type { ChartConfig } from '@/components/ui/chart'
import { VisArea, VisAxis, VisLine, VisXYContainer } from '@unovis/vue'
import { ChartContainer, ChartTooltip, ChartCrosshair, ChartLegendContent, componentToString, ChartTooltipContent } from '@/components/ui/chart'
import DateRangePicker from '@/components/DateRangePicker.vue'
import type { DateRange } from '@/components/DateRangePicker.vue'
import { times, uniq } from 'lodash-es'

interface BaseStatChartData {
    count: number
    size?: number
}

type StatChartType = 'file' | 'share' | 'download' | 'worker'

interface StatApiResponse<T> {
    data: T
}

type ChartDataItem = {
    date: Date
    [key: string]: any
}

type AreaChartConfig = {
    data: ChartDataItem[]
    index: string
    config: ChartConfig
}

function calcDateRange(days: number): DateRange {
    return {
        start: dayjs()
            .subtract(days - 1, 'day')
            .startOf('hour')
            .add(1, 'hour')
            .toDate(),
        end: dayjs().startOf('hour').add(1, 'hour').toDate(),
    }
}

const dateRangeOptions = [
    { label: '最近1天', value: calcDateRange(1) },
    { label: '最近7天', value: calcDateRange(7) },
    { label: '最近30天', value: calcDateRange(30) },
]

const dateRange = ref<DateRange>(calcDateRange(7))
const { data, isLoading } = useQuery({
    queryKey: computed(() => ['stat', dateRange?.value?.start?.getTime(), dateRange?.value?.end?.getTime()]),
    queryFn: async () => {
        const fetchStat = async <T,>(type: StatChartType) => {
            const response = await $fetch<StatApiResponse<T>>('/api/stat', {
                method: 'POST',
                body: {
                    type,
                    start: dateRange.value.start ? Math.floor(dateRange.value.start.getTime() / 1000) : undefined,
                    end: dateRange.value.end ? Math.floor(dateRange.value.end.getTime() / 1000) : undefined,
                },
            })
            return response.data
        }

        const [file, share, download, workerResponse] = await Promise.all([
            fetchStat<Record<string, BaseStatChartData>>('file'),
            fetchStat<Record<string, BaseStatChartData>>('share'),
            fetchStat<Record<string, BaseStatChartData>>('download'),
            fetchStat<Record<string, Record<string, number>>>('worker'),
        ])

        return { file, share, download, worker: workerResponse }
    },
})

const { t } = useI18n()

const chartTabs = computed(() => {
    return [
        {
            label: t('page.about.file'),
            value: 'file',
            total: Object.values(data.value?.file || {}).reduce((acc, curr) => acc + curr.count, 0),
        },
        {
            label: t('page.about.share'),
            value: 'share',
            total: Object.values(data.value?.share || {}).reduce((acc, curr) => acc + curr.count, 0),
        },
        {
            label: t('page.about.download'),
            value: 'download',
            total: Object.values(data.value?.download || {}).reduce((acc, curr) => acc + curr.count, 0),
        },
        {
            label: t('page.about.task'),
            value: 'worker',
            total: Object.values(data.value?.worker || {}).reduce(
                (sum, taskGroups) => sum + Object.values(taskGroups).reduce((count, taskCount) => count + taskCount, 0),
                0
            ),
        },
    ]
})

const currentChartTab = ref<StatChartType>('file')
const currentChartData = computed((): AreaChartConfig => {
    const diffHours = dayjs(dateRange.value.end).diff(dayjs(dateRange.value.start), 'hour')

    if (['worker']?.includes(currentChartTab.value)) {
        const workerData = data.value?.worker ?? {}
        const taskTypes = uniq(Object.values(workerData).flatMap((item) => Object.keys(item)))
        const config: ChartConfig = Object.fromEntries(
            taskTypes.map((item, index) => [item, { color: ['#4ade80', '#f87171', '#38bdf8', '#a78bfa', '#ea580c'][index % 5], label: item }])
        )

        return {
            data: times(diffHours, (index) => {
                const date = dayjs(dateRange.value.end).subtract(index, 'hour')
                const item = workerData[date.unix()] ?? {}

                return {
                    date: date.toDate(),
                    ...Object.fromEntries(taskTypes.map((taskType) => [taskType, item[taskType] ?? 0])),
                }
            }),
            index: 'date' as const,
            config,
        }
    }

    if (['share', 'download']?.includes(currentChartTab.value)) {
        return {
            data: times(diffHours, (index) => {
                const date = dayjs(dateRange.value.end).subtract(index, 'hour')
                const { count = 0, size } = data?.value?.[currentChartTab.value]?.[date?.unix()] || {}
                return { date: date.toDate(), count }
            }),
            index: 'date' as const,
            config: {
                count:
                    currentChartTab.value === 'share'
                        ? { color: '#ea580c', label: t('page.about.share') }
                        : { color: '#a3e635', label: t('page.about.download') },
            },
        }
    }
    return {
        data: times(diffHours, (index) => {
            const date = dayjs(dateRange.value.end).subtract(index, 'hour')
            const { count = 0, size } = data?.value?.[currentChartTab.value]?.[date?.unix()] || {}
            return { date: date.toDate(), count, size }
        }),
        index: 'date' as const,
        config: {
            size: { color: '#38bdf8', label: t('page.about.fileSize') },
            count: { color: '#a78bfa', label: t('page.about.fileNum') },
        },
    }
})
</script>

<template>
    <div class="flex items-center justify-between gap-3">
        <div class="font-semibold">{{ t('page.about.analysis') }}</div>
        <DateRangePicker v-model="dateRange" :options="dateRangeOptions" class="bg-white/50" />
    </div>
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
                            currentChartTab = tab.value as 'file' | 'worker' | 'share' | 'download'
                        }
                    "
                >
                    <div class="opacity-75 text-xs">{{ tab.label }}</div>
                    <div class="text-lg font-semibold">{{ tab.total }}</div>
                </div>
            </div>
            <ChartContainer :config="currentChartData.config" class="h-64 w-full p-5" :cursor="false">
                <VisXYContainer :data="currentChartData.data" :x-domain="[dateRange.end, dateRange.start]">
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
                        :num-ticks="12"
                        :tick-format="
                            (d: Date) => {
                                return dayjs(d).format('MM-DD HH:00')
                            }
                        "
                    />
                    <ChartTooltip />
                    <ChartCrosshair
                        :key="currentChartTab"
                        :template="
                            componentToString(currentChartData.config, ChartTooltipContent, {
                                class: 'w-[14rem]',
                                labelFormatter: (d) => {
                                    return dayjs(d).format('MM-DD HH:00')
                                },
                                valueFormatter: (value, key) => {
                                    if (key === 'size' && typeof value === 'number') {
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
