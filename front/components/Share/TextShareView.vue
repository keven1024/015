<script setup lang="ts">
import AsyncButton from '@/components/ui/button/AsyncButton.vue'
import { toast } from 'vue-sonner'
import MarkdownRender from '@/components/MarkdownRender.vue'
import showDrawer from '~/lib/showDrawer'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'

const { t } = useI18n()
const props = defineProps<{
    data: any
}>()

const { getShareToken } = useMyAppShare()

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
    <div class="flex max-h-full w-full flex-col gap-3">
        <template v-if="!previewText">
            <div class="w-full">
                <AsyncButton @click="handlePreview" class="w-full">{{ t('page.shareView.textShare.viewBtn') }}</AsyncButton>
            </div>
        </template>
        <template v-else>
            <div class="flex justify-end">
                <CopyButton :value="previewText" />
            </div>
            <MarkdownRender :markdown="previewText" class="rounded-md bg-white/70 p-3 w-full max-w-full min-h-80 overflow-y-auto" />
        </template>
    </div>
</template>
