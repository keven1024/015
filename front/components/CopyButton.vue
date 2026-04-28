<script setup lang="ts">
import { Button } from '@/components/ui/button'
import asyncWait from '~/lib/asyncWait'
import { toast } from 'vue-sonner'
import { LucideCheck, LucideCopy } from 'lucide-vue-next'

const isCopy = ref(false)
const props = defineProps<{
    value: string
}>()
const { copy } = useClipboard()
const { t } = useI18n()
</script>

<template>
    <Button
        variant="outline"
        class="transition-all duration-300"
        size="icon"
        @click="
            async () => {
                await copy(props?.value)
                isCopy = true
                toast.success(t('common.copySuccess'))
                await asyncWait(3000)
                isCopy = false
            }
        "
    >
        <component :is="isCopy ? LucideCheck : LucideCopy" class="size-1/2" />
    </Button>
</template>
