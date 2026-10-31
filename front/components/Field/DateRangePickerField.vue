<script setup lang="ts">
import type { DateRange } from '@/components/DateRangePicker.vue'
import type { RuleExpression } from 'vee-validate'
import { Label } from '@/components/ui/label'
import DateRangePicker from '@/components/DateRangePicker.vue'

defineOptions({
    inheritAttrs: false,
})

const props = defineProps<{
    name: string
    label?: string
    rules?: RuleExpression<DateRange>
}>()

const { value, errorMessage } = useField<DateRange>(props.name, props.rules)
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="props.label">{{ props.label }}</Label>
        <DateRangePicker v-model="value" :aria-invalid="!!errorMessage || undefined" v-bind="$attrs" />
    </div>
</template>
