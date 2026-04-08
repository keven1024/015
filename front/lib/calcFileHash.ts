import { noop } from 'lodash-es'
import { createSHA1 } from 'hash-wasm'

interface CalcFileHashProps {
    file: File
    onProgress?: (current: number) => void
    chunkSize?: number
    engine?: 'native' | 'wasm'
}

const calcFileHash = async (props: CalcFileHashProps) => {
    const { file, onProgress = noop, chunkSize = 100, engine = 'native' } = props || {}

    if (engine === 'native') {
        const buffer = await file.arrayBuffer()
        return calcNativeHash(buffer)
    }

    const chunkBytes = chunkSize * 1024 * 1024
    const hasher = await createSHA1()
    let offset = 0
    while (offset < file.size) {
        const buffer = await file.slice(offset, offset + chunkBytes).arrayBuffer()
        hasher.update(new Uint8Array(buffer))
        offset += chunkBytes
        onProgress(Math.min(offset, file.size) / file.size)
    }
    return hasher.digest('hex')
}

export const calcNativeHash = async (buffer: BufferSource) => {
    const hashBuffer = await crypto.subtle.digest('SHA-1', buffer)
    return Array.from(new Uint8Array(hashBuffer))
        .map((b) => b.toString(16).padStart(2, '0'))
        .join('')
}

export default calcFileHash
