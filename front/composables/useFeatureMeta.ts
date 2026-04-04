import { LucideShare, LucideImageMinus, LucideArrowRightLeft } from 'lucide-vue-next'
import useMyAppConfig from '@/composables/useMyAppConfig'
import type { FileHandleKey, TextHandleKey } from '../components/Preprocessing/types'

export type FeatureKey = FileHandleKey | TextHandleKey

export type FeatureMeta = {
    key: FeatureKey
    label: string
    icon: any
    className: string
}

const allFeatureMeta = (t: (key: string) => string): FeatureMeta[] => [
    {
        key: 'file-share',
        label: t('page.upload.file.handleType.file-share'),
        icon: LucideShare,
        className: 'bg-green-300',
    },
    {
        key: 'file-image-compress',
        label: t('page.upload.file.handleType.file-image-compress'),
        icon: LucideImageMinus,
        className: 'bg-red-300',
    },
    {
        key: 'file-image-convert',
        label: t('page.upload.file.handleType.file-image-convert'),
        icon: LucideArrowRightLeft,
        className: 'bg-purple-300',
    },
    {
        key: 'text-share',
        label: t('page.upload.text.handleType.text-share'),
        icon: LucideShare,
        className: 'bg-green-300',
    },
]

export function useFeatureMeta() {
    const { t } = useI18n()
    const appConfig = useMyAppConfig()

    return computed(() => {
        const enabledKeys = appConfig.value?.features ?? []
        return allFeatureMeta(t).filter((meta) => enabledKeys.includes(meta.key))
    })
}
