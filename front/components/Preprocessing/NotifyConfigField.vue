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
import { parseCurl } from 'sweet-curl-parser'

interface WebhookItem {
    id: string
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
    headers: [string, string][]
    body?: string
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
                        <InputField
                            :name="`notify_webhooks.${index}.url`"
                            :label="t('page.shareOptions.notify.webhookUrl')"
                            rules="required|url"
                            @blur="
                                (e: FocusEvent) => {
                                    const input = (e?.target as HTMLInputElement)?.value
                                    if (!input.startsWith('curl ')) return
                                    try {
                                        const { success, data } = parseCurl(input) || {}
                                        if (!success) return
                                        const { url, method, headers, body } = data || {}
                                        setFieldValue(`notify_webhooks.${index}.url`, url?.fullUrl)
                                        setFieldValue(`notify_webhooks.${index}.method`, method?.toUpperCase())
                                        setFieldValue(
                                            `notify_webhooks.${index}.headers`,
                                            headers?.map((h: any) => [h.name, h.value])
                                        )
                                        if (body) setFieldValue(`notify_webhooks.${index}.body`, body)
                                        expandedAdvanced = new Set([...expandedAdvanced, index])
                                    } catch {}
                                }
                            "
                        />
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
                        :config="{
                            key: {
                                placeholder: t('page.shareOptions.notify.webhookHeaderKey'),
                                enum: ['Content-Type', 'User-Agent', 'Authorization', 'Accept', 'Content-Length'],
                            },
                            value: {
                                placeholder: t('page.shareOptions.notify.webhookHeaderValue'),
                                component: [
                                    [
                                        (key: string) => key === 'Content-Type',
                                        ({ ...props }) =>
                                            h(SelectField, { ...props, options: [{ value: 'text/plain' }, { value: 'application/json' }] }),
                                    ],
                                ],
                            },
                        }"
                    />
                    <TextareaField
                        :name="`notify_webhooks.${index}.body`"
                        :label="t('page.shareOptions.notify.webhookBody')"
                        :rows="4"
                        placeholder='{"key": "value"}'
                    />
                </div>
            </div>
            <div class="flex justify-start">
                <Button
                    type="button"
                    size="sm"
                    @click="
                        setFieldValue('notify_webhooks', [
                            ...((values.notify_webhooks as WebhookItem[]) || []),
                            { url: '', method: 'POST', headers: [], body: '' },
                        ])
                    "
                >
                    <LucidePlus class="size-4" />
                    {{ t('common.add') }}
                </Button>
            </div>
        </div>
    </div>
</template>
