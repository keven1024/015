<script setup lang="ts">
import SwitchField from '../Field/SwitchField.vue'
import InputField from '../Field/InputField.vue'
import SelectField from '../Field/SelectField.vue'
import FormButton from '../Field/FormButton.vue'
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
        <div class="flex flex-col gap-3">
            <h2 class="text-lg font-bold">{{ t('fileshare.title') }}</h2>
            <div class="flex flex-row items-center gap-2 text-sm">
                <SelectField
                    name="download_nums"
                    :label="t('fileshare.downloadNums')"
                    :options="[
                        { label: t('fileshare.downloadOptions.1time'), value: 1 },
                        { label: t('fileshare.downloadOptions.2times'), value: 2 },
                        { label: t('fileshare.downloadOptions.3times'), value: 3 },
                        { label: t('fileshare.downloadOptions.5times'), value: 5 },
                        { label: t('fileshare.downloadOptions.10times'), value: 10 },
                    ]"
                />
                {{ t('fileshare.or') }}
                <SelectField
                    name="expire_time"
                    :label="t('fileshare.expireTime')"
                    :options="[
                        { label: t('fileshare.expireOptions.5min'), value: 5 },
                        { label: t('fileshare.expireOptions.1hour'), value: 60 },
                        { label: t('fileshare.expireOptions.1day'), value: 1440 },
                        { label: t('fileshare.expireOptions.3days'), value: 4320 },
                    ]"
                />
                {{ t('fileshare.expireAfter') }}
            </div>
            <div class="flex flex-col gap-1">
                <div class="flex flex-row gap-3 min-h-9">
                    <SwitchField
                        name="has_pickup_code"
                        :label="t('fileshare.pickupCode')"
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
                        :label="t('fileshare.passwordProtection')"
                        :rules="
                            (value: boolean) => {
                                if (!!value) {
                                    setFieldValue('has_pickup_code', false)
                                }
                                return true
                            }
                        "
                    />
                    <InputField v-if="!!values.has_password" name="password" :placeholder="t('fileshare.passwordPlaceholder')" rules="required" />
                </div>
                <div class="flex flex-row gap-3 min-h-9">
                    <SwitchField name="has_notify" :label="t('fileshare.downloadNotify')" />
                    <InputField v-if="!!values.has_notify" name="notify_email" :placeholder="t('fileshare.emailPlaceholder')" rules="required" />
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
