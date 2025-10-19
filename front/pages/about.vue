<script setup lang="ts">
import AboutBaseInfo from '@/components/About/AboutBaseInfo.vue'
import AboutVersionView from '@/components/About/AboutVersionView.vue'

const { t } = useI18n()
</script>

<template>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 my-5 flex flex-col gap-5">
        <div class="text-xl font-normal">{{ t('about.title') }}</div>
        <AboutBaseInfo />
        <AboutChartView />
        <AboutVersionView />
    </div>
</template>
