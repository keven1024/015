<script lang="ts" setup>
import { get } from 'lodash-es'
import VeeForm from '@/components/VeeForm.vue'
import MarkdownInputField from '@/components/Field/MarkdownInputField.vue'

import FormButton from '@/components/Field/FormButton.vue'
import Button from '@/components/ui/button/Button.vue'
import showDrawer from '@/lib/showDrawer'
import { h } from 'vue'
import TextShareDrawer from '@/components/Drawer/TextShareDrawer.vue'
import { cx } from 'class-variance-authority'
</script>
<template>
    <VeeForm v-slot="{ setValues, getValues }">
        <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 gap-5 flex flex-col">
            <div class="text-xl font-normal">输入文本</div>
            <div class="relative">
                <MarkdownInputField name="text" placeholder="使用我们的文本处理器轻松分享，翻译，总结，生成图片，询问大模型"
                    class="max-h-[50vh] min-h-40 overflow-y-auto max-w-full [&>*]:pr-10" rules="required" />
                <Button variant="ghost" size="icon" :class="cx('absolute right-2 top-2 hover:bg-black/10 transition-all duration-300',
                    get(getValues(), 'text')?.length > 0 ? 'opacity-100' : 'opacity-0 pointer-events-none'
                )" @click="() => {
                    setValues({ text: '' })
                }">
                    <LucideX />
                </Button>
            </div>
            <div class="flex flex-row gap-3">
                <FormButton @click="async (form) => {
                    const { text } = form?.values || {}
                    showDrawer({ render: ({ hide }) => h(TextShareDrawer, { hide, text }) })
                }">
                    <LucideShare class="size-4" />提交
                </FormButton>
            </div>
        </div>
    </VeeForm>
</template>