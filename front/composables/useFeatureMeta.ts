import { LucideShare, LucideImageMinus, LucideArrowRightLeft } from 'lucide-vue-next'
import useMyAppConfig from '@/composables/useMyAppConfig'
import type { FileHandleKey, TextHandleKey } from '../components/Preprocessing/types'
import generateRandomColors from '@/lib/generateRandomColors'

export type FeatureKey = FileHandleKey | TextHandleKey

export type FeatureMeta = {
    key: FeatureKey
    label: string
    icon: any
}

const allFeatureMeta = (t: (key: string) => string): FeatureMeta[] => [
    {
        key: 'file-share',
        label: t('page.upload.file.handleType.file-share'),
        icon: LucideShare,
    },
    {
        key: 'file-image-compress',
        label: t('page.upload.file.handleType.file-image-compress'),
        icon: LucideImageMinus,
    },
    {
        key: 'file-image-convert',
        label: t('page.upload.file.handleType.file-image-convert'),
        icon: LucideArrowRightLeft,
    },
    {
        key: 'text-share',
        label: t('page.upload.text.handleType.text-share'),
        icon: LucideShare,
    },
]

export function useFeatureMeta() {
    const { t } = useI18n()
    const appConfig = useMyAppConfig()

    return computed(() => {
        const enabledKeys = appConfig.value?.features ?? []
        const result = allFeatureMeta(t).filter((meta) => enabledKeys.includes(meta.key))
        const colors = generateRandomColors(result.length)
        return result.map((meta, index) => ({ ...meta, style: { backgroundColor: colors[index] } }))
    })
}
