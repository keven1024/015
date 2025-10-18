<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Markdown } from 'tiptap-markdown'
import Placeholder from '@tiptap/extension-placeholder'
import BubbleMenuView from './BubbleMenu/BubbleMenuView.vue'
import { cx } from 'class-variance-authority'

const props = defineProps<{
    modelValue?: string
    placeholder?: string
    class?: string
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
            emit('update:modelValue', (editor.value as any)?.storage?.markdown?.getMarkdown() ?? '')
        },
    })
})
watch(
    () => props.modelValue,
    (value) => {
        if (value !== (editor.value as any)?.storage?.markdown?.getMarkdown()) {
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
        :editor="editor as any"
        :class="
            cx(
                'prose prose-sm bg-white/50 rounded-md p-2 [&>*]:outline-none prose-p:my-1 prose-headings:my-2 prose-pre:mb-0 prose-blockquote:border-black/50 selection:bg-primary/20 max-w-full',
                props.class
            )
        "
    >
    </editor-content>
    <BubbleMenuView :editor="editor as any" />
</template>
