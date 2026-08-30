<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useQueryClient } from '@tanstack/vue-query'
import showDrawer from '~/lib/showDrawer'
import { toast } from 'vue-sonner'
import PasswallShareDrawer from '~/components/Drawer/PasswallShareDrawer.vue'
import { LucideChevronDown, LucideDownload, LucideFileArchive, LucideListChecks } from '@lucide/vue'
import type { DownloadArchiveTarget } from '@/composables/useMyAppShare'

const { t } = useI18n()
const props = defineProps<{
    data: any
}>()

const queryClient = useQueryClient()
const { downloadFile, getShareToken } = useMyAppShare()
const token = ref<string>()
const selectedFiles = ref<any[]>([])
const files = computed<any[]>(() => props?.data?.files || [])
const selectedFileIds = computed(() => selectedFiles.value.map((file) => file.id))

const handleDownload = async (target?: DownloadArchiveTarget, fileIds?: string[]) => {
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
        downloadFile(token.value, fileIds, target)
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
            <h2 class="text-sm font-medium">{{ t('page.shareView.fileShare.fileList') }}</h2>
            <div class="flex items-center gap-2">
                <Button
                    class="size-8"
                    type="button"
                    variant="outline"
                    size="icon"
                    :disabled="files.length < 2"
                    @click="() => (selectedFiles = selectedFiles.length === files.length ? [] : [...files])"
                >
                    <LucideListChecks class="size-4" />
                </Button>
                <DropdownMenu v-if="(props?.data?.files?.length || 0) > 1">
                    <DropdownMenuTrigger as-child>
                        <Button size="sm" :disabled="selectedFiles.length < 2">
                            <LucideDownload class="size-4" />
                            {{ t('page.shareView.fileShare.packageDownload') }}
                            <LucideChevronDown class="size-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuItem @select="handleDownload('zip', selectedFileIds)">
                            <LucideFileArchive class="size-4" />
                            ZIP
                        </DropdownMenuItem>
                        <DropdownMenuItem @select="handleDownload('tar.gz', selectedFileIds)">
                            <LucideFileArchive class="size-4" />
                            TAR.GZ
                        </DropdownMenuItem>
                        <DropdownMenuItem @select="handleDownload('tar.zst', selectedFileIds)">
                            <LucideFileArchive class="size-4" />
                            TAR.ZST
                        </DropdownMenuItem>
                        <DropdownMenuItem @select="handleDownload('tar.s2', selectedFileIds)">
                            <LucideFileArchive class="size-4" />
                            TAR.S2
                        </DropdownMenuItem>
                        <DropdownMenuItem @select="handleDownload('tar.snappy', selectedFileIds)">
                            <LucideFileArchive class="size-4" />
                            TAR.SNAPPY
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
        <ShareFileInfoList v-model="selectedFiles" class="w-full" :files="files">
            <template #default="{ file }">
                <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    :title="t('page.about.download')"
                    :aria-label="t('page.about.download')"
                    @click.stop="handleDownload(undefined, [file.id])"
                >
                    <LucideDownload class="size-4" />
                </Button>
            </template>
        </ShareFileInfoList>
    </div>
</template>
