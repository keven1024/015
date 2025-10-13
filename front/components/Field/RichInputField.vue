<template>
    <Tiptap :model-value="jsonValue" @update:model-value="(v) => setValue(JSON.stringify(v))" />
</template>

<script setup lang="ts">
import Tiptap from '@/components/Tiptap.vue'
import type { RuleExpression } from 'vee-validate'
const props = defineProps<{
    name: string
    rules?: RuleExpression<string>
}>()
const { value, setValue } = useField<string>(props.name, props.rules)
const jsonValue = computed(() => {
    try {
        return value.value ? JSON.parse(value.value) : {}
    } catch (error) {
        return {}
    }
})
</script>
