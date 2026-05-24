<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

const props = withDefaults(
    defineProps<{
        name: string
        label?: string
        placeholder?: string
        rows?: number
    }>(),
    {
        rows: 3,
    }
)

const { value } = useField<string>(props.name)
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <Textarea v-model="value" :placeholder="placeholder" :rows="rows" v-bind="$attrs" />
    </div>
</template>
