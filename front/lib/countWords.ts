function countWords(text: string): number {
    const trimmed = text?.trim()
    if (!trimmed) return 0
    const cjk = trimmed.match(/[\u4e00-\u9fff\u3040-\u30ff\uac00-\ud7af]/g)?.length ?? 0
    const latin = trimmed.replace(/[\u4e00-\u9fff\u3040-\u30ff\uac00-\ud7af]/g, ' ').match(/\S+/g)?.length ?? 0
    return cjk + latin
}

export default countWords
