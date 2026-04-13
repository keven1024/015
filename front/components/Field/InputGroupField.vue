<script setup lang="ts">
import { Label } from '@/components/ui/label'
import type { RuleExpression } from 'vee-validate'
const props = defineProps<{
    name: string
    label?: string
    rules?: RuleExpression<string[]>
}>()
const { value, setValue } = useField<string[]>(props.name, props?.rules)
const addInput = ref('')
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <div v-for="(item, index) in value" :key="`${index}-${item}`" class="flex flex-row gap-2 items-center">
            <Input
                :model-value="item"
                @update:model-value="(v: string | number) => setValue(value.map((o, i) => (i === index ? String(v) : o)))"
                v-bind="$attrs"
            />
            <Button variant="ghost" size="icon" @click="setValue(value.filter((_, i) => i !== index))"><LucideX class="size-4" /></Button>
        </div>
        <div class="flex flex-row gap-2 items-center">
            <Input v-model="addInput" v-bind="$attrs" />
            <Button
                size="icon"
                @click="
                    () => {
                        setValue([...(value || []), addInput])
                        addInput = ''
                    }
                "
                ><LucidePlus class="size-4"
            /></Button>
        </div>
    </div>
</template>
