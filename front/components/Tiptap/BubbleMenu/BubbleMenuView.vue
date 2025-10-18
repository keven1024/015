<script lang="ts" setup>
import type { Editor } from '@tiptap/vue-3'
import { BubbleMenu } from '@tiptap/vue-3/menus'
import { LucideAArrowUp, LucideBold, LucideCode, LucideItalic, LucideSparkles, LucideStrikethrough } from 'lucide-vue-next'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import TextInlineAIDrawer from '@/components/Drawer/TextInlineAIDrawer.vue'
import { cx } from 'class-variance-authority'

import { getHTMLFromFragment } from '@tiptap/vue-3'
import showDrawer from '~/lib/showDrawer'

const props = defineProps<{
    editor: Editor
}>()

const show = ref(false)

const menus = [
    {
        type: 'button',
        label: '询问AI',
        icon: LucideSparkles,
        onClick: () => {
            show.value = false
            const { from, to } = props.editor.state.selection || {}
            showDrawer({
                render: ({ ...rest }) =>
                    h(TextInlineAIDrawer, {
                        ...rest,
                        data: { html: getHTMLFromFragment(props.editor.state.doc.slice(from, to).content, props.editor.schema) },
                    }),
            })
        },
    },
    {
        type: 'icon',
        label: '加粗',
        icon: LucideBold,
        onClick: () => {
            props.editor?.chain().focus().toggleBold().run()
        },
    },
    {
        type: 'icon',
        label: '斜体',
        icon: LucideItalic,
        onClick: () => {
            props.editor?.chain().focus().toggleItalic().run()
        },
    },
    {
        type: 'icon',
        label: '删除线',
        icon: LucideStrikethrough,
        onClick: () => {
            props.editor?.chain().focus().toggleStrike().run()
        },
    },
    {
        type: 'icon',
        label: '代码',
        icon: LucideCode,
        onClick: () => {
            props.editor?.chain().focus().toggleCode().run()
        },
    },
]
</script>
<template>
    <bubble-menu
        v-if="editor"
        :editor="editor as any"
        :options="{
            placement: 'bottom',
            offset: 8,
            onShow: () => {
                show = true
            },
            onHide: () => {
                show = false
            },
        }"
    >
        <div :class="cx('bg-white rounded-md overflow-hidden', show ? 'block' : 'hidden')">
            <div class="border border-black/10 bg-primary/30 text-primary p-1 flex flex-row gap-0.5 shadow-md">
                <template v-for="menu in menus" :key="menu.label">
                    <Button v-if="menu.type === 'button'" variant="ghost" size="sm" @click="menu.onClick">
                        <component :is="menu.icon" /> {{ menu.label }}
                    </Button>
                    <TooltipProvider v-if="menu.type === 'icon'">
                        <Tooltip>
                            <TooltipTrigger as-child>
                                <Button variant="ghost" class="!size-8" @click="menu.onClick">
                                    <component :is="menu.icon" />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                {{ menu.label }}
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </template>
            </div>
        </div>
    </bubble-menu>
</template>
