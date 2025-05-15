<script setup lang="ts">
import { LucideShare, LucideImage, LucideBot, LucideLanguages, LucideFileText, LucideImageMinus, LucideArrowRightLeft, LucideImagePlus, LucideAudioLines, LucideListMusic } from 'lucide-vue-next'
import { cx } from 'class-variance-authority'
import { isObject } from 'lodash-es';
import showDrawer from '~/lib/showDrawer';
import FileShareHandle from '~/components/Preprocessing/FileShareHandle.vue';

const props = defineProps<{
    hide: () => void
    file: File
    onFileHandle: ({ type, config }: { type: string, config: any }) => void
}>()

const isImage = computed(() => props.file.type.startsWith('image/'))
const isVideo = computed(() => props.file.type.startsWith('video/'))
const isAudio = computed(() => props.file.type.startsWith('audio/'))
const isMedia = computed(() => isImage.value || isVideo.value || isAudio.value)

const isPDF = computed(() => props.file.type.startsWith('application/pdf'))
const isDOC = computed(() => props.file.type.startsWith('application/msword'))
const isXLS = computed(() => props.file.type.startsWith('application/vnd.ms-excel'))
const isPPT = computed(() => props.file.type.startsWith('application/vnd.ms-powerpoint'))
const isDocument = computed(() => isPDF.value || isDOC.value || isXLS.value || isPPT.value)
const actions = [
    {
        label: '分享文件', icon: LucideShare, className: 'bg-green-300', onClick: () => {
            showDrawer({ render: ({ hide }) => h(FileShareHandle, { ...props, hide }) })
        }
    },
    isImage.value && {
        label: '图片压缩', icon: LucideImageMinus, className: 'bg-red-300', onClick: () => {
            console.log('复制链接')
        }
    },
    isImage.value && {
        label: '图片翻译', icon: LucideLanguages, className: 'bg-orange-300', onClick: () => {
            console.log('复制链接')
        }
    },
    isImage.value && {
        label: '图片超分', icon: LucideImagePlus, className: 'bg-cyan-300', onClick: () => {
            console.log('复制链接')
        }
    },
    (isAudio.value || isVideo.value) && {
        label: '转文本', icon: LucideListMusic, className: 'bg-cyan-300', onClick: () => {
            console.log('复制链接')
        }
    },
    isAudio.value && {
        label: '语音克隆', icon: LucideAudioLines, className: 'bg-cyan-300', onClick: () => {
            console.log('复制链接')
        }
    },
    (isDocument.value || isMedia.value) && {
        label: '格式转换', icon: LucideArrowRightLeft, className: 'bg-purple-300', onClick: () => {
            console.log('复制链接')
        }
    },
]?.filter(isObject) as { label: string, icon: any, className: string, onClick: () => void }[]

</script>
<template>
    <div class="flex flex-col gap-5 p-5">
        <div class="flex flex-row gap-5">
            <div v-for="item in actions" :key="item.label" class="flex flex-col items-center gap-2" @click="() => {
                props?.hide()
                item?.onClick()
            }">
                <div :class="cx('size-14 flex justify-center items-center rounded-full', item?.className)">
                    <component :is="item?.icon" />
                </div>
                <div class="text-sm">{{ item?.label }}</div>
            </div>
        </div>
    </div>
</template>
