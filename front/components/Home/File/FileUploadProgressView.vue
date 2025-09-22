<script setup lang="ts">
import { LucidePlay, LucideSettings, LucideSquare } from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import FileUploadBlockProgressView from '@/components/FileUploadBlockProgressView.vue'
import { motion } from 'motion-v'
import getFileSize from '~/lib/getFileSize'
import { cx } from 'class-variance-authority'
import asyncWorker from '@/lib/asyncWorker'
import calcFileHashWorker from '@/lib/calcFileHashWorker?worker'
import { clamp, get, groupBy, isEmpty, isNumber, isString, sample, shuffle, times } from 'lodash-es'
import { nanoid } from 'nanoid'
import { toast } from 'vue-sonner'
import dayjs from 'dayjs'
import showDrawer from '~/lib/showDrawer'
import FileUploadSpeedInfoView from './FileUploadSpeedInfoView.vue'
import getFileChunk from '~/lib/getFileChunk'
import type { FileHandleKey } from '~/components/Preprocessing/types'

const props = defineProps<{
    data: { file: File[]; config: Record<string, any>; handle_type: FileHandleKey }
}>()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const form = useFormContext()

const selectedFile = ref()
const uploadfiles = ref<
    {
        fileId: string
        id?: string // 后端文件id
        file: File
        status: 'start' | 'pause' | 'finish' | 'error'
        hash?: string
        procressType: 'hash' | 'create' | 'upload' | 'finish'
        uploadInfo?: {
            chunks: Record<number, { status: 'success' | 'error' | 'processing'; createdAt: number }>
            chunkLength: number
            ChunkSize: number
        }
        queue: {
            taskId: string
            taskType: 'hash' | 'create' | 'chunk' | 'upload' | 'finish'
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
const activeTaskList = computed(() => uploadfiles.value.filter((r) => r.queue.length > 0 && r.status === 'start'))
const activeTaskAllQueue = computed(() => activeTaskList.value.flatMap((r) => r.queue))
const batchNum = ref(3)

const totalTaskStatus = computed(() => {
    if (uploadfiles.value.some((r) => r.status === 'start')) {
        return 'start'
    }
    if (uploadfiles.value.some((r) => r.status === 'pause')) {
        return 'pause'
    }
    return 'disabled'
})
const totalUploadProgress = computed(() => {
    const successCount = uploadfiles.value.reduce((acc, curr) => {
        const { status, uploadInfo } = curr || {}
        if (status === 'finish') return acc
        const { chunks } = uploadInfo || {}
        return acc + Object.entries(chunks || {}).filter(([index, chunk]) => chunk.status === 'success').length
    }, 0)
    const totalCount = uploadfiles.value.reduce((acc, curr) => {
        const { status, uploadInfo } = curr || {}
        if (status === 'finish') return acc
        const { chunkLength } = uploadInfo || {}
        return acc + (chunkLength || 0)
    }, 0)
    return ((successCount || 0) / (totalCount || 0)) * 100
})

const counter = useInterval(1000)
const speedChartData = ref<Record<number, { createdAt: number; value: number }[]>>({})
watch(counter, () => {
    const speed = uploadfiles.value?.flatMap((item) => {
        const { chunks, ChunkSize } = item?.uploadInfo || {}
        return Object.entries(chunks || {})
            .filter(([index, chunk]) => chunk.status === 'success' && dayjs().unix() - 60 < chunk.createdAt)
            ?.map(([index, chunk]) => {
                const { createdAt } = chunk || {}
                return {
                    createdAt,
                    value: ChunkSize || 0,
                }
            })
    })
    speedChartData.value = groupBy(speed, 'createdAt')
})

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
        while (activeTaskAllQueue.value.length > 0) {
            const taskList = uploadfiles?.value?.filter((r) => r.queue.length > 0 && r.status === 'start')
            await Promise.all(
                times(batchNum.value, async (i: number) => {
                    const file = sample(taskList)
                    const { fileId, queue } = file || {}
                    if (!fileId) return
                    const task = get(queue, '0')
                    if (!task) return
                    const { taskId, index, queueType, taskType, retry } = task || {}
                    if (procressTaskList.value.has(taskId)) return
                    procressTaskList.value.set(taskId, task)

                    const uploadFileIndex = uploadfiles.value.findIndex((r) => r.fileId === file?.fileId)

                    if (queueType === 'async') {
                        if (!!retry && retry >= 3) {
                            toast.error('上传错误', {
                                description: `文件 ${file?.file?.name} 的${index}分块经过多次尝试依然上传失败, 我们已经终止该文件上传`,
                            })
                            uploadfiles.value[uploadFileIndex].status = 'error'
                            return
                        }
                        uploadfiles.value[uploadFileIndex]?.queue.shift()
                    }
                    // console.log('[start]', taskType, taskId, queueType)
                    try {
                        if (taskType === 'hash') {
                            await handleHash(fileId)
                        }
                        if (taskType === 'create') {
                            await handleCreate(fileId)
                        }
                        if (taskType === 'chunk') {
                            await handleChunk(fileId)
                        }
                        if (taskType === 'upload' && isNumber(index)) {
                            await handleUpload(fileId, index)
                        }
                        if (taskType === 'finish') {
                            await handleFinish(fileId)
                        }
                        // console.log('[finish]', taskType, taskId)
                        if (queueType === 'sync') {
                            uploadfiles.value[uploadFileIndex]?.queue.shift()
                        }
                    } catch (error) {
                        console.log('error', error)
                        // todo 重新加入队列
                        if (queueType === 'async') {
                            uploadfiles.value[uploadFileIndex]?.queue.push({ ...task, retry: (task?.retry || 0) + 1 })
                            toast.warning('上传错误', {
                                description: `文件 ${file?.file?.name} 的${index}分块上传失败, 我们将在稍后再次尝试上传`,
                            })
                        }
                        if (queueType === 'sync') {
                            uploadfiles.value[uploadFileIndex].status = 'error'
                            toast.error('上传错误', {
                                description: `文件${file?.file?.name}上传失败, 请重试`,
                            })
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
    if (activeTaskAllQueue.value.length > 0 && !isLoading.value) {
        execute()
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
    const { hash, file } = uploadfile || {}
    if (!hash) return
    uploadfile.procressType = 'create'
    const { size, type = 'application/octet-stream' } = file || {}
    const createData = await $fetch<{
        data: {
            id: string
            type: 'init' | 'already'
            chunk_size: number
        }
    }>('/api/file/create', {
        method: 'POST',
        body: {
            size,
            mime_type: type || 'application/octet-stream',
            hash,
        },
    })
    const { id, chunk_size, type: createType } = createData?.data || {}
    uploadfile.id = id
    uploadfile.uploadInfo = {
        chunks: {},
        chunkLength: Math.ceil(size / chunk_size),
        ChunkSize: chunk_size,
    }
    if (createType !== 'init') {
        // 文件存在
        uploadfile.status = 'finish'
        uploadfile.procressType = 'finish'
        uploadfile.queue = []
        return
    }
}

const handleChunk = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    const { chunkLength } = uploadfile.uploadInfo || {}
    if (!chunkLength) return

    uploadfile.procressType = 'upload'
    const tasks = shuffle(
        times(chunkLength, (i: number) => ({ taskType: 'upload' as const, queueType: 'async' as const, taskId: nanoid(), index: i }))
    )
    uploadfile.queue.push(...tasks)
}

const handleUpload = async (fileId: string, index: number) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    const { id, uploadInfo } = uploadfile || {}
    const { chunkLength, ChunkSize } = uploadInfo || {}
    if (!ChunkSize || !id) return

    if (uploadfile.procressType !== 'upload') uploadfile.procressType = 'upload'
    let chunkInfo = uploadfile.uploadInfo?.chunks?.[index]
    if (!chunkInfo) {
        uploadfile.uploadInfo!.chunks[index] = {
            status: 'processing',
            createdAt: dayjs().unix(),
        }
    }

    const chunk = await getFileChunk(uploadfile.file, index * ChunkSize, ChunkSize)
    const formData = new FormData()
    formData.append('file', new Blob([chunk]))
    formData.append('index', (index + 1).toString())
    formData.append('id', id)
    const res = await $fetch<{
        code: number
    }>('/api/file/slice', {
        method: 'POST',
        body: formData,
    })
    const { code } = res || {}
    if (code !== 200) {
        throw new Error('上传失败')
    }
    uploadfile.uploadInfo!.chunks[index].status = 'success'
    if (Object.entries(uploadfile.uploadInfo!.chunks || {}).filter(([index, chunk]) => chunk.status === 'success').length === chunkLength) {
        uploadfile.queue.push({ taskType: 'finish', queueType: 'sync', taskId: nanoid() })
    }
}

const handleFinish = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    uploadfile.procressType = 'finish'
    const { id } = uploadfile || {}
    const res = await $fetch<{
        code: number
    }>('/api/file/finish', {
        method: 'POST',
        body: {
            id,
        },
    })
    if (res?.code !== 200) {
        throw new Error('上传失败')
    }
    uploadfile.status = 'finish'
}

watchEffect(() => {
    if (!isEmpty(uploadfiles.value) && uploadfiles.value.every((r) => r.status === 'finish')) {
        // console.log('change', uploadfiles.value)
        form.setFieldValue(
            'files',
            uploadfiles.value.map((item) => {
                const { id, file } = item || {}
                return {
                    id,
                    file,
                }
            })
        )
        emit('change', 'result')
    }
})

onUnmounted(() => {
    uploadfiles.value = []
})

const handleShowSpeedInfo = () => {
    showDrawer({
        render: ({ hide }) => h(FileUploadSpeedInfoView, { hide }),
    })
}
</script>

<template>
    <div class="grid grid-cols-4 gap-5">
        <div class="rounded-xl p-3 bg-white/80 flex flex-col gap-2 col-span-4 md:col-span-3 h-32 md:h-auto">
            <div class="flex flex-col gap-1">
                <div @click="handleShowSpeedInfo" class="flex flex-row gap-1 items-center text-xs opacity-70">
                    总上传进度
                    <LucideInfo class="size-3" />
                </div>
                <div class="text-2xl font-bold">
                    {{
                        `${
                            getFileSize(
                                Object.entries(speedChartData)
                                    ?.filter((r) => dayjs().unix() - 60 < parseInt(r[0]))
                                    ?.reduce((acc, curr) => acc + curr[1]?.reduce((_acc, _curr) => _acc + _curr.value, 0), 0) / 60
                            ) ?? 0
                        }/s`
                    }}
                </div>
            </div>
            <div class="flex-1 relative overflow-hidden flex flex-row gap-0.5 justify-end items-end">
                <motion.div
                    class="w-2 shrink-0 bg-primary relative"
                    :style="{
                        height: `${(i[1]?.reduce((acc, curr) => acc + curr.value, 0) / Math.max(...Object.entries(speedChartData)?.map((r) => r[1]?.reduce((acc, curr) => acc + curr.value, 0)))) * 100}%`,
                    }"
                    :layoutId="i[0]"
                    v-for="i in Object.entries(speedChartData)"
                    :key="i[0]"
                    :initial="{ x: 10, opacity: 0 }"
                    :animate="{ x: 0, opacity: 1 }"
                    :transition="{ duration: 1 }"
                >
                </motion.div>
                <!-- <BarChart class="h-full" :data="data" index="time" :categories="['value']" :showTooltip="false"
            :showLegend="false" :showXAxis="false" :showYAxis="false" :showGrid="false" :groupMaxWidth="10" /> -->
            </div>
        </div>
        <div class="rounded-xl col-span-4 md:col-span-1 bg-white/80 h-32 md:h-auto md:aspect-square flex flex-col gap-2 relative overflow-hidden">
            <div class="absolute top-0 left-0 w-full h-full z-[0] flex flex-col justify-end">
                <div class="w-full bg-green-100 border-t border-green-500" :style="`height: ${totalUploadProgress || 0}%`"></div>
            </div>
            <div class="flex flex-col gap-2 justify-between p-3 h-full relative z-[1]">
                <div class="flex flex-col gap-1">
                    <div class="text-xs opacity-70">总上传进度</div>
                    <div class="text-4xl font-bold">{{ (totalUploadProgress || 0).toFixed(1) }}%</div>
                </div>
                <div class="flex flex-row gap-2">
                    <Button
                        class="aspect-square hover:text-white p-0 bg-green-200 hover:bg-green-300 text-green-500"
                        :disabled="['start', 'disabled'].includes(totalTaskStatus)"
                        @click="
                            () => {
                                uploadfiles.forEach((r) => {
                                    r.status = 'start'
                                })
                            }
                        "
                    >
                        <LucidePlay class="size-1/2" />
                    </Button>
                    <Button
                        class="aspect-square hover:text-white p-0 bg-orange-200 hover:bg-orange-300 text-orange-500"
                        :disabled="['pause', 'disabled'].includes(totalTaskStatus)"
                        @click="
                            () => {
                                uploadfiles.forEach((r) => {
                                    r.status = 'pause'
                                })
                            }
                        "
                    >
                        <LucideSquare class="size-1/2" />
                    </Button>
                    <!-- <Button class="aspect-square bg-blue-200 hover:bg-blue-300 text-blue-500 hover:text-white p-0">
                        <LucideSettings class="size-1/2" />
                    </Button> -->
                </div>
            </div>
        </div>
        <div class="col-span-4 flex flex-col bg-white/80 rounded-xl p-3 text-md gap-5">
            <div>文件列表</div>
            <div class="flex flex-col -mx-3 text-sm">
                <div class="grid grid-cols-[2fr_6rem_6rem] md:grid-cols-[2fr_6rem_6rem_4fr] gap-2 border-b border-black/20 pb-2 px-3">
                    <div>文件名</div>
                    <div>文件大小</div>
                    <div @click="handleShowSpeedInfo" class="flex flex-row gap-1 items-center">
                        上传速度
                        <LucideInfo class="size-3" />
                    </div>
                    <div class="hidden md:block">进度</div>
                </div>
                <div
                    :class="
                        cx(
                            'grid grid-cols-[2fr_6rem_6rem] md:grid-cols-[2fr_6rem_6rem_4fr] gap-2 py-2 border-b border-black/20 items-center hover:bg-primary/30 px-3 cursor-pointer',
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
                    <div class="flex flex-row gap-2 items-center grow md:grow-0 overflow-hidden">
                        <Button
                            class="size-8 p-0 hover:bg-white/50"
                            variant="ghost"
                            :disabled="['finish']?.includes(item?.procressType)"
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
                            <LucideArrowUpFromLine class="size-4 text-green-500" v-if="item?.status === 'start'" />
                            <LucideSquare class="size-4 text-gray-500" v-if="item?.status === 'pause'" />
                            <LucideCircleX class="size-4 text-red-500" v-if="item?.status === 'error'" />
                            <LucideCheckCircle class="size-4 text-green-500" v-if="item?.status === 'finish'" />
                        </Button>
                        <div class="truncate">{{ item?.file?.name }}</div>
                    </div>
                    <div>{{ getFileSize(item?.file?.size) }}</div>
                    <div>
                        {{
                            `${getFileSize(
                                (Object.entries(item?.uploadInfo?.chunks || {})?.filter(
                                    ([, chunk]) => chunk.status === 'success' && dayjs().unix() - 60 < chunk.createdAt
                                )?.length /
                                    60) *
                                    (item?.uploadInfo?.ChunkSize || 0)
                            )} /s`
                        }}
                    </div>
                    <div
                        class="flex flex-row gap-2 items-center col-span-3 md:col-span-1"
                        v-if="['hash', 'create', 'chunk']?.includes(item?.procressType)"
                    >
                        <LucideLoaderCircle class="size-4 animate-spin" />
                        <div>正在{{ item?.procressType }}中...</div>
                    </div>
                    <div class="flex flex-row gap-2 items-center col-span-3 md:col-span-1" v-if="item?.procressType === 'finish'">
                        {{ item?.status === 'finish' ? '云端已有相同hash文件, 秒传成功' : null }}
                        {{ item?.status === 'error' ? '上传失败，请稍后重试' : null }}
                    </div>
                    <div class="flex flex-row gap-2 items-center col-span-3 md:col-span-1" v-if="item?.procressType === 'upload'">
                        <div class="rounded-full bg-white/50 w-full h-2 overflow-hidden border border-white">
                            <div
                                class="h-full bg-primary"
                                :style="`width: ${clamp(
                                    ((Object.entries(item?.uploadInfo?.chunks || {})?.filter(([index, chunk]) => chunk.status === 'success')
                                        ?.length || 0) /
                                        (item?.uploadInfo?.chunkLength || 0)) *
                                        100,
                                    0,
                                    100
                                )}%`"
                            ></div>
                        </div>
                        <div>
                            {{
                                clamp(
                                    ((Object.entries(item?.uploadInfo?.chunks || {})?.filter(([index, chunk]) => chunk.status === 'success')
                                        ?.length || 0) /
                                        (item?.uploadInfo?.chunkLength || 0)) *
                                        100,
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
                <div>
                    区块： {{ selectedUploadfile?.uploadInfo?.chunkLength }} x {{ getFileSize(selectedUploadfile?.uploadInfo?.ChunkSize as number) }}
                </div>
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
