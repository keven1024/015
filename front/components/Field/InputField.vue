<template>
    <Field :name="props.name" v-slot="{ field }">
        <div class="flex flex-col gap-2">
            <Label v-if="props.label">{{ props.label }}</Label>
            <Input v-bind="{ ...field, ...$attrs }" />
        </div>
    </Field>
</template>
<script setup lang="ts">
import { Input } from '@/components/ui/input'
const props = defineProps<{
    name: string
    label?: string
}>()
</script>
