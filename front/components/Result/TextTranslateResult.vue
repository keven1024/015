<script setup lang="ts">
import VeeForm from '@/components/VeeForm.vue'
import MarkdownInputField from '@/components/Field/MarkdownInputField.vue'
import SelectField from '@/components/Field/SelectField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import MarkdownRender from '@/components/MarkdownRender.vue'
import { Button } from '@/components/ui/button'
import { useClipboard } from '@vueuse/core'
import { toast } from 'vue-sonner'
import type { handleTextComponentProps } from './types'

const props = defineProps<handleTextComponentProps>()
console.log(props?.data)

const { t } = useI18n()
const { copy } = useClipboard()
const translatedText = ref('')

const languageOptions = computed(() => [
    { label: t('page.result.textTranslate.language.auto'), value: 'auto' },
    { label: t('page.result.textTranslate.language.zh-CN'), value: 'zh-CN' },
    { label: t('page.result.textTranslate.language.en'), value: 'en' },
    { label: t('page.result.textTranslate.language.ja'), value: 'ja' },
    { label: t('page.result.textTranslate.language.ko'), value: 'ko' },
])

const providerOptions = computed(() => [
    { label: 'Google', value: 'google' },
    { label: 'Microsoft', value: 'microsoft' },
    { label: 'DeepSeek', value: 'deepseek' },
    // { label: 'DeepLX', value: 'deeplx' },
])

const handleCopyText = async (text?: string) => {
    if (!text) return
    await copy(text)
    toast.success('复制成功')
}

const handleRetranslate = async () => {
    translatedText.value = ''
}
</script>

<template>
    <BaseCard class="flex flex-col gap-4" :title="t('page.result.textTranslate.title')" :showBackButton="true">
        <VeeForm v-slot="{ values, setFieldValue }" :initialValues="{ ...data?.config, input: data?.text }">
            <div class="flex flex-col gap-3">
                <div class="flex flex-col gap-2">
                    <div class="flex flex-row justify-between items-end">
                        <div class="flex flex-col gap-1">
                            <Label class="text-xs">{{ t('page.result.textTranslate.sourceLanguage') }}</Label>
                            <SelectField
                                name="source"
                                class="bg-white/70"
                                :placeholder="t('page.result.textTranslate.sourceLanguage')"
                                default-value="auto"
                                :options="languageOptions"
                                rules="required"
                            />
                        </div>
                        <div class="flex flex-row gap-1">
                            <Button variant="outline" class="bg-white/70" size="icon" @click="handleCopyText(values.input)">
                                <LucideCopy class="size-4" />
                            </Button>
                        </div>
                    </div>
                    <MarkdownInputField name="input" rules="required" class="max-h-[30vh] min-h-40 overflow-y-auto max-w-full" />
                </div>

                <div class="flex flex-row justify-center items-center gap-2">
                    <Button class="px-10">
                        <LucideArrowUpDown />
                    </Button>
                </div>

                <div class="flex flex-col gap-2">
                    <div class="flex justify-between gap-3 items-end">
                        <div class="flex flex-row gap-3">
                            <div class="flex flex-row justify-between">
                                <div class="flex flex-col gap-1">
                                    <Label class="text-xs">{{ t('page.result.textTranslate.targetLanguage') }}</Label>
                                    <SelectField
                                        name="target"
                                        class="bg-white/70"
                                        :placeholder="t('page.result.textTranslate.targetLanguage')"
                                        :options="languageOptions"
                                        rules="required"
                                    />
                                </div>
                            </div>
                            <div class="flex flex-col gap-1">
                                <Label class="text-xs">{{ t('page.result.textTranslate.provider') }}</Label>
                                <SelectField
                                    name="provider"
                                    class="bg-white/70"
                                    :placeholder="t('page.result.textTranslate.provider')"
                                    default-value="mtranslate"
                                    :options="providerOptions"
                                    rules="required"
                                />
                            </div>
                        </div>

                        <Button variant="outline" class="bg-white/70" size="icon" @click="handleCopyText(translatedText)">
                            <LucideCopy class="size-4" />
                        </Button>
                    </div>
                    <div class="rounded-md bg-white/50 min-h-48 p-2">
                        <MarkdownRender :markdown="translatedText" class="prose prose-sm max-w-none min-h-40 max-h-[30vh]" />
                    </div>
                </div>
                <FormButton @click="handleRetranslate" class="w-full">
                    <LucideLanguages class="size-4" />
                    {{ t('page.result.textTranslate.retranslate') }}
                </FormButton>
            </div>
        </VeeForm>
    </BaseCard>
</template>
