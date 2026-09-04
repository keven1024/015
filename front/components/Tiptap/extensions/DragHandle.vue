<script setup lang="ts">
import { DragHandle } from '@tiptap/extension-drag-handle-vue-3'
import type { Editor } from '@tiptap/vue-3'
import {
    LucideArrowDownToLine,
    LucideCodeXml,
    LucideHeading1,
    LucideHeading2,
    LucideHeading3,
    LucideHeading4,
    LucideList,
    LucideListOrdered,
    LucideQuote,
    LucideRefreshCcw,
    LucideSigma,
    LucideText,
    LucideTrash2,
} from '@lucide/vue'
import type { Component } from 'vue'
import Button from '@/components/ui/button/Button.vue'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { isNumber } from 'lodash-es'

const { editor } = defineProps<{
    editor: Editor
}>()

const activeNode = ref<{ type: string; from: number; to: number }>()

type BlockMenuItem = {
    label: string
    type: 'paragraph' | 'heading' | 'bulletList' | 'orderedList' | 'codeBlock' | 'blockquote' | 'blockMath'
    level?: 1 | 2 | 3 | 4
    icon: Component
}

const handleBlock = ({ type, level }: BlockMenuItem, insert = false) => {
    const chain = editor.chain().focus()
    const { from, to, type: currentType } = activeNode.value || {}

    if (!isNumber(to) || !isNumber(from) || !currentType) {
        return
    }
    if (!insert && activeNode.value) {
        chain.setTextSelection(from + 1)
    }

    if (type === 'paragraph') {
        if (insert) {
            chain.setTextSelection(to).insertContent({ type: 'paragraph' }).run()
        } else {
            chain.lift(type).setParagraph().run()
        }
        return
    }

    if (type === 'heading') {
        if (insert) {
            chain
                .setTextSelection(to)
                .insertContent({ type: 'heading', attrs: { level: level ?? 1 } })
                .run()
        } else {
            chain.setHeading({ level: level ?? 1 }).run()
        }
        return
    }

    if (type === 'bulletList') {
        if (insert) {
            chain
                .setTextSelection(to)
                .insertContent({ type: 'bulletList', content: [{ type: 'listItem', content: [{ type: 'paragraph' }] }] })
                .run()
        } else {
            chain.toggleBulletList().run()
        }
        return
    }

    if (type === 'orderedList') {
        if (insert) {
            chain
                .setTextSelection(to)
                .insertContent({ type: 'orderedList', content: [{ type: 'listItem', content: [{ type: 'paragraph' }] }] })
                .run()
        } else {
            chain.toggleOrderedList().run()
        }
        return
    }

    if (type === 'codeBlock') {
        if (insert) {
            chain.setTextSelection(to).insertContent({ type: 'codeBlock' }).run()
        } else {
            chain.setCodeBlock().run()
        }
        return
    }

    if (insert) {
        chain
            .setTextSelection(to)
            .insertContent({ type: 'blockquote', content: [{ type: 'paragraph' }] })
            .run()
    } else {
        chain.toggleBlockquote().run()
    }
}

type AddMenuItem =
    | { type: 'sub'; label: string; icon: Component; children: (BlockMenuItem & { handle: () => void })[] }
    | { type: 'separator' }
    | { type: 'item'; label: string; icon: Component; class?: string; handle: () => void }

const blockMenu: BlockMenuItem[] = [
    { label: '文本', type: 'paragraph', icon: LucideText },
    { label: '标题 1', type: 'heading', level: 1, icon: LucideHeading1 },
    { label: '标题 2', type: 'heading', level: 2, icon: LucideHeading2 },
    { label: '标题 3', type: 'heading', level: 3, icon: LucideHeading3 },
    { label: '标题 4', type: 'heading', level: 4, icon: LucideHeading4 },
    { label: '项目符号列表', type: 'bulletList', icon: LucideList },
    { label: '有序列表', type: 'orderedList', icon: LucideListOrdered },
    { label: '代码', type: 'codeBlock', icon: LucideCodeXml },
    { label: '引用', type: 'blockquote', icon: LucideQuote },
]
const activeNodeMenu = computed(() => blockMenu.find((item) => item.type === activeNode.value?.type))
const addMenu: AddMenuItem[] = [
    { type: 'sub', label: '转换成', icon: LucideRefreshCcw, children: blockMenu.map((item) => ({ ...item, handle: () => handleBlock(item) })) },
    {
        type: 'sub',
        label: '向下插入',
        icon: LucideArrowDownToLine,
        children: blockMenu.map((item) => ({ ...item, handle: () => handleBlock(item, true) })),
    },
    { type: 'separator' },
    {
        type: 'item',
        label: '删除',
        icon: LucideTrash2,
        class: 'text-destructive focus:text-destructive',
        handle: () => {
            if (activeNode.value) {
                const { from, to } = activeNode.value
                editor.chain().focus().deleteRange({ from, to }).run()
                return
            }

            editor.chain().focus().deleteSelection().run()
        },
    },
]
const setActiveNode = ({ node, pos }: { node: { nodeSize: number; type: { name: string } } | null; pos: number }) => {
    activeNode.value = node ? { type: node.type.name, from: pos, to: pos + node.nodeSize } : undefined
}
</script>

<template>
    <DragHandle
        :editor="editor"
        :on-node-change="setActiveNode"
        :compute-position-config="{
            placement: 'left-start',
            strategy: 'fixed',
        }"
    >
        <div class="flex h-full items-center gap-0.5 pr-1 pt-1">
            <DropdownMenu>
                <DropdownMenuTrigger as-child>
                    <Button type="button" variant="ghost" size="icon" class="size-5" @mousedown.stop>
                        <LucidePlus class="size-3.5" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent side="right" align="start">
                    <DropdownMenuLabel class="text-xs px-1 py-0.5 opacity-70">{{ activeNodeMenu?.label ?? activeNode?.type }}</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <template v-for="(menu, index) in addMenu" :key="index">
                        <DropdownMenuSub v-if="menu.type === 'sub'">
                            <DropdownMenuSubTrigger class="gap-2">
                                <component :is="menu.icon" class="size-4" />
                                {{ menu.label }}
                            </DropdownMenuSubTrigger>
                            <DropdownMenuSubContent>
                                <DropdownMenuItem v-for="item in menu.children" :key="item.label" @select="item.handle">
                                    <component :is="item.icon" class="size-4" />
                                    {{ item.label }}
                                </DropdownMenuItem>
                            </DropdownMenuSubContent>
                        </DropdownMenuSub>
                        <DropdownMenuSeparator v-else-if="menu.type === 'separator'" />
                        <DropdownMenuItem v-else :class="menu.class" @select="menu.handle">
                            <component :is="menu.icon" class="size-4" />
                            {{ menu.label }}
                        </DropdownMenuItem>
                    </template>
                </DropdownMenuContent>
            </DropdownMenu>
            <Button type="button" variant="ghost" size="icon" class="size-5 cursor-grab active:cursor-grabbing">
                <LucideGripVertical class="size-3.5" />
            </Button>
        </div>
    </DragHandle>
</template>
