<script setup lang="ts">
import { Button } from '@/components/ui/button'
const isLoading = ref(false)

const props = withDefaults(defineProps<{
    disabled?: boolean
    onClick?: (event?: Event) => Promise<void>
}>(), {
    disabled: false,
    onClick: async () => {}
})

const handleClick = async (event?: Event) => {
    if (isLoading.value) return
    
    isLoading.value = true
    try {
        await props.onClick?.(event)
    }finally {
        isLoading.value = false
    }
}
</script>

<template>
    <Button :onClick="handleClick" :disabled="isLoading || disabled">
        <slot />
    </Button>
</template>
