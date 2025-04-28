<script setup lang="ts">
import { useDropZone } from '@vueuse/core'
const dropZoneRef = ref()

const props = defineProps<{
  accept?: string[]
}>()

const accept = computed(() => (props?.accept || ['*'])?.join(','))

const emit = defineEmits<{
  (e: 'onChange', file: File): void
}>()

const { isOverDropZone } = useDropZone(dropZoneRef, {
  onDrop: (file) => {
    if (file?.[0]) {
      emit('onChange', file?.[0])
    }
  },
  // 指定要接收的数据类型
  dataTypes: (types) => {
    for (const type of types) {
      for (const acceptType of accept.value.split(',')) {
        if (acceptType === '*') {
          return true
        }
        if (acceptType?.endsWith('*')) {
          const [acceptTypePrefix,] = acceptType?.split('/')
          if (!acceptTypePrefix) {
            return true
          }
          if (type?.startsWith(acceptTypePrefix)) {
            return true
          }
        }
        if (acceptType === type) {
          return true
        }
      }
    }
    return false
  },
  // 控制多文件拖放
  multiple: false,
  // 是否阻止未处理事件的默认行为
  preventDefaultForUnhandled: false,
})

const { open, onChange } = useFileDialog({
  accept: accept.value, // Set to accept only image files
  directory: false,
})
onChange((files) => {
  if (files?.[0]) {
    emit('onChange', files?.[0])
  }
})
</script>

<template>
  <div ref="dropZoneRef" @click="open">
    <slot :isOverDropZone="isOverDropZone" />
  </div>
</template>