<script setup lang="ts">
import { LucideSquare, LucideInfo, LucideFolders, LucideArrowUpFromLine, LucideCircleX, LucideCheckCircle, LucideLoaderCircle } from '@lucide/vue'
import Button from '@/components/ui/button/Button.vue'

import getFileSize from '~/lib/getFileSize'
import { cx } from 'class-variance-authority'
import asyncWorker from '@/lib/asyncWorker'
import calcFileHashWorker from '@/lib/calcFileHashWorker?worker'
import { detectSupportedEngines } from '@/lib/calcFileHash'
import { clamp, get, isEmpty, isNumber, sample, shuffle, throttle, times } from 'lodash-es'
import { nanoid } from 'nanoid'
import { toast } from 'vue-sonner'
import dayjs from 'dayjs'
import showDrawer from '~/lib/showDrawer'
import FileUploadTotalSpeedView from './FileUploadTotalSpeedView.vue'
import FileUploadTotalProgressControlView from './FileUploadTotalProgressControlView.vue'
import FileUploadSpeedInfoView from './FileUploadSpeedInfoView.vue'
import getFileChunk from '~/lib/getFileChunk'
import type { FileHandleKey } from '~/components/Preprocessing/types'
import asyncWait from '@/lib/asyncWait'
import axios from 'axios'
import type { SelectedFile, Uploadfile } from './types'
import FileUploadDetailView from './FileUploadDetailView.vue'

const props = defineProps<{
    data: { file: File[]; config: Record<string, any>; handle_type: FileHandleKey }
}>()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const form = useFormContext()
const { t } = useI18n()

const selectedFile = ref<SelectedFile>(null)
const uploadfiles = ref<Uploadfile[]>([])

const procressTaskList = ref<Map<string, any>>(new Map())
const activeTaskList = computed(() => uploadfiles.value.filter((r) => r.queue.length > 0 && r.status === 'start'))
const activeTaskAllQueue = computed(() => activeTaskList.value.flatMap((r) => r.queue))
const batchNum = ref(3)

