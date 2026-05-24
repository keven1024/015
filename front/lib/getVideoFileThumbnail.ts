async function getVideoFileThumbnail(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const video = document.createElement('video')
        const objectUrl = URL.createObjectURL(file)

        video.muted = true
        video.playsInline = true
        video.preload = 'metadata'

        video.onloadedmetadata = () => {
            video.currentTime = video.duration * 0.1
        }

        video.onseeked = async () => {
            try {
                // WebCodecs: capture a VideoFrame from the video element
                const frame = new VideoFrame(video)
                const bitmap = await createImageBitmap(frame)
                frame.close()

                const canvas = new OffscreenCanvas(bitmap.width, bitmap.height)
                const ctx = canvas.getContext('2d')!
                ctx.drawImage(bitmap, 0, 0)
                bitmap.close()

                const blob = await canvas.convertToBlob({ type: 'image/jpeg', quality: 0.8 })
                URL.revokeObjectURL(objectUrl)
                resolve(URL.createObjectURL(blob))
            } catch (e) {
                URL.revokeObjectURL(objectUrl)
                reject(e)
            }
        }

        video.onerror = () => {
            URL.revokeObjectURL(objectUrl)
            reject(new Error('Video load failed'))
        }

        video.src = objectUrl
    })
}

export default getVideoFileThumbnail
