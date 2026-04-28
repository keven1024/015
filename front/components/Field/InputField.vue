<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import type { RuleExpression } from 'vee-validate'

const props = defineProps<{
    name: string
    label?: string
    rules?: RuleExpression<string>
}>()

const { value, errorMessage } = useField<string>(props.name, props.rules)
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="props.label">{{ props.label }}</Label>
        <Input v-model="value" :aria-invalid="!!errorMessage || undefined" v-bind="$attrs" />
    </div>
</template>
