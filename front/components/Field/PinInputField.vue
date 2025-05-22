<script setup lang="ts">
import { PinInput, PinInputGroup, PinInputSlot } from '@/components/ui/pin-input'
import type { RuleExpression } from 'vee-validate'
const props = withDefaults(defineProps<{
    name: string
    label?: string
    placeholder?: string
    class?: string
    length?: number
    rules?: RuleExpression<string>
}>(), {
    placeholder: '○',
    length: 4
})
const { value, setValue } = useField<string>(props?.name, props?.rules)
</script>
<template>
    <div class="flex flex-col gap-2 items-center">
        <Label v-if="props.label">{{ props.label }}</Label>
        <PinInput :model-value="value?.split('')" @update:model-value="(v) => {
            setValue(v.join(''))
        }" :placeholder="props.placeholder">
            <PinInputGroup>
                <PinInputSlot v-for="(id, index) in props.length" :key="id" :index="index" />
            </PinInputGroup>
        </PinInput>
    </div>
</template>
