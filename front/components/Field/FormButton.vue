<template>
  <AsyncButton type="button" @click="async (e) => {
    e?.preventDefault()
    await props.onClick?.(form)
  }" :disabled="!isValid">
    <slot />
  </AsyncButton>
</template>
<script setup lang="ts">
import { AsyncButton } from '~/components/ui/button'
import { useFormContext } from 'vee-validate'
const form = useFormContext()

const isValid = ref(false)

watchEffect(async () => {
  const { valid } = await form?.validate()
  isValid.value = valid
})

const props = withDefaults(defineProps<{
  onClick?: (form: ReturnType<typeof useFormContext>) => Promise<void>
}>(), {
  onClick: async () => {}
})
</script>
