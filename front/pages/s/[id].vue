<script setup lang="ts">
import { LucideAlertCircle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { cx } from 'class-variance-authority'
import { times } from 'lodash-es'
import dayjs from 'dayjs'
import FileShareView from '@/components/Share/FileShareView.vue'
import TextShareView from '@/components/Share/TextShareView.vue'
const route = useRoute()
const router = useRouter()

const id = computed(() => route.params.id)

const { state, isLoading } = useAsyncState(async () => {
    const data = await $fetch<{
        code: number
        data: {
            id?: string
            expire_at?: number
        }
    }>(`/api/share/${id.value}`)
    return data?.data
}, null)

const isExpired = computed(() => {
    const { expire_at } = state.value || {}
    return !state || !expire_at || dayjs(expire_at * 10e2).isBefore(dayjs())
})

</script>

<template>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 my-5">
        <div v-if="isLoading" class="flex flex-col gap-5 items-center">
            <Skeleton :class="cx('w-40 h-5 rounded-full', i > 0 && '!w-20')" v-for="i in times(3)" :key="i" />
        </div>
        <template v-else>
            <div v-if="isExpired" class="flex flex-col gap-5 items-center">
                <LucideAlertCircle :size="48" class="text-orange-500 rounded-full bg-orange-500/30 p-2" />
                <div class="text-xl ">此链接已过期。</div>
                <Button @click="() => {
                    router.push('/')
                }">返回首页</Button>
            </div>
            <template v-else>
                <FileShareView :data="state" />
                <TextShareView :data="state" />
            </template>
        </template>
    </div>

</template>
