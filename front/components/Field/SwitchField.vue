<script setup lang="ts">
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
const props = defineProps<{
    name: string
    label?: string
    rules?: string
}>()
const { value } = useField<boolean>(props.name, props?.rules)
</script>

<template>
  <div class="flex flex-row gap-2 items-center">
    <Label v-if="label">{{ label }}</Label>
    <Switch v-model="value" />
  </div>
</template>