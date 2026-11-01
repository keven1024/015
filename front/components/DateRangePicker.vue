<script setup lang="ts">
import { parseDate } from '@internationalized/date'
import type { DateRange as RekaDateRange, DateValue } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { CalendarDays, ChevronDown } from '@lucide/vue'
import dayjs from 'dayjs'
import { computed, ref } from 'vue'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

defineOptions({
    inheritAttrs: false,
})

export type DateRange = {
    start?: Date
    end?: Date
}

export type DateRangeOption = {
    label?: string
    value: DateRange
}

const props = withDefaults(
    defineProps<{
        placeholder?: string
        class?: HTMLAttributes['class']
        options?: DateRangeOption[]
    }>(),
    {
        placeholder: '-',
    }
)

const model = defineModel<DateRange>({
    default: () => ({
        start: undefined,
        end: undefined,
    }),
})
const open = ref(false)

const calendarRange = computed<RekaDateRange>(() => ({
    start: model.value.start ? parseDate(dayjs(model.value.start).format('YYYY-MM-DD')) : undefined,
    end: model.value.end ? parseDate(dayjs(model.value.end).format('YYYY-MM-DD')) : undefined,
}))

function formatDate(value?: Date) {
    if (!value) return ''

    return dayjs(value).format('YYYY-MM-DD')
}

function onRangeChange(value: RekaDateRange) {
    const nextValue: DateRange = {
        start: value?.start ? dayjs(value.start.toString()).toDate() : undefined,
        end: value?.end ? dayjs(value.end.toString()).toDate() : undefined,
    }

    model.value = nextValue

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
            <RangeCalendar :model-value="calendarRange" @update:model-value="onRangeChange" />
        </PopoverContent>
    </Popover>
</template>
