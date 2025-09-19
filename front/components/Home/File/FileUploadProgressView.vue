<script setup lang="ts">
import { LucidePlay, LucideSettings, LucideSquare } from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import dayjs from 'dayjs'
import FileUploadBlockProgressView from '@/components/FileUploadBlockProgressView.vue'
import { motion } from 'motion-v'
import { filesize } from 'filesize'
import { cx } from 'class-variance-authority'
import asyncWait from '@/lib/asyncWait'
import asyncWorker from '@/lib/asyncWorker'
import calcFileHashWorker from '@/lib/calcFileHashWorker?worker'
import { clamp, get, isEmpty, isNumber, isString, sample, shuffle, times } from 'lodash-es'
import { nanoid } from 'nanoid'
import asyncRetry from '@/lib/asyncRetry'
import asyncTimeout from '@/lib/asyncTimeout'

const props = defineProps<{
    data: { file: File[]; config: any; handle_type: string }
}>()

const selectedFile = ref()
const uploadfiles = ref<
    {
        fileId: string
        file: File
        status: 'start' | 'pause' | 'finish' | 'error'
        hash?: string
        procressType: 'hash' | 'create' | 'upload'
        uploadInfo?: {
            chunks: Record<number, { status: 'success' | 'error' | 'processing'; createdAt: number }>
            chunkLength: number
        }
        queue: {
            taskId: string
            taskType: 'hash' | 'create' | 'chunk' | 'upload'
            queueType: 'sync' | 'async' // sync任务禁止并发
            index?: number
            retry?: number
        }[]
    }[]
>([])
const selectedUploadfile = computed(() => uploadfiles.value.find((item) => item.fileId === selectedFile.value))
const selectedUploadfileChunk = computed(() => Object.values(selectedUploadfile.value?.uploadInfo?.chunks || {}))
const selectedUploadfileViewMode = ref<'progress' | 'heatmap'>('progress')

const procressTaskList = ref<Map<string, any>>(new Map())
const taskList = computed(() => uploadfiles.value.filter((r) => r.queue.length > 0 && r.status === 'start').flatMap((r) => r.queue))
const batchNum = ref(3)

onMounted(() => {
    props.data.file.forEach((file) => {
        uploadfiles.value.push({
            fileId: `${file.name}_${file.size}`,
            file,
            status: 'start',
            procressType: 'hash',
            queue: [
                { taskType: 'hash', queueType: 'sync', taskId: nanoid() },
                { taskType: 'create', queueType: 'sync', taskId: nanoid() },
                { taskType: 'chunk', queueType: 'sync', taskId: nanoid() },
            ],
        })
    })
})

const { error, execute, isLoading } = useAsyncState(
    async () => {
        while (taskList.value.length > 0) {
            console.log('hasTask')
            const taskList = uploadfiles?.value?.filter((r) => r.queue.length > 0 && r.status === 'start')
            await Promise.all(
                times(batchNum.value, async (i: number) => {
                    const file = sample(taskList)
                    const fileId = file?.fileId as string
                    const task = get(file?.queue, '0')
                    if (!task) return
                    const { taskId, index, queueType, taskType } = task || {}
                    if (procressTaskList.value.has(taskId)) return
                    procressTaskList.value.set(taskId, task)

                    const uploadFileIndex = uploadfiles.value.findIndex((r) => r.fileId === file?.fileId)

                    if (queueType === 'async') {
                        uploadfiles.value[uploadFileIndex]?.queue.shift()
                    }
                    console.log('[start]', taskType, taskId, queueType)
                    try {
                        if (taskType === 'hash') {
                            await asyncRetry(() => asyncTimeout(() => handleHash(fileId), 10000))
                        }
                        if (taskType === 'create') {
                            await asyncRetry(() => asyncTimeout(() => handleCreate(fileId), 10000))
                        }
                        if (taskType === 'chunk') {
                            await asyncRetry(() => asyncTimeout(() => handleChunk(fileId), 10000))
                        }
                        if (taskType === 'upload' && isNumber(index)) {
                            await asyncRetry(() => asyncTimeout(() => handleUpload(fileId, index), 10000))
                        }
                        console.log('[finish]', taskType, taskId)
                        if (queueType === 'sync') {
                            uploadfiles.value[uploadFileIndex]?.queue.shift()
                        }
                    } catch (error) {
                        console.log('error', error)
                        // todo 重新加入队列
                        if (queueType === 'async') {
                            uploadfiles.value[uploadFileIndex]?.queue.push({ ...task, retry: (task?.retry || 0) + 1 })
                        }
                        if (queueType === 'sync') {
                            uploadfiles.value[uploadFileIndex].status = 'error'
                        }
                    } finally {
                        procressTaskList.value.delete(taskId)
                    }
                })
            )
        }
    },
    null,
    { immediate: false }
)

