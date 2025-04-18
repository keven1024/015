<template>
    <Field :name="name" v-slot="{ field }">
        <Tiptap :modelValue="field.value" @update:modelValue="field.onChange" />
    </Field>
</template>

<script setup lang="ts">
import Tiptap from '@/components/Tiptap.vue'

const props = defineProps<{
    name: string
}>()
</script>
