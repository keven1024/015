<script setup lang="ts">
import { LucideLink2Off } from '@lucide/vue'
import { Button } from '@/components/ui/button'

defineProps<{
    title: string
    description?: string
}>()

const router = useRouter()
const { t } = useI18n()
</script>

<template>
    <div class="flex min-h-96 flex-col items-center justify-center gap-5 text-center">
        <LucideLink2Off aria-hidden="true" class="size-12 text-red-500 dark:text-red-400" :stroke-width="1.5" />
        <div class="space-y-1">
            <h1 class="text-xl font-medium">{{ title }}</h1>
            <p v-if="description" class="text-sm text-muted-foreground">{{ description }}</p>
        </div>
        <div class="flex gap-2">
            <Button @click="router.push('/')">{{ t('btn.backToHome') }}</Button>
        </div>
    </div>
</template>
