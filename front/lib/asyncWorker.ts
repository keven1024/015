const asyncWorker = (w: new () => Worker, opts: { data: any }) => {
    const { data } = opts || {}
    return new Promise<MessageEvent>((resolve, reject) => {
        const worker = new w()
        worker.postMessage(data || {})
        worker.onmessage = (e: MessageEvent) => {
            resolve(e)
        }
        worker.onerror = (e: ErrorEvent) => {
            reject(e)
        }
    })
}
export default asyncWorker
