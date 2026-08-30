<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import { LucideTrash2 } from '@lucide/vue'

const props = defineProps<{
    title: string
    desc: string
    btnLabel: string
    onClick: () => Promise<void>
}>()
</script>

<template>
    <div class="flex flex-col gap-5">
        <div class="space-y-2">
            <h2 class="text-xl font-bold">{{ title }}</h2>
            <p class="text-sm text-muted-foreground">{{ desc }}</p>
        </div>
        <AsyncButton variant="destructive" @click="props.onClick">
            {{ btnLabel }}
        </AsyncButton>
    </div>
</template>
