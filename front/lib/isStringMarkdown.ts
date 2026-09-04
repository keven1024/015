import MarkdownItAsync from 'markdown-it-async'

const md = MarkdownItAsync()

const blockTokens = new Set(['heading_open', 'bullet_list_open', 'ordered_list_open', 'blockquote_open', 'fence', 'code_block', 'hr'])
const inlineTokens = new Set(['strong_open', 'em_open', 'code_inline', 'link_open', 'image'])

const isStringMarkdown = (text: string) => {
    const tokens = md.parse(text, {})

    return (
        tokens.some(({ type }) => blockTokens.has(type)) ||
        tokens.some(({ type, children }) => type === 'inline' && children?.some(({ type }) => inlineTokens.has(type)))
    )
}

export default isStringMarkdown
