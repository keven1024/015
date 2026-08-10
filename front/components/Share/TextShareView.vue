<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import { cx } from 'class-variance-authority'
import { toast } from 'vue-sonner'
import MarkdownRender from '@/components/MarkdownRender.vue'
import showDrawer from '~/lib/showDrawer'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'

const { t } = useI18n()
const props = defineProps<{
    data: any
}>()

const { getShareToken } = useMyAppShare()

const textShareInfo = computed(() => {
    return [
        { label: t('page.shareView.textShare.needPassword'), type: 'bool' as const, value: props?.data?.has_password ?? false },
        {
            label: t('page.shareView.textShare.expireTime'),
            type: 'countdown' as const,
            value: props?.data?.expire_at ?? 0,
        },
        { label: t('page.shareView.textShare.remainingViews'), type: 'string' as const, value: props?.data?.download_nums ?? 0 },
    ]
})
const previewText = ref<string | null>(null)

const handlePreview = async () => {
    try {
        let token = null
        if (props?.data?.has_password) {
            token = await showDrawer({
                render: ({ ...rest }) => h(PasswallShareDrawer, { ...rest, share_id: props?.data?.id }),
            })
        } else {
            token = await getShareToken(props?.data?.id)
        }
        const r = await $fetch<{
            code: number
            data: {
                text: string
            }
        }>(`/api/download?token=${token}`)
        previewText.value = r?.data?.text
    } catch (error: any) {
        toast.error(error?.data?.message || error?.message || error)
    }
}
</script>
<template>
    <div :class="cx('flex flex-col max-h-full', !!previewText ? 'gap-3' : 'gap-16 items-center')">
        <div :class="cx('flex flex-row w-full', !!previewText ? 'justify-between' : 'justify-center')">
            <h1 class="text-xl">{{ t('page.shareView.textShare.title') }}</h1>
            <CopyButton v-if="!!previewText" :value="previewText as string" />
        </div>
        <template v-if="!previewText">
            <ShareInfoCards :items="textShareInfo" />
            <div class="w-full">
                <AsyncButton @click="handlePreview" class="w-full">{{ t('page.shareView.textShare.viewBtn') }}</AsyncButton>
            </div>
        </template>
        <template v-else>
            <MarkdownRender :markdown="previewText" class="rounded-md bg-white/70 p-3 w-full max-w-full min-h-80 overflow-y-auto" />
        </template>
    </div>
</template>
