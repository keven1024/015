<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { AutocompleteAnchor, AutocompleteContent, AutocompleteInput, AutocompleteItem, AutocompleteRoot, AutocompleteViewport } from 'reka-ui'
import type { Component } from 'vue'
import InputField from '../Field/InputField.vue'

type KvInputValueComponentConfig = [(key: string) => boolean, Component]

type KvInputConfig = {
    key?: {
        placeholder?: string
        enum?: string[]
    }
    value?: {
        placeholder?: string
        component?: KvInputValueComponentConfig[]
        default?: Component
    }
}

const defaultConfig = {
    key: {},
    value: {
        default: InputField,
    },
} satisfies Required<KvInputConfig>

const props = defineProps<{
    name: string
    label?: string
    config?: KvInputConfig
}>()

const { t } = useI18n()
const config = computed(() => {
    return {
        key: { ...defaultConfig.key, ...(props.config?.key ?? {}) },
        value: { ...defaultConfig.value, ...(props.config?.value ?? {}) },
    }
})

const { value, setValue } = useField<[string, string][]>(props.name)

const updateKey = (index: number, nextKey: string | number) => {
    const next = [...(value.value ?? [])]
    next[index] = [String(nextKey), next[index]?.[1] ?? '']
    setValue(next)
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <Label v-if="label">{{ label }}</Label>
        <div v-for="([key, _], index) in value" class="flex flex-row gap-2 items-center">
            <AutocompleteRoot class="basis-40 relative" :model-value="String(key ?? '')" @update:model-value="(v) => updateKey(index, v)">
                <AutocompleteAnchor>
                    <AutocompleteInput
                        class="w-full placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 border-input flex h-9 min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
                        :placeholder="config.key.placeholder"
                    >
                    </AutocompleteInput>
                </AutocompleteAnchor>
                <AutocompleteContent
                    v-if="config.key?.enum"
                    class="bg-popover border rounded-md shadow-md z-50 w-(--reka-autocomplete-trigger-width) absolute inset-x-0"
                >
                    <AutocompleteViewport class="p-1">
                        <AutocompleteItem
                            v-for="opt in config.key?.enum"
                            :key="opt"
                            :value="opt"
                            class="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none data-highlighted:bg-accent data-highlighted:text-accent-foreground"
                        >
                            {{ opt }}
                        </AutocompleteItem>
                    </AutocompleteViewport>
                </AutocompleteContent>
            </AutocompleteRoot>
            <div class="flex-1">
                <component
                    :is="config.value.component?.find(([isMatchCom]) => isMatchCom(key))?.[1] ?? config.value.default"
                    :name="`${props.name}[${index}][1]`"
                    :placeholder="config.value.placeholder"
                    class="w-full"
                />
            </div>
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
            {{ t('common.add') }}
        </Button>
    </div>
</template>
