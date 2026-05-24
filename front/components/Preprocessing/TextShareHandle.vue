<script setup lang="ts">
import InputField from '../Field/InputField.vue'
import SelectField from '../Field/SelectField.vue'
import SwitchField from '../Field/SwitchField.vue'
import FormButton from '../Field/FormButton.vue'
import NotifyConfigField from './NotifyConfigField.vue'
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
        <div class="flex flex-col gap-3 max-h-[75vh]">
            <h2 class="text-lg font-bold">{{ t('page.shareOptions.text.title') }}</h2>
            <div class="flex flex-col gap-3 flex-1 overflow-y-auto">
                <div class="flex flex-row items-center gap-2 text-sm">
                    <SelectField
                        name="download_nums"
                        :label="t('page.shareOptions.text.viewNums')"
                        :options="[
                            { label: t('page.shareOptions.text.viewOptions.xview', [1]), value: 1 },
                            { label: t('page.shareOptions.text.viewOptions.xview', [2]), value: 2 },
                            { label: t('page.shareOptions.text.viewOptions.xview', [3]), value: 3 },
                            { label: t('page.shareOptions.text.viewOptions.xview', [5]), value: 5 },
                            { label: t('page.shareOptions.text.viewOptions.xview', [10]), value: 10 },
                        ]"
                    />
                    {{ t('page.shareOptions.text.or') }}
                    <SelectField
                        name="expire_time"
                        :label="t('page.shareOptions.text.expireTime')"
                        :options="[
                            { label: t('page.shareOptions.text.expireOptions.5min'), value: 5 },
                            { label: t('page.shareOptions.text.expireOptions.1hour'), value: 60 },
                            { label: t('page.shareOptions.text.expireOptions.1day'), value: 1440 },
                            { label: t('page.shareOptions.text.expireOptions.3days'), value: 4320 },
                        ]"
                    />
                    {{ t('page.shareOptions.text.expireAfter') }}
                </div>
                <div class="flex flex-col gap-1">
                    <div class="flex flex-row gap-3 min-h-9">
                        <SwitchField
                            name="has_pickup_code"
                            :label="t('page.shareOptions.text.pickupCode')"
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
                            :label="t('page.shareOptions.text.passwordProtection')"
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
                            :placeholder="t('page.shareOptions.text.passwordPlaceholder')"
                            rules="required"
                        />
                    </div>
                    <NotifyConfigField :switchLabel="t('page.shareOptions.text.readNotify')" />
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