watchEffect(async () => {
    if (taskList.value.length > 0 && !isLoading.value) {
        console.log('task队列已更新', `当前任务 ${taskList.value.length} 个 正在消化`)
        execute()
        console.log('开始执行任务...')
    }
})

const handleHash = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    uploadfile.procressType = 'hash'
    const res = await asyncWorker(calcFileHashWorker, { data: { file: uploadfile.file } })
    const { hash } = res?.data || {}
    uploadfile.hash = hash
}

const handleCreate = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    uploadfile.procressType = 'create'
    await asyncWait(1000)
    console.log('create', fileId)
}

const handleChunk = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    uploadfile.procressType = 'upload'
    uploadfile.uploadInfo = {
        chunks: {},
        chunkLength: 1000,
    }
    const tasks = shuffle(times(1000, (i: number) => ({ taskType: 'upload' as const, queueType: 'async' as const, taskId: nanoid(), index: i })))
    uploadfile.queue.push(...tasks)
    console.log('chunk', fileId)
}

const handleUpload = async (fileId: string, index: number) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    if (uploadfile.procressType !== 'upload') uploadfile.procressType = 'upload'
    let chunkInfo = uploadfile.uploadInfo?.chunks?.[index]
    if (!chunkInfo) {
        uploadfile.uploadInfo!.chunks[index] = {
            status: 'processing',
            createdAt: Date.now(),
        }
    }
    await asyncWait(500)
    if (index % 3 === 0) {
        uploadfile.uploadInfo!.chunks[index].status = 'error'
        throw new Error('上传失败')
    }
    uploadfile.uploadInfo!.chunks[index].status = 'success'
    console.log('upload', fileId, index)
}
</script>