const speedChartData = ref<Record<number, { fileId: string; index: number; value: number }[]>>({})

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
            await asyncWait(10) // 让出主线程，防止空转导致ui卡死
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
                            toast.error(t('page.progress.file.uploadError'), {
                                description: t('page.progress.file.chunkUploadFailed', [file?.file?.name, index]),
                            })
                            uploadfiles.value[uploadFileIndex]!.status = 'error'
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
                            toast.warning(t('page.progress.file.uploadError'), {
                                description: t('page.progress.file.chunkUploadRetry', [file?.file?.name, index]),
                            })
                        }
                        if (queueType === 'sync') {
                            uploadfiles.value[uploadFileIndex]!.status = 'error'
                            toast.error(t('page.progress.file.uploadError'), {
                                description: t('page.progress.file.fileUploadFailed', [file?.file?.name]),
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

const LARGE_FILE_THRESHOLD = 500 * 1024 * 1024 // 500 MB

const handleHash = async (fileId: string) => {
    const uploadfile = uploadfiles.value.find((item) => item.fileId === fileId)
    if (!uploadfile?.file) return
    uploadfile.procressType = 'hash'
    const supportedEngines = detectSupportedEngines()
    if (supportedEngines.length === 0) {
        throw new Error(t('page.progress.file.hashEngineNotFound'))
    }
    const preferredEngine = uploadfile.file.size >= LARGE_FILE_THRESHOLD ? 'wasm' : 'native'
    const res = await asyncWorker(calcFileHashWorker, {
        data: { file: uploadfile.file, engine: supportedEngines.includes(preferredEngine) ? preferredEngine : supportedEngines?.[0] },
    })
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
            chunks: number[]
        }
    }>('/api/file/create', {
        method: 'POST',
        body: {
            size,
            mime_type: type || 'application/octet-stream',
            hash,
        },
    })
    const { id, chunk_size, type: createType, chunks = [] } = createData?.data || {}
    uploadfile.id = id
    uploadfile.uploadInfo = {
        chunks: Object.fromEntries(chunks.map((index: number) => [index - 1, { status: 'success', createdAt: dayjs().unix() }])),
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
    const { data: res } = await axios.post('/api/file/slice', formData, {
        onUploadProgress: throttle((progressEvent) => {
            const { rate } = progressEvent || {}
            const timestamp = dayjs().unix()
            speedChartData.value[timestamp] = [...(speedChartData.value[timestamp] || []), { fileId, index, value: rate || 0 }]
        }, 1000),
    })
    const { code } = res || {}
    if (code !== 200) {
        throw new Error(t('page.progress.file.uploadFailed'))
    }
    uploadfile.uploadInfo!.chunks![index]!.status = 'success'
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
        throw new Error(t('page.progress.file.uploadFailed'))
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
    <BaseCard class="grid grid-cols-4 gap-5">
        <FileUploadTotalSpeedView :speedChartData="speedChartData" />
        <FileUploadTotalProgressControlView :uploadfiles="uploadfiles" />
        <div class="col-span-4 flex flex-col bg-white/80 rounded-xl p-3 text-md gap-5">
            <div class="flex flex-row gap-2 items-center">
                <LucideFolders class="size-4" />
                {{ t('page.progress.file.fileList') }}
            </div>
            <div class="flex flex-col -mx-3 text-sm">
                <div class="grid grid-cols-[2fr_6rem_6rem] md:grid-cols-[2fr_6rem_6rem_4fr] gap-2 border-b border-black/20 pb-2 px-3">
                    <div>{{ t('page.progress.file.fileName') }}</div>
                    <div>{{ t('page.progress.file.fileSize') }}</div>
                    <div @click="handleShowSpeedInfo" class="flex flex-row gap-1 items-center">
                        {{ t('page.progress.file.uploadSpeed') }}
                        <LucideInfo class="size-3" />
                    </div>
                    <div class="hidden md:block">{{ t('page.progress.file.progress') }}</div>
                </div>
                <div
                    :class="
                        cx(
                            'grid grid-cols-[2fr_6rem_6rem] md:grid-cols-[2fr_6rem_6rem_4fr] gap-2 py-2 border-b border-black/20 items-center hover:bg-primary/30 px-3 cursor-pointer',
                            selectedFile === item?.fileId && 'bg-primary text-white hover:bg-primary!'
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
                                        uploadfiles[index]!.status = 'pause'
                                    } else {
                                        uploadfiles[index]!.status = 'start'
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
                        <template v-if="item?.status === 'start' && item?.procressType === 'upload'">
                            {{
                                `${getFileSize(
                                    Object.entries(speedChartData)
                                        ?.filter(([key, value]) => value.some((r) => r.fileId === item?.fileId) && dayjs().unix() - 1 === Number(key))
                                        ?.reduce((acc, curr) => acc + curr[1]?.reduce((acc, curr) => acc + curr.value, 0), 0) ?? 0
                                )} /s`
                            }}
                        </template>
                    </div>
                    <div
                        class="flex flex-row gap-2 items-center col-span-3 md:col-span-1"
                        v-if="['hash', 'create', 'chunk']?.includes(item?.procressType)"
                    >
                        <LucideLoaderCircle class="size-4 animate-spin" />
                        <div>{{ t(`page.progress.file.processing.${item?.procressType}`) }}</div>
                    </div>
                    <div class="flex flex-row gap-2 items-center col-span-3 md:col-span-1" v-if="item?.procressType === 'finish'">
                        {{ item?.status === 'finish' ? t('page.progress.file.instantUploadSuccess') : null }}
                        {{ item?.status === 'error' ? t('page.progress.file.uploadFailedRetry') : null }}
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
        <FileUploadDetailView :uploadfiles="uploadfiles" :selectedFile="selectedFile" />
    </BaseCard>
</template>
