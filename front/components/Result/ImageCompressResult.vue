<script setup lang="ts">
import { useQueries, useQuery } from '@tanstack/vue-query'
import { AsyncButton } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { filesize } from 'filesize'
import useMyAppShare from '@/composables/useMyAppShare'
import { toast } from 'vue-sonner'
import type { handleFileComponentProps } from './types'
import { get } from 'lodash-es'

const { t } = useI18n()
const emit = defineEmits<{
    (e: 'change', key: string): void
}>()
const props = defineProps<handleFileComponentProps>()
const fileIds = computed(() => props?.data?.files?.map((item) => item.id))
const { data: taskIds } = useQuery({
    queryKey: ['create-image-compress', fileIds.value],
    queryFn: async () => {
        return await Promise.all(
            props?.data?.files?.map(async (file) => {
                const { id } = file || {}
                if (!id) return
                const data = await $fetch<{
                    code: number
                    data: {
                        id?: string
                    }
                }>(`/api/task/image:compress`, {
                    method: 'POST',
                    body: {
                        file_id: id,
                    },
                })
                return data?.data?.id
            })
        )
    },
    staleTime: Infinity,
    enabled: !!fileIds.value && fileIds.value?.length > 0,
})

const taskResults = useQueries({
    queries: computed(
        () =>
            taskIds?.value?.filter(Boolean).map((taskId) => {
                return {
                    queryKey: ['task-image-compress', taskId],
                    queryFn: async () => {
                        const data = await $fetch<{
                            code: number
                            data: {
                                result: {
                                    old_file: {
                                        id: string
                                        size: number
                                    }
                                    new_file: {
                                        id: string
                                        size: number
                                    }
                                }[]
                                status: 'success' | 'retry' | 'archived'
                                err?: {
                                    message?: string
                                    retry?: number
                                    max_retry?: number
                                }
                            }
                        }>(`/api/task/${taskId}`)
                        return data?.data
                    },
                    enabled: !!taskId,
                }
            }) ?? []
    ),
})

const totalSize = computed(() => {
    return taskResults.value.reduce(
        (acc, item) => {
            const { new_file, old_file } = get(item, 'data.result.0') || {}
            return {
                oldSize: acc.oldSize + (old_file?.size ?? 0),
                newSize: acc.newSize + (new_file?.size ?? 0),
            }
        },
        { oldSize: 0, newSize: 0 }
    )
})

const { downloadFileByShareId, createFileShare } = useMyAppShare()

const { counter, pause } = useInterval(2000, { controls: true })

watch(
    () => counter.value,
    () => {
        taskResults.value.forEach((item) => {
            if (['success', 'archived'].includes(item.data?.status ?? '')) return
            item.refetch()
        })

        if (taskResults.value.every((item) => ['success', 'archived'].includes(item.data?.status ?? ''))) {
            pause()
        }
    }
)
</script>
<template>
    <BaseCard class="flex flex-col gap-3" :title="t('page.result.imageCompress.title')" :showBackButton="true">
        <div class="flex flex-row gap-3">
            <div class="rounded-xl flex flex-col bg-white/70 px-3 py-2 gap-1 basis-2/3">
                <div class="text-sm font-semibold">{{ t('page.result.imageCompress.totalSize') }}</div>
                <div class="text-2xl font-light flex flex-row items-center gap-1">
                    <span v-if="totalSize.oldSize > 0" class="opacity-75">{{ filesize(totalSize.oldSize) }}</span>
                    <span v-else><Skeleton class="w-12 h-10" /></span>
                    <LucideChevronsRight class="size-6" />
                    <span v-if="totalSize.newSize > 0">{{ filesize(totalSize.newSize) }}</span>
                    <span v-else><Skeleton class="w-12 h-10" /></span>
                    <div
                        v-if="totalSize.oldSize > 0 && totalSize.newSize > 0"
                        class="rounded flex flex-row items-center bg-green-100 text-green-600 p-1 py-0.5 text-sm"
                    >
                        <LucideArrowDown class="size-4" />
                        {{ ((1 - totalSize.newSize / totalSize.oldSize) * 100).toFixed(2) }}%
                    </div>
                </div>
            </div>
            <div class="rounded-xl flex flex-col bg-white/70 px-3 py-2 gap-1 basis-1/3">
                <div class="text-sm font-semibold">{{ t('page.result.imageCompress.task') }}</div>
                <div class="text-3xl font-light">{{ taskResults.length }}</div>
            </div>
        </div>
        <div v-for="(item, index) in props?.data?.files" class="flex flex-row rounded-xl bg-white/70 p-3 justify-between w-full gap-3">
            <div class="flex flex-row gap-2 items-center w-full overflow-hidden">
                <div class="*:h-12 overflow-hidden">
                    <FileIcon :file="item?.file" />
                </div>
                <div class="flex flex-col gap-0.5 flex-1 overflow-hidden">
                    <div class="truncate w-auto">{{ item?.file?.name }}</div>
                    <div class="text-xs opacity-50">{{ filesize(item?.file?.size ?? 0) }}</div>
                </div>
            </div>
            <div class="flex items-center justify-center" v-if="!taskResults?.[index]?.data">
                <Skeleton class="w-16 h-12" />
            </div>
            <div class="flex flex-row gap-1 items-center text-sm shrink-0" v-if="taskResults?.[index]?.data?.status === 'retry'">
                <LucideLoader2 class="size-4 animate-spin" />
                {{
                    t('page.result.imageCompress.retry', [
                        taskResults?.[index]?.data?.err?.retry ?? 0,
                        taskResults?.[index]?.data?.err?.max_retry ?? 0,
                    ])
                }}
            </div>
            <div class="flex items-center justify-center shrink-0" v-if="taskResults?.[index]?.data?.status === 'archived'">
                <div class="text-sm text-red-500 px-2 py-1 rounded-md bg-red-100">{{ t('page.result.imageCompress.failed') }}</div>
            </div>
            <div class="flex flex-row gap-2 items-center shrink-0" v-if="taskResults?.[index]?.data?.status === 'success'">
                <div class="flex flex-col gap-1 items-center">
                    <div class="rounded flex flex-row items-center bg-green-100 text-green-600 px-1 text-xs">
                        <LucideArrowDown class="size-4" />
                        {{
                            (
                                (1 -
                                    (taskResults?.[index]?.data?.result?.[0]?.new_file?.size ?? 0) /
                                        (taskResults?.[index]?.data?.result?.[0]?.old_file?.size ?? 0)) *
                                100
                            ).toFixed(2)
                        }}%
                    </div>
                    <div class="text-xs opacity-50">{{ filesize(taskResults?.[index]?.data?.result?.[0]?.new_file?.size ?? 0) }}</div>
                </div>
                <AsyncButton
                    size="icon"
                    @click="
                        async () => {
                            const { new_file } = taskResults?.[index]?.data?.result?.[0] || {}
                            if (!new_file?.id) return
                            const data = await createFileShare({
                                files: [{ id: new_file?.id as string, file_name: item?.file?.name }],
                                config: {
                                    download_nums: 1,
                                    expire_time: 60,
                                    has_pickup_code: false,
                                    has_password: false,
                                },
                            })
                            const { id } = data?.data || {}
                            if (!id) {
                                return
                            }
                            try {
                                await downloadFileByShareId(id)
                            } catch (error: any) {
                                toast.error(error?.data?.message || error?.message || error)
                            }
                        }
                    "
                    ><LucideArrowDown
                /></AsyncButton>
            </div>
        </div>
    </BaseCard>
</template>
