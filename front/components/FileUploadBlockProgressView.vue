<script setup lang="ts">
import type { Graphics } from 'pixi.js'
import type { PixiExpose } from '@/components/Pixi.vue'

const props = withDefaults(
    defineProps<{
        size?: number
        data?: { chunks: Record<number, { status: 'success' | 'error' | 'processing'; createdAt: number }>; chunkLength: number }
    }>(),
    {
        size: 16,
        data: () => ({ chunks: [], chunkLength: 0 }),
    }
)
const pixiRef = ref<PixiExpose>()
const width = computed(() => pixiRef.value?.width || 0)
const height = computed(() => pixiRef.value?.height || 0)
watchEffect(() => {
    console.log('data', props.data)
})
const renderGraphics = (graphics: Graphics) => {
    graphics.clear()
    Object.entries(props?.data?.chunks || {})?.map(([index, item]) => {
        const { status, createdAt } = item || {}
        const size = width.value / props.data?.chunkLength
        const x = Number(index) * size
        graphics.rect(x, 0, size, height.value)
        let color = 0x60a5fa
        if (status === 'success') {
            color = 0x4ade80
        }
        if (status === 'error') {
            color = 0xf87171
        }
        graphics.fill(color)
    })
}
</script>

<template>
    <Pixi ref="pixiRef">
        <graphics @render="renderGraphics" />
    </Pixi>
</template>
