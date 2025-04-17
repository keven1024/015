<template>
    <Field :name="props.name" v-slot="{ field }">
        <div class="flex flex-col gap-2">
            <Label v-if="props.label">{{ props.label }}</Label>
            <div class="border rounded-md">
                <Tiptap 
                    :modelValue="field.value" 
                    @update:modelValue="field.onChange" 
                />
            </div>
        </div>
    </Field>
</template>

<script setup lang="ts">
import Tiptap from '@/components/Tiptap.vue'

const props = defineProps<{
    name: string
    label?: string
    placeholder?: string
}>()
</script>
