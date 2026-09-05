<script setup lang="ts">
import 'katex/dist/katex.min.css'
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Markdown } from '@tiptap/markdown'
import Placeholder from '@tiptap/extension-placeholder'
import Mathematics from '@tiptap/extension-mathematics'
import { cx } from 'class-variance-authority'
import countWords from '@/lib/countWords'
import CodeBlockShiki from 'tiptap-extension-code-block-shiki'
import Button from '@/components/ui/button/Button.vue'
import { NodeRange } from '@tiptap/extension-node-range'
import TiptapDragHandle from './extensions/DragHandle.vue'
import { TableKit } from '@tiptap/extension-table'
import { createPasteHandler } from './extensions/PasteHandler'
import { ListKit } from '@tiptap/extension-list'
import Image from '@tiptap/extension-image'

const { t } = useI18n()

const props = defineProps<{
    modelValue?: string
    placeholder?: string
    class?: string
}>()
const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
}>()

const editor = shallowRef<Editor>()
const plainText = computed(() => {
    if (!props.modelValue) return ''
    return editor.value?.getText() ?? ''
})

const wordCount = computed(() => countWords(plainText.value))

const clearContent = () => {
    editor.value?.commands.clearContent()
    emit('update:modelValue', '')
}

onMounted(() => {
    editor.value = new Editor({
        content: props.modelValue,
        extensions: [
            TableKit,
            ListKit,
            StarterKit.configure({
                codeBlock: false,
                dropcursor: {
                    color: 'var(--color-blue-300)',
                    width: 2,
                },
            }),
            CodeBlockShiki.configure({
                defaultTheme: 'tokyo-night',
            }),
            Markdown.configure({
                markedOptions: {
                    gfm: true, // GitHub Flavored Markdown
                },
            }),
            Mathematics.configure({
                katexOptions: {
                    throwOnError: false,
                },
            }),
            Placeholder.configure({
                placeholder: props.placeholder ?? '',
            }),
            // CommandsPlugin,
            NodeRange.configure({
                // macOS 按 Cmd、Windows/Linux 按 Ctrl 后跨 block 框选
                key: 'Mod',
            }),
            Image.configure({
                resize: {
                    enabled: true,
                    // directions: ['top', 'bottom', 'left', 'right'], // can be any direction or diagonal combination
                    // minWidth: 50,
                    // minHeight: 50,
                    alwaysPreserveAspectRatio: true,
                },
            }),
            createPasteHandler(t),
        ],
        editorProps: {
            attributes: {
                spellcheck: 'false',
                autocorrect: 'off',
                autocapitalize: 'off',
            },
        },
        onUpdate: () => {
            emit('update:modelValue', (editor.value as any)?.getMarkdown())
        },
    })
})
watch(
    () => props.modelValue,
    (value) => {
        if (value !== (editor.value as any)?.getMarkdown()) {
            editor.value?.commands.setContent(value ?? '', { contentType: 'markdown' })
        }
    }
)
onUnmounted(() => {
    editor.value?.destroy()
})
</script>
<template>
    <div :class="['relative bg-white/50 rounded-md overflow-hidden [&_.tiptap]:h-full', props.class]">
        <editor-content
            :editor="editor as any"
            :class="[
                'prose prose-sm *:outline-none prose-p:my-1 prose-headings:my-2 prose-pre:mb-0 prose-blockquote:border-black/50',
                'selection:bg-primary/20 max-w-full w-full [&_.is-node-active]:bg-blue-100',
                'prose-task:list-none prose-task-item:flex prose-task-item:items-center prose-task-item:gap-2 prose-task-item:before:content-none ',
                'prose-inline-code:rounded-sm prose-inline-code:bg-orange-500/10 prose-inline-code:text-rose-700 prose-inline-code:px-1 prose-inline-code:py-0.5 prose-inline-code:text-xs prose-inline-code:before:content-none prose-inline-code:after:content-none',
            ]"
        >
        </editor-content>
        <TiptapDragHandle v-if="editor" :editor="editor" />
    </div>
    <Button
        type="button"
        variant="ghost"
        size="icon"
        :class="
            cx(
                'absolute right-2 top-2 hover:bg-black/10 transition-all duration-300',
                modelValue?.length && modelValue.length > 0 ? 'opacity-100' : 'opacity-0 pointer-events-none'
            )
        "
        @click="clearContent"
    >
        <LucideX class="size-4" />
    </Button>
    <div
        v-if="plainText.length > 0"
        class="absolute bottom-2 right-3 flex justify-end px-2 py-1 text-xs text-gray-400 select-none bg-white rounded-md"
    >
        {{
            `${wordCount.length} ${t('common.length')} · ${wordCount.characters}
            ${t('common.words')}`
        }}
    </div>
</template>
