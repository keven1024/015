<script lang="ts" setup>
import { VisSingleContainer, VisDonut } from '@unovis/vue'
import { withDefaults, defineProps } from 'vue'

const props = withDefaults(defineProps<{
    value?: number
    size?: number
    color?: string
}>(), {
    value: 0,
    size: 40,
    color: '#06b6d4'
})
const data = computed(() => {
    const progress = Math.min(Math.max(props.value, 0), 100)
    return [progress, 100 - progress ]
})
const getValue = (d: number) => d
const getColor = (d: number, i: number) => [props?.color, 'transparent'][i]
</script>

<template>
    <div :style="{ width: `${size}px`, height: `${size}px`, '--vis-donut-background-color': 'transparent' }">
        <VisSingleContainer :data="data" :width="size" :height="size">
            <VisDonut :value="getValue" :radius="size / 2" :arcWidth="size / 5" :cornerRadius="20" :color="getColor" />
        </VisSingleContainer>
    </div>
</template>
