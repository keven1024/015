<script lang="ts" setup>
import CircularProgress from '@/components/CircularProgress.vue'
import { chunk, get, shuffle, times } from 'lodash-es'
import { cx } from 'class-variance-authority'
import calcFileHash from '@/lib/calcFileHash'
import { filesize } from 'filesize'
import { toast } from 'vue-sonner'
const props = defineProps<{
    data: { file: File; config: any; handle_type: string }
}>()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

const form = useFormContext()

const step = ref<'hash' | 'upload'>('hash')
const calcHashTime = ref(0)
const chunkSize = ref(0)
const fileSliceUploadStatusList = ref<
    {
        status: string
        index: number
    }[]
>([])

const successCount = computed(() => fileSliceUploadStatusList.value.filter((item) => item.status === 'success').length)
const alreadyUploadSize = computed(() => successCount.value * chunkSize.value)
const uploadProgress = computed(() => Math.round((alreadyUploadSize.value / (props?.data?.file?.size || 0)) * 100))

const { error } = useAsyncState(async () => {
    const { file } = props.data || {}
    if (!file) return
    const { size, type = 'application/octet-stream' } = file || {}
    const now = Date.now()
    const hash = await calcFileHash({ file })
    if (hash) {
        step.value = 'upload'
        calcHashTime.value = Date.now() - now
    }

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
            mime_type: get(type, '', 'application/octet-stream'),
            hash,
        },
    })
    const { id, chunk_size, type: createType } = createData?.data || {}
    if (createType !== 'init') {
        // 文件存在
        form.setFieldValue('file_id', id)
        emit('change', 'result')
        return
    }
    chunkSize.value = chunk_size
    const chunks = Math.ceil(size / chunk_size)
    fileSliceUploadStatusList.value = times(chunks, (i) => ({
        status: 'pending',
        index: i,
    }))

    const readChunk = (start: number): Promise<ArrayBuffer> => {
        const fileReader = new FileReader()
        return new Promise((resolve, reject) => {
            const chunk = file.slice(start, start + chunk_size)
            fileReader.onload = (e) => resolve(e.target?.result as ArrayBuffer)
            fileReader.onerror = reject
            fileReader.readAsArrayBuffer(chunk)
        })
    }

    const chunkedUploadTasks = chunk(shuffle([...fileSliceUploadStatusList.value]), 3)
    for (let i = 0; i < chunkedUploadTasks?.length; i++) {
        await Promise.all(
            chunkedUploadTasks?.[i]?.map(async (item: any) => {
                const { index } = item || {}
                try {
                    const chunk = await readChunk(index * chunk_size)
                    // console.log('chunk', chunk)
                    const formData = new FormData()
                    formData.append('file', new Blob([chunk]))
                    formData.append('index', index + 1)
                    formData.append('id', id)
                    fileSliceUploadStatusList.value[index].status = 'uploading'
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
                    fileSliceUploadStatusList.value[index].status = 'success'
                } catch (error) {
                    console.log('error', error)
                    // fileSliceStatusList.value[index].status = 'error'
                }
            })
        )
    }
    const r = await $fetch<{
        code: number
    }>('/api/file/finish', {
        method: 'POST',
        body: {
            id,
        },
    })
    if (r?.code !== 200) {
        throw new Error('上传失败')
    }
    form.setFieldValue('file_id', id)
    emit('change', 'result')
}, null)

watch(error, (newVal) => {
    if (newVal) {
        toast.error('上传失败')
    }
    emit('change', 'input')
})
</script>
<template>
    <div class="flex flex-col gap-3">
        <div class="text-xl font-normal">正在上传</div>
        <div class="flex flex-col items-center gap-4 md:flex-row md:justify-evenly">
            <div :class="cx('flex flex-row items-center gap-5', step !== 'hash' && 'opacity-50')">
                <div class="flex flex-col gap-0.5 items-center">
                    <div class="text-xs opacity-50">1.计算hash</div>
                    <div class="text-3xl font-light">
                        <LucideLoaderCircle v-if="step === 'hash'" class="size-5 my-1 animate-spin" />
                        <div v-else>{{ calcHashTime }}ms</div>
                    </div>
                </div>
            </div>
            <div :class="cx('flex flex-row items-center gap-5 min-w-32', step !== 'upload' && 'opacity-50')">
                <CircularProgress :size="80" :value="step !== 'upload' ? 0 : uploadProgress" />
                <div class="flex flex-col gap-0.5 items-center">
                    <div class="text-xs opacity-50">2.上传文件</div>
                    <div class="text-3xl font-light">{{ step !== 'upload' ? 0 : uploadProgress }}%</div>
                    <div class="text-sm opacity-50" v-if="alreadyUploadSize">
                        {{ filesize(alreadyUploadSize) }} / {{ filesize(data?.file?.size) }}
                    </div>
                </div>
            </div>
        </div>
        <div class="flex flex-row gap-2 items-baseline">
            <div class="text-md font-normal">详细信息</div>
            <div class="text-xs opacity-50" v-if="step === 'upload' && fileSliceUploadStatusList?.length > 0">
                当前正在上传分块{{ `${successCount}/${Math.ceil(data?.file?.size / chunkSize)}` }} 并发:3
            </div>
        </div>
        <div class="flex flex-row flex-wrap gap-1">
            <div
                v-for="i in fileSliceUploadStatusList"
                :class="
                    cx(
                        'rounded size-4 ',
                        i.status === 'pending' && 'bg-white/90',
                        i.status === 'uploading' && 'bg-yellow-500',
                        i.status === 'success' && 'bg-green-500'
                    )
                "
            />
        </div>
    </div>
</template>
