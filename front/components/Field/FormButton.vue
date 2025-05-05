<template>
  <Button type="button" @click="(e) => {
    e.preventDefault()
    emit('click', form)
  }" :disabled="!isValid">
    <slot />
  </Button>
</template>
<script setup lang="ts">
import { Button } from '~/components/ui/button'
import { useFormContext } from 'vee-validate'
const form = useFormContext()

const isValid = ref(false)

watchEffect(async () => {
  const { valid } = await form?.validate()
  isValid.value = valid
})

const emit = defineEmits<{
  (e: 'click', form: ReturnType<typeof useFormContext>): void
}>()
</script>
