<script setup lang="ts">
import { Button } from '@/components/ui/button'
import FilePreviewView from '@/components/FilePreviewView.vue'
import { Input } from '@/components/ui/input'
import { useClipboard } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { useQuery } from '@tanstack/vue-query'
import useMyAppShare from '@/composables/useMyAppShare'
import useMyAppConfig from '@/composables/useMyAppConfig'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn' // 导入中文语言包
import showDrawer from '@/lib/showDrawer'
import QrCoreDrawer from '@/components/Drawer/QrCoreDrawer.vue'
import { h } from 'vue'
import { cx } from 'class-variance-authority'
dayjs.extend(relativeTime) // 扩展 relativeTime 插件
dayjs.locale('zh-cn') // 设置语言为中文

const props = defineProps<{
    data: { files: { id: string; file: File }[]; config: any; handle_type: string }
}>()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const { createFileShare } = useMyAppShare()
const { data } = useQuery({
    queryKey: ['create-share', ...props?.data?.files?.map((item) => item.id)],
    queryFn: async () => {
        const { files, config } = props?.data || {}
        const data = await createFileShare({
            files: files?.map((item) => {
                const { id, file } = item || {}
                return { id, name: file.name }
            }),
            config,
        })
        return data?.map((item) => item?.data)
    },
})

const selectedFile = ref<string | undefined>()
const selectedFileShare = computed(() => {
    return data?.value?.find((item) => item?.id === selectedFile.value)
})
watchEffect(() => {
    if (data?.value && data?.value?.length === 1 && !!data?.value?.[0]?.id) {
        selectedFile.value = data.value[0].id
    }
})

watchEffect(() => {
    console.log('data', data?.value)
})

const appConfig = useMyAppConfig()
const getShareUrl = (id: string) => {
    return `${appConfig?.value?.site_url}/s/${id}`
}

const { copy } = useClipboard()
</script>

<template>
    <div class="flex flex-col gap-3">
        <h2 class="text-lg">上传成功</h2>
        <div class="flex flex-col gap-3 items-center">
            <div v-if="data?.length === 1" class="flex flex-col h-30 items-center">
                <FilePreviewView :value="props?.data?.files?.[0]?.file" />
            </div>
            <div v-else class="flex flex-col gap-2 w-full p-5 bg-white/20 backdrop-blur-xl rounded-md">
                <div class="text-sm font-semibold">文件列表</div>
                <div
                    v-for="file in data"
                    :class="
                        cx(
                            'flex flex-row justify-between items-center gap-1 rounded-md p-2 border border-black/10 w-full cursor-pointer',
                            selectedFile === file?.id && 'bg-primary text-white'
                        )
                    "
                    @click="selectedFile = file?.id"
                >
                    <div class="flex flex-row items-center gap-2 flex-1 min-w-0">
                        <FileIcon
                            :file="props?.data?.files?.[data?.findIndex((i) => i?.id === file?.id) as number]?.file"
                            :class="cx('!size-7 !rounded-md shrink-0', selectedFile === file?.id && '!bg-white/50')"
                        />
                        <div class="text-sm flex-1 truncate">{{ file?.file_name }}</div>
                    </div>
                    <div class="flex flex-row items-center gap-2 shrink-0">
                        <Button
                            variant="outline"
                            :class="cx('bg-white/70', selectedFile === file?.id && '!bg-white/30 border-none hover:text-white/80')"
                            size="icon"
                            @click="
                                () => {
                                    copy(getShareUrl(file?.id as string))
                                    toast.success('复制成功')
                                }
                            "
                        >
                            <LucideCopy />
                        </Button>
                        <Button
                            variant="outline"
                            :class="cx('bg-white/70', selectedFile === file?.id && '!bg-white/30 border-none hover:text-white/80')"
                            size="icon"
                            @click="
                                () => {
                                    showDrawer({
                                        render: ({ ...rest }) =>
                                            h(QrCoreDrawer, {
                                                ...rest,
                                                data: getShareUrl(file?.id as string),
                                            }),
                                    })
                                }
                            "
                        >
                            <LucideQrCode />
                        </Button>
                    </div>
                </div>
            </div>
            <div v-if="!!selectedFileShare" class="flex flex-col md:flex-row gap-5 rounded-md p-5 bg-white/20 backdrop-blur-xl w-full">
                <div class="flex flex-col gap-2 flex-1">
                    <div class="text-sm font-semibold">信息</div>
                    <div class="grid grid-cols-2 gap-2">
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">下载次数</div>
                            <div class="text-3xl font-light">{{ selectedFileShare?.download_nums }}</div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                            <div class="text-xs font-semibold">过期时间</div>
                            <div class="text-md font-light">
                                {{ dayjs((selectedFileShare?.expire_at || 0) * 1000).format('YYYY-MM-DD HH:mm:ss') }}
                            </div>
                        </div>
                        <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1" v-if="selectedFileShare?.pickup_code">
                            <div class="flex flex-row justify-between w-full items-center">
                                <div class="text-xs font-semibold">提取码</div>
                                <Button
                                    variant="outline"
                                    class="bg-white/70 p-0 size-6"
                                    size="icon"
                                    @click="
                                        () => {
                                            copy(selectedFileShare?.pickup_code as string)
                                            toast.success('复制成功')
                                        }
                                    "
                                >
                                    <LucideCopy class="size-3" />
                                </Button>
                            </div>
                            <div class="flex flex-row gap-2">
                                <div v-for="s in selectedFileShare?.pickup_code" class="text-2xl font-light">
                                    {{ s }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex flex-col gap-5 flex-1">
                    <div class="text-sm font-semibold">链接</div>
                    <div class="flex flex-row gap-2">
                        <Input :model-value="getShareUrl(selectedFileShare?.id as string)" class="bg-white/70" readonly />
                        <Button
                            variant="outline"
                            class="bg-white/70"
                            size="icon"
                            @click="
                                () => {
                                    copy(getShareUrl(selectedFileShare?.id as string))
                                    toast.success('复制成功')
                                }
                            "
                        >
                            <LucideCopy />
                        </Button>

                        <Button
                            variant="outline"
                            class="bg-white/70"
                            size="icon"
                            @click="
                                () => {
                                    showDrawer({
                                        render: ({ ...rest }) =>
                                            h(QrCoreDrawer, {
                                                ...rest,
                                                data: getShareUrl(selectedFileShare?.id as string),
                                            }),
                                    })
                                }
                            "
                        >
                            <LucideQrCode />
                        </Button>
                    </div>
                </div>
            </div>
            <Button
                class="w-40 hover:bg-primary/90"
                @click="
                    () => {
                        emit('change', 'input')
                    }
                "
            >
                返回首页
            </Button>
        </div>
    </div>
</template>
