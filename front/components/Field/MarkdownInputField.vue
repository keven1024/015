<template>
    <Tiptap v-model="value" />
</template>

<script setup lang="ts">
import Tiptap from '@/components/Tiptap.vue'

const props = defineProps<{
    name: string
    rules?: string
}>()
const { value, } = useField<string>(props.name, props.rules)
</script>
