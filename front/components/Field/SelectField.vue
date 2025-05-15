<script setup lang="ts">
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import type { RuleExpression } from 'vee-validate'
type SelectValue = string | number
const props = defineProps<{
    name: string
    placeholder?: string
    label?: string
    rules?: RuleExpression<SelectValue>
    options?: {
        label?: string
        value: SelectValue
    }[]
}>()
const { value } = useField<SelectValue>(props.name, props?.rules)
</script>

<template>
    <Select v-model="value">
        <SelectTrigger>
            <SelectValue :placeholder="placeholder" />
        </SelectTrigger>
        <SelectContent>
            <SelectGroup>
                <SelectLabel class="text-xs" v-if="label">{{ label }}</SelectLabel>
                <SelectItem v-for="item in options" :key="item.value" :value="item.value">
                    {{ item?.label ?? item.value }}
                </SelectItem>
            </SelectGroup>
        </SelectContent>
    </Select>
</template>