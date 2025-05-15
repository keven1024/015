<template>
    <Tiptap v-model="value" />
</template>

<script setup lang="ts">
import Tiptap from '@/components/Tiptap.vue'
import type { RuleExpression } from 'vee-validate'
const props = defineProps<{
    name: string
    rules?: RuleExpression<string>
}>()
const { value, } = useField<string>(props.name, props.rules)
</script>
