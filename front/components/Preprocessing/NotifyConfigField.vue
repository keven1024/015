<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { useFormContext } from 'vee-validate'
import SelectField from '../Field/SelectField.vue'
import SwitchField from '../Field/SwitchField.vue'
import InputGroupField from '../Field/InputGroupField.vue'
import InputField from '../Field/InputField.vue'
import KvInputField from '../Field/KvInputGroupField.vue'
import TextareaField from '../Field/TextareaField.vue'

interface WebhookItem {
    id: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
    headers: Record<string, any>
    bodyType: 'none' | 'form-data' | 'raw'
    body: string
}

const { t } = useI18n()
const { values, setFieldValue } = useFormContext()
const expandedAdvanced = ref<Set<number>>(new Set())
</script>

<template>
    <div class="flex flex-col gap-3">
        <div class="flex flex-row gap-3 min-h-9 items-center">
            <SwitchField name="has_notify" :label="t('page.shareOptions.file.downloadNotify')" />
            <SelectField
                v-if="values.has_notify"
                name="notify_types"
                :placeholder="t('page.shareOptions.notify.notifyVia')"
                multiple
                :options="[
                    { label: t('page.shareOptions.notify.email'), value: 'email' },
                    { label: t('page.shareOptions.notify.webhook'), value: 'webhook' },
                ]"
            />
        </div>
        <div v-if="!!values.has_notify && values.notify_types?.includes('email')">
            <InputGroupField
                name="notify_emails"
                :placeholder="t('page.shareOptions.notify.emailPlaceholder')"
                :label="t('page.shareOptions.notify.email')"
                rules="email"
            />
        </div>
        <div v-if="!!values.has_notify && values.notify_types?.includes('webhook')" class="flex flex-col gap-2">
            <Label>Webhook</Label>
            <div v-for="(_, index) in (values.notify_webhooks as WebhookItem[]) || []" :key="index" class="flex flex-col gap-2 border rounded-md p-3">
                <div class="flex flex-row gap-2 items-end">
                    <div class="flex flex-col gap-2">
                        <Label>{{ t('page.shareOptions.notify.webhookMethod') }}</Label>
                        <SelectField
                            :name="`notify_webhooks.${index}.method`"
                            :label="t('page.shareOptions.notify.webhookMethod')"
                            default-value="POST"
                            :options="[
                                { label: 'GET', value: 'GET' },
                                { label: 'POST', value: 'POST' },
                                { label: 'PUT', value: 'PUT' },
                                { label: 'PATCH', value: 'PATCH' },
                                { label: 'DELETE', value: 'DELETE' },
                            ]"
                            class="w-28"
                        />
                    </div>
                    <div class="flex-1">
                        <InputField :name="`notify_webhooks.${index}.url`" :label="t('page.shareOptions.notify.webhookUrl')" rules="required|url" />
                    </div>
                    <Button
                        type="button"
                        variant="outline"
                        size="icon"
                        @click="
                            expandedAdvanced = new Set(
                                expandedAdvanced.has(index)
                                    ? [...expandedAdvanced].filter((expandedIndex) => expandedIndex !== index)
                                    : [...expandedAdvanced, index]
                            )
                        "
                    >
                        <LucideSettings class="size-4" />
                    </Button>
                    <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        class="bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white"
                        @click="
                            () => {
                                setFieldValue(
                                    'notify_webhooks',
                                    ((values.notify_webhooks as WebhookItem[]) || []).filter((_, itemIndex) => itemIndex !== index)
                                )
                                expandedAdvanced = new Set(
                                    [...expandedAdvanced]
                                        .filter((expandedIndex) => expandedIndex !== index)
                                        .map((expandedIndex) => (expandedIndex > index ? expandedIndex - 1 : expandedIndex))
                                )
                            }
                        "
                    >
                        <LucideTrash class="size-4" />
                    </Button>
                </div>

                <div v-show="expandedAdvanced.has(index)" class="flex flex-col gap-3 rounded-md border p-3">
                    <KvInputField
                        :name="`notify_webhooks.${index}.headers`"
                        :label="t('page.shareOptions.notify.webhookHeaders')"
                        :key-placeholder="t('page.shareOptions.notify.webhookHeaderKey')"
                        :value-placeholder="t('page.shareOptions.notify.webhookHeaderValue')"
                    />
                    <div v-if="((values.notify_webhooks as WebhookItem[]) || [])[index]?.method === 'POST'" class="flex flex-col gap-3">
                        <Label class="w-20 pt-2">{{ t('page.shareOptions.notify.webhookBody') }}</Label>
                        <div class="flex flex-1 flex-col gap-2">
                            <SelectField
                                :name="`notify_webhooks.${index}.bodyType`"
                                :placeholder="t('page.shareOptions.notify.webhookBodyType')"
                                :label="t('page.shareOptions.notify.webhookBodyType')"
                                :options="[
                                    { label: t('page.shareOptions.notify.webhookBodyTypeNone'), value: 'none' },
                                    { label: t('page.shareOptions.notify.webhookBodyTypeFormData'), value: 'form-data' },
                                    { label: t('page.shareOptions.notify.webhookBodyTypeRaw'), value: 'raw' },
                                ]"
                            />
                            <TextareaField
                                v-if="['form-data', 'raw'].includes(((values.notify_webhooks as WebhookItem[]) || [])[index]?.bodyType || '')"
                                :name="`notify_webhooks.${index}.body`"
                                :placeholder="t('page.shareOptions.notify.webhookBody')"
                            />
                        </div>
                    </div>
                </div>
            </div>
            <div class="flex justify-start">
                <Button
                    type="button"
                    size="sm"
                    @click="
                        setFieldValue('notify_webhooks', [
                            ...((values.notify_webhooks as WebhookItem[]) || []),
                            { url: '', method: 'POST', headers: {}, bodyType: 'none', body: '' },
                        ])
                    "
                >
                    <LucidePlus class="size-4" />
                    添加
                </Button>
            </div>
        </div>
    </div>
</template>
