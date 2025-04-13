<template>
    <editor-content :editor="editor" class="prose prose-sm bg-white/50 rounded-md p-2 [&>*]:outline-none" />
</template>

<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Markdown } from 'tiptap-markdown';

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
    extensions: [StarterKit, Markdown],
    onUpdate: () => {
        // HTML
        emit('update:modelValue', editor.value?.storage?.markdown?.getMarkdown() ?? '')
      }
  })
})
onUnmounted(() => { 
    editor.value?.destroy()
})
</script>