<template>
    <div class="grid grid-cols-4 gap-5">
        <div class="rounded-xl p-3 bg-white/80 flex flex-col gap-2 col-span-3">
            <div class="flex flex-col gap-1">
                <div class="text-xs opacity-70">总上传速度</div>
                <div class="text-2xl font-bold">144.0MB/s</div>
            </div>
            <div class="flex-1 relative overflow-hidden flex flex-row gap-0.5 justify-end items-end">
                <motion.div
                    class="w-2 shrink-0 bg-primary relative"
                    :style="{ height: (i.value + 1) * 4 + 'px' }"
                    :layoutId="i.time"
                    v-for="i in data"
                    :key="i.time"
                    :initial="{ x: 10, opacity: 0 }"
                    :animate="{ x: 0, opacity: 1 }"
                    :transition="{ duration: 1 }"
                >
                </motion.div>
                <!-- <BarChart class="h-full" :data="data" index="time" :categories="['value']" :showTooltip="false"
            :showLegend="false" :showXAxis="false" :showYAxis="false" :showGrid="false" :groupMaxWidth="10" /> -->
            </div>
        </div>
        <div class="rounded-xl bg-white/80 aspect-square flex flex-col gap-2 relative overflow-hidden">
            <div class="absolute top-0 left-0 w-full h-full z-[0] flex flex-col justify-end">
                <div class="w-full bg-green-100 h-1/2 border-t border-green-500"></div>
            </div>
            <div class="flex flex-col gap-2 justify-between p-3 h-full relative z-[1]">
                <div class="flex flex-col gap-1">
                    <div class="text-xs opacity-70">总上传进度</div>
                    <div class="text-4xl font-bold">44.0%</div>
                </div>
                <div class="flex flex-row gap-2">
                    <Button class="aspect-square bg-green-200 hover:bg-green-300 text-green-500 hover:text-white p-0">
                        <LucidePlay class="size-1/2" />
                    </Button>
                    <Button
                        class="aspect-square bg-red-200 hover:bg-red-300 text-red-500 hover:text-white p-0"
                        @click="
                            () => {
                                console.log('暂停')
                            }
                        "
                    >
                        <LucideSquare class="size-1/2" />
                    </Button>
                    <Button class="aspect-square bg-blue-200 hover:bg-blue-300 text-blue-500 hover:text-white p-0">
                        <LucideSettings class="size-1/2" />
                    </Button>
                </div>
            </div>
        </div>
        <div class="col-span-4 flex flex-col bg-white/80 rounded-xl p-3 text-md gap-5">
            <div>文件列表</div>
            <div class="flex flex-col -mx-3 text-sm">
                <div class="grid grid-cols-[2fr_5rem_5rem_4fr] gap-2 border-b border-black/20 pb-2 px-3">
                    <div>文件名</div>
                    <div>文件大小</div>
                    <div>上传速度</div>
                    <div>进度</div>
                </div>
                <div
                    :class="
                        cx(
                            'grid grid-cols-[2fr_5rem_5rem_4fr] gap-2 py-2 border-b border-black/20 items-center hover:bg-primary/30 px-3 cursor-pointer',
                            selectedFile === item?.fileId && 'bg-primary text-white hover:!bg-primary'
                        )
                    "
                    v-for="(item, index) in uploadfiles"
                    @click="
                        () => {
                            if (selectedFile === item?.fileId) {
                                selectedFile = null
                            } else {
                                selectedFile = item?.fileId
                            }
                        }
                    "
                >
                    <div class="flex flex-row gap-2 items-center grow-0 overflow-hidden">
                        <Button
                            class="size-8 p-0 hover:bg-white/50"
                            variant="ghost"
                            @click="
                                (e: Event) => {
                                    e.stopPropagation()
                                    if (item?.status === 'start') {
                                        uploadfiles[index].status = 'pause'
                                    } else {
                                        uploadfiles[index].status = 'start'
                                    }
                                }
                            "
                        >
                            <LucidePlay class="size-4 text-green-500" v-if="item?.status === 'start'" />
                            <LucidePause class="size-4 text-orange-500" v-if="item?.status === 'pause'" />
                        </Button>
                        <div class="truncate">{{ item?.file?.name }}</div>
                    </div>
                    <div>{{ filesize(item?.file?.size) }}</div>
                    <div>100MB/s</div>
                    <div class="flex flex-row gap-2 items-center" v-if="item?.procressType !== 'upload'">
                        <LucideLoaderCircle class="size-4 animate-spin" />
                        <div>正在{{ item?.procressType }}中...</div>
                    </div>
                    <div class="flex flex-row gap-2 items-center" v-else>
                        <div class="rounded-full bg-white/50 w-full h-2 overflow-hidden border border-white">
                            <div
                                class="h-full bg-primary"
                                :style="`width: ${clamp(((Object.keys(item?.uploadInfo?.chunks || {})?.length || 0) / (item?.uploadInfo?.chunkLength || 0)) * 100, 0, 100)}%`"
                            ></div>
                        </div>
                        <div>
                            {{
                                clamp(
                                    ((Object.keys(item?.uploadInfo?.chunks || {})?.length || 0) / (item?.uploadInfo?.chunkLength || 0)) * 100,
                                    0,
                                    100
                                )?.toFixed(2)
                            }}%
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-span-4 flex flex-col bg-white/80 rounded-xl p-3 gap-5" v-if="selectedFile">
            <div>上传详情</div>
            <div class="grid grid-cols-3 text-sm gap-3">
                <div>区块： {{ selectedUploadfile?.uploadInfo?.chunkLength }} x 256.0 KiB</div>
                <div class="truncate col-span-2">hash: {{ selectedUploadfile?.hash }}</div>
                <div>已完成: {{ selectedUploadfileChunk?.filter((r) => r.status === 'success')?.length || 0 }}</div>
                <div>已丢弃: {{ selectedUploadfileChunk?.filter((r) => r.status === 'error')?.length || 0 }}</div>
                <div>
                    待完成:
                    {{ (selectedUploadfile?.uploadInfo?.chunkLength || 0) - (selectedUploadfileChunk?.length || 0) }}
                </div>
                <div class="col-span-3 flex flex-row justify-between items-center">
                    <div class="text-md font-bold">{{ selectedUploadfileViewMode === 'progress' ? '区块进度条' : '区块热力图' }}</div>
                    <Button
                        size="sm"
                        class="ml-auto text-xs"
                        @click="selectedUploadfileViewMode = selectedUploadfileViewMode === 'progress' ? 'heatmap' : 'progress'"
                    >
                        <LucideArrowDownUp class="size-4" />
                        {{ selectedUploadfileViewMode === 'progress' ? '热力图' : '进度条' }}
                    </Button>
                </div>
                <div class="h-7 col-span-3 flex flex-row gap-2 items-center" v-if="selectedUploadfileViewMode === 'progress'">
                    <div class="flex-1 h-full">
                        <FileUploadBlockProgressView :data="selectedUploadfile?.uploadInfo" />
                    </div>
                    {{
                        clamp(((selectedUploadfileChunk?.length || 0) / (selectedUploadfile?.uploadInfo?.chunkLength || 0)) * 100, 0, 100)?.toFixed(
                            2
                        )
                    }}%
                </div>
                <div v-if="selectedUploadfileViewMode === 'heatmap'" class="col-span-3 bg-black/5 rounded p-2">
                    <FileUploadHeatMapView :size="12" :data="selectedUploadfile?.uploadInfo" />
                </div>
            </div>
        </div>
    </div>
</template>
