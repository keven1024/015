<script setup lang="ts">
import type { Graphics } from 'pixi.js'
import type { PixiExpose } from '@/components/Pixi.vue'

let props = withDefaults(
    defineProps<{
        size?: number
        gap?: number
        data?: { chunks: Record<number, { status: 'success' | 'error' | 'processing'; createdAt: number }>; chunkLength: number }
    }>(),
    {
        size: 16,
        gap: 4,
        data: () => ({ chunks: [], chunkLength: 0 }),
    }
)
const pixiRef = ref<PixiExpose>()
const width = computed(() => pixiRef.value?.width)
const height = computed(() => pixiRef.value?.height)

const size = ref(props.size)
const gap = ref(props.gap)
const itemsPerRow = ref(0)
const itemsPerCol = computed(() => Math.ceil(props?.data?.chunkLength / itemsPerRow.value))
watchEffect(() => {
    if (!width?.value || !height?.value || !pixiRef.value) return
    itemsPerRow.value = (width.value + gap.value) / (props.size + gap.value)
    if (itemsPerRow.value - Math.floor(itemsPerRow.value) > 0.5) {
        itemsPerRow.value = Math.floor(itemsPerRow.value)
    } else {
        itemsPerRow.value = Math.ceil(itemsPerRow.value)
    }
    size.value = (width.value - (itemsPerRow.value - 1) * gap.value) / itemsPerRow.value
    pixiRef.value.contentHeight = itemsPerCol.value * (size.value + gap.value) - gap.value
})

const renderGraphics = (graphics: Graphics) => {
    graphics.clear()
    Object.entries(props?.data?.chunks || {})?.map(([index, item]) => {
        const { status, createdAt } = item || {}
        const row = Math.floor(Number(index) / itemsPerRow.value)
        const col = Number(index) % itemsPerRow.value
        const x = col * (size.value + gap.value)
        const y = row * (size.value + gap.value)
        graphics.roundRect(x, y, size.value, size.value, 4)
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
