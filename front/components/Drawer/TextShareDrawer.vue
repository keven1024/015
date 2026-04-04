<script setup lang="ts">
import { cx } from 'class-variance-authority'
import showDrawer from '@/lib/showDrawer'
import TextShareHandle from '@/components/Preprocessing/TextShareHandle.vue'
import { useFeatureMeta, type FeatureKey } from '@/composables/useFeatureMeta'

const props = defineProps<{
    hide: () => void
    text: string
    onTextHandle: ({ type, config }: { type: string; config: any }) => void
}>()

const featureMeta = useFeatureMeta()

type ActionHandler = {
    condition?: () => boolean
    onClick: () => void
}

const actionHandlers: Partial<Record<FeatureKey, ActionHandler>> = {
    'text-share': {
        onClick: () => showDrawer({ render: ({ hide }) => h(TextShareHandle, { ...props, hide }) }),
    },
    // 'text-image-generate': {
    //     label: '生成配图', icon: LucideImage, className: 'bg-red-300',
    //     onClick: () => { console.log('复制链接') }
    // },
    // 'text-ai-ask': {
    //     label: '问大模型', icon: LucideBot, className: 'bg-blue-300',
    //     onClick: () => { console.log('复制链接') }
    // },
    // 'text-translate': {
    //     label: '文本翻译', icon: LucideLanguages, className: 'bg-orange-300',
    //     onClick: () => { console.log('复制链接') }
    // },
}

const actions = computed(() =>
    featureMeta.value
        .filter((meta) => {
            const handler = actionHandlers?.[meta.key]
            return handler && (!handler.condition || handler.condition())
        })
        .map((meta) => ({ ...meta, onClick: actionHandlers[meta.key]!.onClick }))
)
</script>
<template>
    <div class="flex flex-col gap-5 p-5 overflow-x-auto">
        <div class="flex flex-row gap-2">
            <div
                v-for="item in actions"
                :key="item.key"
                class="flex flex-col items-center gap-2 max-w-20"
                @click="
                    () => {
                        props?.hide()
                        item?.onClick()
                    }
                "
            >
                <div :class="cx('size-14 flex justify-center items-center rounded-full mx-3', item?.className)">
                    <component :is="item?.icon" />
                </div>
                <div class="text-xs truncate w-full text-center">{{ item?.label }}</div>
            </div>
        </div>
    </div>
</template>
