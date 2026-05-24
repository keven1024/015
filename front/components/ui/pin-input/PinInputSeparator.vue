<script setup lang="ts">
import { Minus } from '@lucide/vue'
import { Primitive, type PrimitiveProps, useForwardProps } from 'reka-ui'

const props = defineProps<PrimitiveProps>()
const forwardedProps = useForwardProps(props)
</script>

<template>
    <Primitive data-slot="pin-input-separator" v-bind="forwardedProps">
        <slot>
            <Minus />
        </slot>
    </Primitive>
</template>
