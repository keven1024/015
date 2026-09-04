import { h } from 'vue'
import { Extension } from '@tiptap/vue-3'
import { Plugin } from '@tiptap/pm/state'
import ConfirmActionDrawer from '~/components/Drawer/ConfirmActionDrawer.vue'
import showDrawer from '~/lib/showDrawer'
import getMDRenderHtml from '~/lib/getMDRenderHtml'
import isStringMarkdown from '~/lib/isStringMarkdown'

export const PasteHandler = Extension.create({
    name: 'pasteHandler',
    addProseMirrorPlugins() {
        return [
            new Plugin({
                props: {
                    handlePaste: (_view, event) => {
                        const text = event.clipboardData?.getData('text/plain')
                        if (!text) return false
                        if (!isStringMarkdown(text)) return false

                        const { from, to } = this.editor.state.selection
                        void (async () => {
                            const shouldParse = await showDrawer<boolean>({
                                render: ({ hide }) =>
                                    h(ConfirmActionDrawer, {
                                        title: '检测到 Markdown',
                                        desc: '是否解析并粘贴 Markdown 内容？',
                                        btnLabel: '是',
                                        btnClass: 'bg-primary text-primary-foreground hover:bg-primary/90',
                                        onClick: () => hide(true),
                                    }),
                            })
                            if (!shouldParse) return

                            this.editor
                                .chain()
                                .focus()
                                .setTextSelection({ from, to })
                                .insertContent(await getMDRenderHtml(text))
                                .run()
                        })()

                        return true
                    },
                },
            }),
        ]
    },
})
