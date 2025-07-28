<script setup lang="ts">
import type { NuxtError } from '#app'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'

const { t } = useI18n()

const props = defineProps({
    error: Object as () => NuxtError,
})
console.error(props.error, 'error')
const handleError = () => clearError({ redirect: '/' })
</script>

<template>
    <NuxtLayout>
        <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 flex flex-col items-center justify-center min-h-[50vh] mt-5 gap-10">
            <div class="font-bold text-5xl">{{ error?.statusCode }}</div>
            <Button @click="handleError">{{ t('btn.backToHome') }}</Button>
        </div>
    </NuxtLayout>
</template>
