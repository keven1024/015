<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Markdown } from 'tiptap-markdown'
import Placeholder from '@tiptap/extension-placeholder'
const props = defineProps<{
    modelValue?: string
    placeholder?: string
}>()
const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
}>()

const editor = ref<Editor | undefined>(undefined)
onMounted(() => {
    editor.value = new Editor({
        content: props.modelValue,
        extensions: [
            StarterKit,
            Markdown.configure({
                transformPastedText: true,
                transformCopiedText: true,
            }),
            Placeholder.configure({
                placeholder: props.placeholder ?? '',
            }),
            // CommandsPlugin,
        ],
        onUpdate: () => {
            emit('update:modelValue', editor.value?.storage?.markdown?.getMarkdown() ?? '')
        },
    })
})
watch(
    () => props.modelValue,
    (value) => {
        if (value !== editor.value?.storage?.markdown?.getMarkdown()) {
            editor.value?.commands.setContent(value ?? '')
        }
    }
)
onUnmounted(() => {
    editor.value?.destroy()
})
</script>
<template>
    <editor-content
        :editor="editor"
        class="prose prose-sm bg-white/50 rounded-md p-2 [&>*]:outline-none prose-p:my-1 prose-headings:my-2 prose-pre:mb-0"
    />
</template>
