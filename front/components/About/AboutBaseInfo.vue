<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Skeleton } from '@/components/ui/skeleton'
import getFileSize from '~/lib/getFileSize'
import SparkMD5 from 'spark-md5'
import useMyAppConfig from '@/composables/useMyAppConfig'
import Progress from '~/components/ui/progress/Progress.vue'
import renderI18n from '~/lib/renderI18n'
import { I18nT } from 'vue-i18n'

const appConfig = useMyAppConfig()
const { data, isLoading } = useQuery({
    queryKey: ['about'],
    queryFn: async () => {
        const data = await $fetch<{
            data: {
                file: {
                    maximun: number
                    current: number
                }
                bg_url?: string
                avatar?: string
                name?: string
                email?: string
                url?: string
            }
        }>('/api/about')
        return data?.data
    },
})
const { t } = useI18n()

const genUserAvatar = (email: string) => {
    return `https://www.gravatar.com/avatar/${SparkMD5.hash(email)}?d=retro`
}
</script>

<template>
    <template v-if="isLoading">
        <div class="flex flex-col gap-2">
            <Skeleton class="aspect-[3/1] w-full rounded-xl" />
            <div class="flex flex-col gap-2 items-center">
                <Skeleton class="h-6 w-32 rounded" />
                <Skeleton class="h-4 w-52 rounded" />
            </div>
        </div>
    </template>
    <template v-else>
        <NuxtImg v-if="data?.bg_url" :src="data?.bg_url" class="aspect-[3/1] w-full rounded-xl" fit="cover" />
        <div class="flex flex-col gap-2 items-center">
            <div class="text-xl">{{ renderI18n(appConfig?.site_title ?? {}, 'en') }}</div>
            <div class="text-sm opacity-75 text-center px-5">
                <I18nT keypath="about.powerBy" tag="span">
                    <NuxtLink href="https://github.com/keven1024/015" target="_blank" class="text-primary hover:underline">015</NuxtLink>
                </I18nT>
            </div>
        </div>
    </template>

    <div class="font-semibold">{{ t('about.systemInfo') }}</div>
    <template v-if="isLoading">
        <div class="flex flex-row gap-2">
            <Skeleton class="w-full h-20 rounded-xl" v-for="i in 2" :key="i" />
        </div>
    </template>
    <template v-else>
        <div class="grid grid-cols-2 gap-2">
            <div class="rounded-xl bg-white/50 flex-1 flex flex-col p-3 gap-2">
                <div class="opacity-75 text-xs">{{ t('about.admin') }}</div>
                <div
                    class="flex flex-row gap-2 items-center cursor-pointer"
                    @click="
                        () => {
                            if (data?.url) {
                                navigateTo(data?.url, { external: true })
                                return
                            }
                            if (data?.email) {
                                navigateTo(`mailto:${data?.email ?? ''}`, { external: true })
                                return
                            }
                            return
                        }
                    "
                >
                    <Avatar class="size-10">
                        <AvatarImage v-if="!!data?.avatar || !!data?.email" :src="data?.avatar || genUserAvatar(data?.email as string)" />
                        <AvatarFallback class="bg-black/10 font-bold">
                            {{ data?.name?.charAt(0)?.toUpperCase() }}
                        </AvatarFallback>
                    </Avatar>
                    <div class="flex flex-col">
                        <div class="text-md font-semibold">{{ data?.name }}</div>
                        <div class="text-xs opacity-75">{{ data?.email }}</div>
                    </div>
                </div>
            </div>
            <div class="rounded-xl bg-white/50 flex-1 flex flex-col p-3 gap-2">
                <div class="opacity-75 text-xs">{{ t('about.storage') }}</div>
                <div class="text-right flex flex-row items-baseline">
                    <span class="text-lg font-semibold">{{ getFileSize(data?.file?.current ?? 0) }}</span>
                    <span class="text-md opacity-75">/ {{ getFileSize(data?.file?.maximun ?? 0) }}</span>
                </div>
                <Progress class="h-1" :model-value="((data?.file?.current ?? 0) / (data?.file?.maximun ?? 0)) * 100" />
            </div>
        </div>
    </template>
</template>
