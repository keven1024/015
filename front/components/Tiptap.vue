<template>
    <editor-content :editor="editor" />
</template>

<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const props = defineProps<{
    modelValue: string
}>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const editor = ref<Editor | undefined>(undefined)
onMounted(() => {
  editor.value = new Editor({
    content: props.modelValue,
    extensions: [StarterKit],
    onUpdate: () => {
        // HTML
        emit('update:modelValue', editor.value?.getHTML() ?? '')
      }
  })
})
onUnmounted(() => { 
    editor.value?.destroy()
})
</script>