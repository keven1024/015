import calcFileHash from './calcFileHash'

// 监听主线程消息
self.onmessage = async (e: MessageEvent<{ file: File; engine?: 'native' | 'wasm' }>) => {
    const { file, engine } = e.data || {}
    const hash = await calcFileHash({ file, engine })
    self.postMessage({ hash })
}
