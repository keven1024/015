<script setup lang="ts">
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import type { RuleExpression } from 'vee-validate'
const props = defineProps<{
    name: string
    label?: string
    rules?: RuleExpression<boolean>
}>()
const { value } = useField<boolean>(props.name, props?.rules)
</script>

<template>
  <div class="flex flex-row gap-2 items-center">
    <Label v-if="label">{{ label }}</Label>
    <Switch v-model="value" v-bind="$attrs" />
  </div>
</template>