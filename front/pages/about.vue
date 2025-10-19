<script setup lang="ts">
import { CurveType } from '@unovis/ts'
import { AreaChart } from '@/components/ui/chart-area'
import { cx } from 'class-variance-authority'
import { useQuery } from '@tanstack/vue-query'
import { Skeleton } from '@/components/ui/skeleton'
import AboutChartTooltip from '@/components/AboutChartTooltip.vue'
import getFileSize from '~/lib/getFileSize'
import SparkMD5 from 'spark-md5'
import useMyAppConfig from '@/composables/useMyAppConfig'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import Progress from '~/components/ui/progress/Progress.vue'
import AboutBaseInfo from '@/components/About/AboutBaseInfo.vue'

dayjs.extend(relativeTime)
const { t } = useI18n()
</script>

<template>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 my-5 flex flex-col gap-5">
        <div class="text-xl font-normal">{{ t('about.title') }}</div>
        <AboutBaseInfo />
        <AboutChartView />
    </div>
</template>
