<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useQueryClient } from '@tanstack/vue-query'
import showDrawer from '~/lib/showDrawer'
import { toast } from 'vue-sonner'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'
import { LucideChevronDown, LucideDownload, LucideFileArchive } from '@lucide/vue'

const { t } = useI18n()
const props = defineProps<{
    data: any
}>()

const queryClient = useQueryClient()
const { downloadFile, getShareToken } = useMyAppShare()
const token = ref<string>()

const handleDownload = async (target: 'zip' | 'tar.gz') => {
    const { id } = props?.data || {}
    try {
        if (!token.value) {
            if (props?.data?.has_password) {
                token.value = await showDrawer({
                    render: ({ ...rest }) => h(PasswallShareDrawer, { ...rest, share_id: id }),
                })
            } else {
                token.value = await getShareToken(id)
            }
            if (!token.value) {
                throw new Error(t('page.shareView.fileShare.getTokenFailed'))
            }
        }
        downloadFile(token.value, undefined, target)
    } catch (error: any) {
        toast.error(error?.data?.message || error?.message || error)
    } finally {
        queryClient.invalidateQueries({ queryKey: ['share', id] })
    }
}
</script>

<template>
    <div class="flex w-full flex-col gap-3">
        <div class="flex w-full items-center justify-between gap-3">
            <h2 class="text-sm font-medium">文件列表</h2>
            <DropdownMenu v-if="(props?.data?.files?.length || 0) > 1">
                <DropdownMenuTrigger as-child>
                    <Button size="sm">
                        <LucideDownload class="size-4" />
                        {{ t('page.shareView.fileShare.packageDownload') }}
                        <LucideChevronDown class="size-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuItem @select="handleDownload('zip')">
                        <LucideFileArchive class="size-4" />
                        ZIP
                    </DropdownMenuItem>
                    <DropdownMenuItem @select="handleDownload('tar.gz')">
                        <LucideFileArchive class="size-4" />
                        TAR.GZ
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
        <ShareFileInfoList class="w-full" :files="props?.data?.files || []" />
    </div>
</template>
