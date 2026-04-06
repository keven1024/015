<script setup lang="ts">
import showDrawer from '@/lib/showDrawer'
import FileShareHandle from '@/components/Preprocessing/FileShareHandle.vue'
import ImageConvertHandle from '@/components/Preprocessing/ImageConvertHandle.vue'
import { useFeatureMeta, type FeatureKey } from '@/composables/useFeatureMeta'
import type { FileShareHandleProps } from '../Preprocessing/types'

const props = defineProps<{
    hide: () => void
    file: File[]
    onFileHandle: (props: FileShareHandleProps) => void
}>()

const isImage = computed(() => props.file.every((r) => r?.type?.startsWith('image/')))

const featureMeta = useFeatureMeta()

type ActionHandler = {
    condition?: () => boolean
    onClick: () => void
}

const actionHandlers: Partial<Record<FeatureKey, ActionHandler>> = {
    'file-share': {
        onClick: () => showDrawer({ render: ({ hide }) => h(FileShareHandle, { ...props, hide }) }),
    },
    'file-image-compress': {
        condition: () => isImage.value,
        onClick: () => props.onFileHandle({ type: 'file-image-compress', config: {} }),
    },
    'file-image-convert': {
        condition: () => isImage.value,
        onClick: () => showDrawer({ render: ({ hide }) => h(ImageConvertHandle, { ...props, hide }) }),
    },
    // isImage.value && {
    //     label: '图片翻译', icon: LucideLanguages, className: 'bg-orange-300', onClick: () => {
    //         console.log('复制链接')
    //     }
    // },
    // isImage.value && {
    //     label: '图片超分', icon: LucideImagePlus, className: 'bg-cyan-300', onClick: () => {
    //         console.log('复制链接')
    //     }
    // },
    // (isAudio.value || isVideo.value) && {
    //     label: '转文本', icon: LucideListMusic, className: 'bg-cyan-300', onClick: () => {
    //         console.log('复制链接')
    //     }
    // },
    // isAudio.value && {
    //     label: '语音克隆', icon: LucideAudioLines, className: 'bg-cyan-300', onClick: () => {
    //         console.log('复制链接')
    //     }
    // },
    // (isDocument.value || isMedia.value) && {
    //     label: '格式转换', icon: LucideArrowRightLeft, className: 'bg-purple-300', onClick: () => {
    //         console.log('复制链接')
    //     }
    // },
}

const actions = computed(() =>
    featureMeta.value
        .filter((meta) => {
            const { key } = meta || {}
            const handler = actionHandlers?.[key]
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
                <div class="size-14 flex justify-center items-center rounded-full mx-3" :style="item?.style">
                    <component :is="item?.icon" />
                </div>
                <div class="text-xs truncate w-full text-center">{{ item?.label }}</div>
            </div>
        </div>
    </div>
</template>
