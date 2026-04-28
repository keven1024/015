<script setup lang="ts">
import { Label } from '@/components/ui/label'

const props = defineProps<{
    name: string
    label?: string
    keyPlaceholder?: string
    valuePlaceholder?: string
}>()

const { value, setValue } = useField<Record<string, any>>(props.name)
const addKey = ref('')
const addValue = ref('')

const currentValue = computed(() => (value.value && typeof value.value === 'object' && !Array.isArray(value.value) ? value.value : {}))

const entries = computed(() => Object.entries(currentValue.value))

const removeItem = (key: string) => {
    const nextValue = { ...currentValue.value }
    delete nextValue[key]
    setValue(nextValue)
}

const updateKey = (oldKey: string, nextKeyValue: string | number) => {
    const nextKey = String(nextKeyValue).trim()
    if (!nextKey || nextKey === oldKey) {
        return
    }

    const nextValue = { ...currentValue.value }
    const oldValue = nextValue[oldKey]
    delete nextValue[oldKey]
    nextValue[nextKey] = oldValue
    setValue(nextValue)
}

const updateValue = (key: string, nextItemValue: string | number) => {
    setValue({ ...currentValue.value, [key]: String(nextItemValue) })
}

const addItem = () => {
    const nextKey = addKey.value.trim()
    const nextValue = addValue.value.trim()

    if (!nextKey || !nextValue) {
        return
    }

    setValue({ ...currentValue.value, [nextKey]: nextValue })
    addKey.value = ''
    addValue.value = ''
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <div v-for="[key, itemValue] in entries" :key="key" class="flex flex-row gap-2 items-center">
            <Input :model-value="key" :placeholder="keyPlaceholder" @update:model-value="(nextValue: string | number) => updateKey(key, nextValue)" />
            <Input
                :model-value="String(itemValue ?? '')"
                :placeholder="valuePlaceholder"
                @update:model-value="(nextValue: string | number) => updateValue(key, nextValue)"
            />
            <Button
                type="button"
                variant="ghost"
                size="icon"
                class="bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white"
                @click="removeItem(key)"
            >
                <LucideTrash class="size-4" />
            </Button>
        </div>
        <div class="flex flex-row gap-2 items-center">
            <Input v-model="addKey" :placeholder="keyPlaceholder" />
            <Input v-model="addValue" :placeholder="valuePlaceholder" />
            <Button type="button" size="icon" @click="addItem">
                <LucidePlus class="size-4" />
            </Button>
        </div>
    </div>
</template>
