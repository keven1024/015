<script setup lang="ts">
import type { GenericObject } from 'vee-validate'
import { useForm } from 'vee-validate'

const props = withDefaults(defineProps<{
    initialValues?: GenericObject
    keepValues?: boolean
}>(), {
    initialValues: () => ({}),
    keepValues: false
})
const form = useForm({
    initialValues: props.initialValues,
    keepValuesOnUnmount: props.keepValues
})
defineExpose({
    form: form
})


</script>
<template>
    <slot v-bind="form" />
</template>
