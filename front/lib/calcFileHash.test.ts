import { describe, it, expect } from 'vitest'
import calcFileHash from './calcFileHash'

const makeFile = (content: string) => new File([content], 'test.txt', { type: 'text/plain' })

describe('calcFileHash 引擎一致性', () => {
    it('对于小体积内容，native 和 wasm 生成相同哈希', async () => {
        const file = makeFile('hello world')
        const [native, wasm] = await Promise.all([calcFileHash({ file, engine: 'native' }), calcFileHash({ file, engine: 'wasm' })])
        expect(native).toBe(wasm)
    })

    it('对于二进制内容，native 和 wasm 生成相同哈希', async () => {
        const bytes = new Uint8Array(1024).map((_, i) => i % 256)
        const file = new File([bytes], 'bin.bin', { type: 'application/octet-stream' })
        const [native, wasm] = await Promise.all([calcFileHash({ file, engine: 'native' }), calcFileHash({ file, engine: 'wasm' })])
        expect(native).toBe(wasm)
    })

    it('哈希值应为 40 位十六进制字符串', async () => {
        const file = makeFile('abc')
        const hash = await calcFileHash({ file, engine: 'native' })
        expect(hash).toMatch(/^[0-9a-f]{40}$/)
    })

    it('不同内容应生成不同哈希', async () => {
        const [h1, h2] = await Promise.all([
            calcFileHash({ file: makeFile('foo'), engine: 'native' }),
            calcFileHash({ file: makeFile('bar'), engine: 'native' }),
        ])
        expect(h1).not.toBe(h2)
    })
})
