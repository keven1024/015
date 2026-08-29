<script setup lang="ts">
import { cloneVNode, type VNode } from 'vue'
import dayjs from 'dayjs'
import { LucideCheck, LucideX } from '@lucide/vue'
import { cx } from 'class-variance-authority'
import NumberFlow from '@number-flow/vue'

type ShareInfoItem = {
    label: string
    type: 'countdown' | 'string' | 'bool'
    value: any
    handle?: VNode
}

const VNodeRenderer = (props: { vnode: VNode }) => cloneVNode(props.vnode)

const props = defineProps<{
    items: ShareInfoItem[]
}>()

const finalItems = shallowRef<ShareInfoItem[]>([])

const calcCountdownSeconds = (value: number) => {
    if (!value) {
        return 0
    }

    return Math.max(dayjs(value * 10e2).unix() - dayjs().unix(), 0)
}

const calcCountdown = (value: number) => {
    const seconds = calcCountdownSeconds(value)
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = seconds % 60

    return [hours, minutes, secs]
}

const hasCountdown = () => {
    return props.items.some((item) => item.type === 'countdown' && calcCountdownSeconds(item.value) > 0)
}

const calcFinalItems = () => {
    finalItems.value = props.items.map((item) => {
        if (item.type !== 'countdown') {
            return item
        }

        return {
            ...item,
            value: calcCountdown(item.value),
        }
    })
}

const { pause, resume } = useIntervalFn(
    () => {
        calcFinalItems()
        if (!hasCountdown()) {
            pause()
        }
    },
    1000,
    {
        immediate: false,
    }
)

watch(
    () => props.items,
    () => {
        calcFinalItems()
        if (hasCountdown()) {
            resume()
            return
        }
        pause()
    },
    {
        deep: true,
        immediate: true,
    }
)
</script>

<template>
    <div class="flex w-full flex-col gap-2 md:flex-row">
        <div
            v-for="item in finalItems"
            :key="item.label"
            class="flex min-h-11 flex-1 items-center justify-between gap-1 rounded-xl bg-black/5 hover:bg-black/10 px-3 py-2 md:flex-col md:items-start md:justify-between"
        >
            <div class="text-xs font-semibold">{{ item.label }}</div>
            <div class="flex flex-row justify-between w-full">
                <div v-if="item.type === 'bool'" class="flex items-center flex-1">
                    <span class="grid size-7 place-items-center rounded-full bg-white/50">
                        <component
                            :is="item.value ? LucideCheck : LucideX"
                            :class="cx('size-5', item.value ? 'text-emerald-600' : 'text-zinc-500')"
                        />
                    </span>
                </div>
                <div v-else-if="item.type === 'countdown'" class="flex flex-row items-center flex-1">
                    <template v-for="(i, index) in item.value" :key="index">
                        <NumberFlow :value="i" :trend="0" :format="{ minimumIntegerDigits: 2 }" />
                        <span v-if="Number(index) < item.value.length - 1">:</span>
                    </template>
                </div>
                <div v-else class="text-base font-light leading-none tabular-nums md:text-xl flex-1">
                    {{ item.value }}
                </div>
                <VNodeRenderer v-if="item.handle" :vnode="item.handle" />
            </div>
        </div>
    </div>
</template>
