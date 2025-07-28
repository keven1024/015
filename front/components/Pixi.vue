<script setup lang="ts">
import { Application } from 'vue3-pixi'
import { useElementSize } from '@vueuse/core'

export interface PixiExpose {
    app: Application
    width: number
    height: number
    contentHeight: number
    ref: HTMLElement
}

const containerRef = ref<HTMLElement>()
const app = ref<Application>()
const { width: windowWidth, height: windowHeight } = useWindowSize()
const { width: containerWidth, height: containerHeight } = useElementSize(containerRef)
const width = computed(() => Math.min(containerWidth.value, windowWidth.value))
const contentHeight = ref(0)
const height = computed(() => {
    if (!contentHeight.value) return Math.min(containerHeight.value, windowHeight.value)
    if (containerHeight.value > contentHeight.value) return contentHeight.value
    return contentHeight.value
})
const resolution = computed(() => window?.devicePixelRatio || 1)

watchEffect(() => {
    if (app.value?.app?.renderer) {
        app.value?.app?.renderer?.resize(width.value, height.value)
    }
})

defineExpose({
    app,
    width,
    height,
    contentHeight,
    ref: containerRef,
})
</script>

<template>
    <div ref="containerRef" class="w-full h-full">
        <Application :width="width" :height="height" :resolution="resolution" ref="app" autoDensity antialias :backgroundAlpha="0">
            <slot />
        </Application>
    </div>
</template>
