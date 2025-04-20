<script setup lang="ts">
import { LucideShare, LucideImage, LucideBot, LucideLanguages } from 'lucide-vue-next'
import { cx } from 'class-variance-authority'

const props = defineProps<{
    hide: () => void
}>()
const actions = [
    {
        label: '分享文件', icon: LucideShare, className: 'bg-green-300', onClick: () => {
            console.log('复制链接')
        }
    },
    {
        label: '生成图文', icon: LucideImage, className: 'bg-red-300', onClick: () => {
            console.log('复制链接')
        }
    },
    {
        label: '问大模型', icon: LucideBot, className: 'bg-blue-300', onClick: () => {
            console.log('复制链接')
        }
    },
    {
        label: '文本翻译', icon: LucideLanguages, className: 'bg-orange-300', onClick: () => {
            console.log('复制链接')
        }
    },
]
</script>
<template>
    <div class="flex flex-col gap-5 p-5">
        <div class="flex flex-row gap-5">
            <div v-for="item in actions" :key="item.label" class="flex flex-col items-center gap-2" @click="()=>{
                item?.onClick()
                props?.hide()
            }">
                <div :class="cx('size-14 flex justify-center items-center rounded-full', item?.className)">
                    <component :is="item?.icon" />
                </div>
                <div class="text-sm">{{ item?.label }}</div>
            </div>
        </div>
    </div>
</template>
