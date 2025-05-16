<script lang="ts" setup>
import FileShareResult from '@/components/Result/FileShareResult.vue'

const props = defineProps<{
  data: { file: File, config: any, file_handle_type: string, file_id: string }
}>()

const emit = defineEmits<{
  (e: 'change', key: string): void
}>()

// console.log(props.data)

const handleList = [
  { component: FileShareResult, key: 'file-share' },
  // { component: FileShareResult, key: 'file-share' },
]
const handleComponent = computed(() => {
  return handleList.find((item) => item.key === props?.data?.file_handle_type)?.component
})

</script>
<template>
  <div class="">
    <component :is="handleComponent" :data="data" @change="(key: string) => {
      emit('change', key)
    }" />
  </div>
</template>