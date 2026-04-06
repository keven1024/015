<script setup lang="ts">
import { LucideAlertCircle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import dayjs from 'dayjs'
import FileShareView from '@/components/Share/FileShareView.vue'
import TextShareView from '@/components/Share/TextShareView.vue'
import { useQuery } from '@tanstack/vue-query'
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const id = computed(() => route.params.id)

const { data, isLoading, error } = useQuery({
    queryKey: ['share', id.value],
    queryFn: async () => {
        const data = await $fetch<{
            code: number
            data: {
                id?: string
                expire_at?: number
                type?: string
            }
        }>(`/api/share/${id.value}`)
        return data?.data
    },
    retry: false,
})

const isExpired = computed(() => {
    const { expire_at } = data.value || {}
    return !data || !expire_at || dayjs(expire_at * 10e2).isBefore(dayjs())
})

const componentMap = {
    file: FileShareView,
    text: TextShareView,
}
</script>

<template>
    <BaseCard class="my-5 overflow-hidden">
        <div v-if="isLoading" class="flex flex-col gap-5 items-center">
            <Skeleton class="w-20 h-5 rounded-full" />
            <Skeleton class="w-16 h-16 rounded-xl" />
            <Skeleton class="w-20 h-5 rounded-full" />
            <div class="flex flex-row w-full justify-around">
                <Skeleton class="size-10 rounded-xl" v-for="i in 3" />
            </div>
            <Skeleton class="w-full h-5 rounded-full" />
        </div>
        <template v-else>
            <div v-if="isExpired || !data" class="flex flex-col gap-5 items-center">
                <LucideAlertCircle :size="48" class="text-orange-500 rounded-full bg-orange-500/30 p-2" />
                <div class="text-xl">{{ t('page.shareView.linkExpired') }}</div>
                <Button
                    @click="
                        () => {
                            router.push('/')
                        }
                    "
                    >{{ t('btn.backToHome') }}</Button
                >
            </div>
            <template v-else>
                <component :is="componentMap[data?.type as keyof typeof componentMap] || 'div'" :data="data" />
            </template>
        </template>
    </BaseCard>
</template>
