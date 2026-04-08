<script setup lang="ts">
import type { filePreview } from './Index.vue'

const props = defineProps<{
    file: File | filePreview
}>()

const { state: thumbnailUrl } = useAsyncState(async () => {
    if (props.file instanceof File && props.file.type.startsWith('video/')) {
        return await extractThumbnail(props.file)
    }
    return null
}, null)

async function extractThumbnail(file: File): Promise<string> {
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

onUnmounted(() => {
    if (!!thumbnailUrl.value) {
        URL.revokeObjectURL(thumbnailUrl.value)
    }
})
</script>

<template>
    <div v-if="thumbnailUrl" class="relative grayscale-50 overflow-hidden">
        <img :src="thumbnailUrl" class="object-contain block max-w-full max-h-full" />
        <LucidePlay class="size-[40%] absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white" />
    </div>
    <div v-else>
        <div class="h-16 aspect-video flex justify-center items-center bg-white/80 rounded-sm">
            <LucideVideo class="size-10" />
        </div>
    </div>
</template>
