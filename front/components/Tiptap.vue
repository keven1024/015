<script setup lang="ts">
import { Editor, EditorContent, type JSONContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
const props = defineProps<{
    modelValue?: JSONContent
    placeholder?: string
}>()
const emit = defineEmits<{
    (e: 'update:modelValue', value?: JSONContent): void
}>()

const editor = ref<Editor | undefined>(undefined)
onMounted(() => {
    editor.value = new Editor({
        content: props.modelValue,
        extensions: [
            StarterKit,
            Placeholder.configure({
                placeholder: props.placeholder ?? '',
            }),
            // CommandsPlugin,
        ],
        onUpdate: (v) => {
            emit('update:modelValue', editor.value?.getJSON?.())
        },
    })
})
watch(
    () => props.modelValue,
    (value) => {
        if (value !== editor.value?.getJSON?.()) {
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
