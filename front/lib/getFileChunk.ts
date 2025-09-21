const getFileChunk = (file: File, start: number, chunk_size: number): Promise<ArrayBuffer> => {
    const fileReader = new FileReader()
    return new Promise((resolve, reject) => {
        const chunk = file.slice(start, start + chunk_size)
        fileReader.onload = (e) => resolve(e.target?.result as ArrayBuffer)
        fileReader.onerror = reject
        fileReader.readAsArrayBuffer(chunk)
    })
}

export default getFileChunk
