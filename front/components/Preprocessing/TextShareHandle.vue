<script setup lang="ts">
import SwitchField from '../Field/SwitchField.vue'
import InputField from '../Field/InputField.vue'
import SelectField from '../Field/SelectField.vue'
import FormButton from '../Field/FormButton.vue'
import type { TextShareHandleProps } from './types'
const { t } = useI18n()
const props = defineProps<{
    hide: () => void
    text: string
    onTextHandle: (props: TextShareHandleProps) => void
}>()
</script>

<template>
    <VeeForm v-slot="{ values, setFieldValue }" :initialValues="{ download_nums: 1, expire_time: 1440 }">
        <div class="flex flex-col gap-3">
            <h2 class="text-lg font-bold">{{ t('textshare.title') }}</h2>
            <div class="flex flex-row items-center gap-2 text-sm">
                <SelectField
                    name="download_nums"
                    :label="t('textshare.viewNums')"
                    :options="[
                        { label: t('textshare.viewOptions.1time'), value: 1 },
                        { label: t('textshare.viewOptions.2times'), value: 2 },
                        { label: t('textshare.viewOptions.3times'), value: 3 },
                        { label: t('textshare.viewOptions.5times'), value: 5 },
                        { label: t('textshare.viewOptions.10times'), value: 10 },
                    ]"
                />
                {{ t('textshare.or') }}
                <SelectField
                    name="expire_time"
                    :label="t('textshare.expireTime')"
                    :options="[
                        { label: t('textshare.expireOptions.5min'), value: 5 },
                        { label: t('textshare.expireOptions.1hour'), value: 60 },
                        { label: t('textshare.expireOptions.1day'), value: 1440 },
                        { label: t('textshare.expireOptions.3days'), value: 4320 },
                    ]"
                />
                {{ t('textshare.expireAfter') }}
            </div>
            <div class="flex flex-col gap-1">
                <div class="flex flex-row gap-3 min-h-9">
                    <SwitchField
                        name="has_pickup_code"
                        :label="t('textshare.pickupCode')"
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
                        :label="t('textshare.passwordProtection')"
                        :rules="
                            (value: boolean) => {
                                if (!!value) {
                                    setFieldValue('has_pickup_code', false)
                                }
                                return true
                            }
                        "
                    />
                    <InputField v-if="!!values.has_password" name="password" :placeholder="t('textshare.passwordPlaceholder')" rules="required" />
                </div>
                <div class="flex flex-row gap-3 min-h-9">
                    <SwitchField name="has_notify" :label="t('textshare.readNotify')" />
                    <InputField v-if="!!values.has_notify" name="notify_email" :placeholder="t('textshare.emailPlaceholder')" rules="required" />
                </div>
            </div>
            <FormButton
                @click="
                    async (form) => {
                        onTextHandle({ type: 'text-share', config: values })
                        hide()
                    }
                "
                >{{ t('btn.submit') }}</FormButton
            >
        </div>
    </VeeForm>
</template>
