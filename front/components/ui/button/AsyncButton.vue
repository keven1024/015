<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { cva, cx } from 'class-variance-authority'
import type { ButtonProps } from './Button.vue';
import { omit } from 'lodash-es';
const isLoading = ref(false)

interface AsyncButtonProps extends ButtonProps {
    onClick?: (event?: Event) => Promise<void>
    disabled?: boolean
}

const props = withDefaults(defineProps<AsyncButtonProps>(), {
    disabled: false,
    onClick: async () => { }
})

const AsyncButtonIconVariants = cva('size-4 animate-spin absolute', {
    variants: {
        variant: {
            default: 'text-primary-foreground',
            destructive: 'text-white',
            outline: 'text-accent-foreground',
            secondary: 'text-secondary-foreground',
            ghost: 'text-accent-foreground',
            link: 'text-primary'
        }
    },
    defaultVariants: {
      variant: 'default',
    },
})

const handleClick = async (event?: Event) => {
    if (isLoading.value) return
    isLoading.value = true
    try {
        await props.onClick?.(event)
    } finally {
        isLoading.value = false
    }
}
</script>

<template>
    <Button :onClick="handleClick" :disabled="isLoading || disabled"
        :class="cx(isLoading && 'text-transparent', props.class)"
        v-bind="omit(props, 'onClick', 'disabled')">
        <LucideLoader2 :class="AsyncButtonIconVariants({ variant})" v-if="isLoading" />
        <slot />
    </Button>
</template>
