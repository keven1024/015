<script setup lang="ts">
import type { DateRange, DateValue } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { CalendarDays, ChevronDown } from '@lucide/vue'
import { ref } from 'vue'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

defineOptions({
    inheritAttrs: false,
})

const props = withDefaults(
    defineProps<{
        placeholder?: string
        class?: HTMLAttributes['class']
    }>(),
    {
        placeholder: '-',
    }
)

const model = defineModel<DateRange>('value', {
    default: () => ({
        start: undefined,
        end: undefined,
    }),
})
const open = ref(false)

function formatDate(value?: DateValue) {
    if (!value) return ''

    return [value.year, value.month.toString().padStart(2, '0'), value.day.toString().padStart(2, '0')].join('-')
}

function onRangeChange(value: DateRange) {
    if (value.start && value.end) {
        open.value = false
    }
}
</script>

<template>
    <Popover v-model:open="open">
        <PopoverTrigger as-child>
            <button
                v-bind="$attrs"
                type="button"
                :data-placeholder="!model.start ? '' : undefined"
                :class="
                    cn(
                        'border-input data-[placeholder]:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 flex h-9 min-w-52 items-center justify-between gap-2 rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs outline-none transition-[color,box-shadow] focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50',
                        props.class
                    )
                "
            >
                <span class="flex min-w-0 items-center gap-2">
                    <CalendarDays class="size-4 shrink-0 text-muted-foreground" />
                    <span v-if="model.start" class="truncate flex items-center gap-2">
                        {{ formatDate(model.start) }}
                        <template v-if="model.end"> <LucideMoveRight class="size-3" /> {{ formatDate(model.end) }}</template>
                    </span>
                    <span v-else class="text-muted-foreground">{{ placeholder }}</span>
                </span>
                <ChevronDown class="size-4 shrink-0 text-muted-foreground" />
            </button>
        </PopoverTrigger>

        <PopoverContent class="w-auto p-0">
            <RangeCalendar v-model="model" @update:model-value="onRangeChange" />
        </PopoverContent>
    </Popover>
</template>
