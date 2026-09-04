function countWords(text: string): { length: number; characters: number } {
    const trimmed = text?.trim()
    if (!trimmed) return { length: 0, characters: 0 }
    const cjk_emoji = /[\u4e00-\u9fff\u3040-\u30ff\uac00-\ud7af]|[\uD800-\uDBFF][\uDC00-\uDFFF]/g
    const latin = trimmed.replace(cjk_emoji, ' ').match(/\S+/g)?.length ?? 0
    return {
        length: text?.length ?? 0,
        characters: (trimmed.match(cjk_emoji)?.length ?? 0) + latin,
    }
}

export default countWords
