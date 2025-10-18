<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useClipboard } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { useQuery } from '@tanstack/vue-query'
import useMyAppShare from '@/composables/useMyAppShare'
import useMyAppConfig from '@/composables/useMyAppConfig'
import showDrawer from '@/lib/showDrawer'
import QrCoreDrawer from '@/components/Drawer/QrCoreDrawer.vue'
import dayjs from 'dayjs'
import { h } from 'vue'
import type { TextHandleKey } from '../Preprocessing/types'

const props = defineProps<{
    data: { text: string; config: Record<string, any>; handle_type: TextHandleKey }
}>()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

const { createTextShare } = useMyAppShare()
const { data } = useQuery({
    queryKey: ['create-share', props?.data?.text],
    staleTime: Infinity,
    queryFn: async () => {
        const { config, text } = props?.data || {}
        const data = await createTextShare({
            text,
            config,
        })
        return data?.data
    },
})
const appConfig = useMyAppConfig()
const url = computed(() => {
    const { id } = data?.value || {}
    return `${appConfig?.value?.site_url}/s/${id}`
})

const { copy } = useClipboard()
const { t } = useI18n()
</script>

<template>
    <div class="flex flex-col gap-3">
        <div class="flex flex-row gap-2">
            <div class="flex flex-row justify-between w-full">
                <h2 class="text-lg">{{ t('textshareresult.title') }}</h2>
                <Button
                    variant="outline"
                    class="bg-white/70"
                    size="icon"
                    @click="
                        () => {
                            emit('change', 'input')
                        }
                    "
                >
                    <LucideHome />
                </Button>
            </div>
        </div>
        <div class="flex flex-col md:flex-row gap-5 rounded-md p-5 bg-white/20 backdrop-blur-xl w-full">
            <div class="flex flex-col gap-2 flex-1">
                <div class="text-sm font-semibold">{{ t('textshareresult.info') }}</div>
                <div class="grid grid-cols-2 gap-2">
                    <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
                        <div class="text-xs font-semibold">{{ t('textshareresult.viewNums') }}</div>
                        <div class="text-3xl font-light">{{ data?.download_nums }}</div>
                    </div>
                    <div class="rounded-xl flex flex-col bg-black/5 px-3 py-2 gap-1">
                        <div class="text-xs font-semibold">{{ t('textshareresult.expireTime') }}</div>
                        <div class="text-md font-light">
                            {{ dayjs((data?.expire_at ?? 0) * 1000).format('YYYY-MM-DD HH:mm:ss') }}
                        </div>
                    </div>
                    <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1" v-if="data?.pickup_code">
                        <div class="flex flex-row justify-between w-full items-center">
                            <div class="text-xs font-semibold">{{ t('textshareresult.pickupCode') }}</div>
                            <Button
                                variant="outline"
                                class="bg-white/70 p-0 size-6"
                                size="icon"
                                @click="
                                    () => {
                                        copy(data?.pickup_code as string)
                                        toast.success(t('textshareresult.copySuccess'))
                                    }
                                "
                            >
                                <LucideCopy class="size-3" />
                            </Button>
                        </div>
                        <div class="flex flex-row gap-2">
                            <div v-for="s in data?.pickup_code" class="text-2xl font-light">
                                {{ s }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="flex flex-col gap-5 flex-1">
                <div class="text-sm font-semibold">{{ t('textshareresult.link') }}</div>
                <div class="flex flex-row gap-2">
                    <Input v-model="url" class="bg-white/70" readonly />
                    <Button
                        variant="outline"
                        class="bg-white/70"
                        size="icon"
                        @click="
                            () => {
                                copy(url)
                                toast.success(t('textshareresult.copySuccess'))
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
                                            data: url,
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
        <h2 class="text-md">{{ t('textshareresult.content') }}</h2>
        <MarkdownRender class="prose rounded-md bg-white/70 p-3 w-full max-w-full min-h-[30vh]" :markdown="props?.data?.text" />
    </div>
</template>
