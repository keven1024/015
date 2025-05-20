import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from 'prosemirror-state'
import { marked } from 'marked'

export const MarkdownPaste = Extension.create({
    name: 'markdownPaste',
    addProseMirrorPlugins() {
        return [
            new Plugin({
                key: new PluginKey('markdownPaste'),
                props: {
                    handlePaste: (view, event) => {
                        const clipboardText = event.clipboardData?.getData('text/plain')
                        if (clipboardText) {
                            const html = marked(clipboardText)
                            this.editor.commands.insertContent(html)
                            return true
                        }
                        return false
                    },
                },
            }),
        ]
    },
})