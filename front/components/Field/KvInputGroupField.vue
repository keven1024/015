<script setup lang="ts">
import { Label } from '@/components/ui/label'

const props = defineProps<{
    name: string
    label?: string
    keyPlaceholder?: string
    valuePlaceholder?: string
}>()

const { value, setValue } = useField<[string, string][]>(props.name)

const updateKey = (index: number, nextKey: string | number) => {
    const next = [...(value.value ?? [])]
    next[index] = [String(nextKey), next[index]?.[1] ?? '']
    setValue(next)
}

const updateValue = (index: number, nextVal: string | number) => {
    const next = [...(value.value ?? [])]
    next[index] = [next[index]?.[0] ?? '', String(nextVal)]
    setValue(next)
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <div v-for="([key, itemValue], index) in value" :key="index" class="flex flex-row gap-2 items-center">
            <Input :model-value="key" :placeholder="keyPlaceholder" @update:model-value="(v: string | number) => updateKey(index, v)" />
            <Input
                :model-value="String(itemValue ?? '')"
                :placeholder="valuePlaceholder"
                @update:model-value="(v: string | number) => updateValue(index, v)"
            />
            <Button
                type="button"
                variant="ghost"
                size="icon"
                class="bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white"
                @click="() => setValue(value?.filter((_, i) => i !== index))"
            >
                <LucideTrash class="size-4" />
            </Button>
        </div>
        <Button type="button" class="self-start" size="sm" @click="() => setValue([...(value ?? []), ['', '']])">
            <LucidePlus class="size-4" />
            添加
        </Button>
    </div>
</template>
