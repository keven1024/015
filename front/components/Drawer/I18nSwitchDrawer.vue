<script setup lang="ts">
import { cx } from 'class-variance-authority'
import type { Locale } from '@intlify/core-base'

const props = defineProps<{
    hide: () => void
}>()

const { locales, setLocale, locale: currentLocale, t } = useI18n()

const switchLocale = async (locale: Locale) => {
    await setLocale(locale)
    props.hide()
}
</script>

<template>
    <div class="flex flex-col gap-1 py-2">
        <div class="text-xl font-bold mb-3">{{ t('i18n.switchLocale') }}</div>
        <div
            v-for="locale in locales"
            :key="locale.code"
            :class="cx('rounded-md hover:bg-black/10 p-2 cursor-pointer', currentLocale === locale.code && 'bg-black/10 font-bold')"
            @click="() => switchLocale(locale.code)"
        >
            {{ locale.name }}
        </div>
    </div>
</template>
