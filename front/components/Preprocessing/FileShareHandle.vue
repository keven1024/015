<script setup lang="ts">
import InputField from '../Field/InputField.vue'
import SelectField from '../Field/SelectField.vue'
import SwitchField from '../Field/SwitchField.vue'
import FormButton from '../Field/FormButton.vue'
import NotifyConfigField from './NotifyConfigField.vue'
import type { FileShareHandleProps } from './types'
const { t } = useI18n()
const props = defineProps<{
    hide: () => void
    file: File[]
    onFileHandle: (props: FileShareHandleProps) => void
}>()
</script>

<template>
    <VeeForm v-slot="{ values, setFieldValue }" :initialValues="{ download_nums: 1, expire_time: 1440 }">
        <div class="flex flex-col gap-3 max-h-[75vh]">
            <h2 class="text-lg font-bold">{{ t('page.shareOptions.file.title') }}</h2>
            <div class="flex flex-col gap-3 flex-1 overflow-y-auto">
                <div class="flex flex-row items-center gap-2 text-sm">
                    <SelectField
                        name="download_nums"
                        :label="t('page.shareOptions.file.downloadNums')"
                        :options="[
                            { label: t('page.shareOptions.file.downloadOptions.xdownload', [1]), value: 1 },
                            { label: t('page.shareOptions.file.downloadOptions.xdownload', [2]), value: 2 },
                            { label: t('page.shareOptions.file.downloadOptions.xdownload', [3]), value: 3 },
                            { label: t('page.shareOptions.file.downloadOptions.xdownload', [5]), value: 5 },
                            { label: t('page.shareOptions.file.downloadOptions.xdownload', [10]), value: 10 },
                        ]"
                    />
                    {{ t('page.shareOptions.file.or') }}
                    <SelectField
                        name="expire_time"
                        :label="t('page.shareOptions.file.expireTime')"
                        :options="[
                            { label: t('page.shareOptions.file.expireOptions.5min'), value: 5 },
                            { label: t('page.shareOptions.file.expireOptions.1hour'), value: 60 },
                            { label: t('page.shareOptions.file.expireOptions.1day'), value: 1440 },
                            { label: t('page.shareOptions.file.expireOptions.3days'), value: 4320 },
                        ]"
                    />
                    {{ t('page.shareOptions.file.expireAfter') }}
                </div>
                <div class="flex flex-col gap-1">
                    <div class="flex flex-row gap-3 min-h-9">
                        <SwitchField
                            name="has_pickup_code"
                            :label="t('page.shareOptions.file.pickupCode')"
                            :rules="
                                (value: boolean) => {
                                    if (!!value) {
                                        setFieldValue('has_password', false)
                                    }
                                    return true
                                }
                            "
                        />
                    </div>
                    <div class="flex flex-row gap-3 min-h-9">
                        <SwitchField
                            name="has_password"
                            :label="t('page.shareOptions.file.passwordProtection')"
                            :rules="
                                (value: boolean) => {
                                    if (!!value) {
                                        setFieldValue('has_pickup_code', false)
                                    }
                                    return true
                                }
                            "
                        />
                        <InputField
                            v-if="!!values.has_password"
                            name="password"
                            :placeholder="t('page.shareOptions.file.passwordPlaceholder')"
                            rules="required"
                        />
                    </div>
                    <NotifyConfigField :switchLabel="t('page.shareOptions.file.downloadNotify')" />
                </div>
            </div>
            <FormButton
                @click="
                    async (form) => {
                        onFileHandle({ type: 'file-share', config: values })
                        hide()
                    }
                "
                >{{ t('btn.submit') }}</FormButton
            >
        </div>
    </VeeForm>
</template>
