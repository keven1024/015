import calcFileHash from './calcFileHash'

// 监听主线程消息
self.onmessage = async (e: MessageEvent<{ file: File }>) => {
    const { file } = e.data || {}
    const hash = await calcFileHash({ file })
    self.postMessage({ hash })
}
