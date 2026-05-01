<script setup lang="ts">
import { Label } from '@/components/ui/label'
import type { RuleExpression } from 'vee-validate'
const props = defineProps<{
    name: string
    label?: string
    rules?: RuleExpression<string[]>
}>()
const { t } = useI18n()
const { value, setValue, errorMessage } = useField<string[]>(props.name, props?.rules)
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <div v-for="(item, index) in value" class="flex flex-row gap-2 items-center">
            <Input
                :model-value="item"
                @update:model-value="(v: string | number) => setValue(value.map((o, i) => (i === index ? String(v) : o)))"
                :aria-invalid="!!errorMessage || undefined"
                v-bind="$attrs"
            />
            <Button
                variant="ghost"
                size="icon"
                class="bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white"
                @click="setValue(value.filter((_, i) => i !== index))"
            >
                <LucideTrash class="size-4" />
            </Button>
        </div>
        <Button class="self-start" size="sm" @click="() => setValue([...(value || []), ''])">
            <LucidePlus class="size-4" />
            {{ t('common.add') }}
        </Button>
    </div>
</template>
