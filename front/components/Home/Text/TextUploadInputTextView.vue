<script lang="ts" setup>
import MarkdownInputField from '@/components/Field/MarkdownInputField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import Button from '@/components/ui/button/Button.vue'
import showDrawer from '@/lib/showDrawer'
import { h } from 'vue'
import TextShareDrawer from '@/components/Drawer/TextShareDrawer.vue'
import { cx } from 'class-variance-authority'

const form = useFormContext()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

const handleTextShare = ({ type, config }: { type: string, config: any }) => {
    form?.setFieldValue('handle_type', type)
    form?.setFieldValue('config', config)
    emit('change', 'result')
}
</script>
<template>
    <div class="gap-5 flex flex-col">
        <div class="text-xl font-normal">输入文本</div>
        <div class="relative">
            <MarkdownInputField name="text" placeholder="使用我们的文本处理器轻松分享，翻译，总结，生成图片，询问大模型"
                class="max-h-[50vh] min-h-40 overflow-y-auto max-w-full [&>*]:pr-10" rules="required" />
            <Button variant="ghost" size="icon" :class="cx('absolute right-2 top-2 hover:bg-black/10 transition-all duration-300',
                form?.values.text?.length > 0 ? 'opacity-100' : 'opacity-0 pointer-events-none'
            )" @click="() => {
                form?.setValues({ text: '' })
            }">
                <LucideX />
            </Button>
        </div>
        <div class="flex flex-row gap-3">
            <FormButton @click="async (form) => {
                const { text } = form?.values || {}
                showDrawer({ render: ({ hide }) => h(TextShareDrawer, { hide, text, onTextHandle: handleTextShare }) })
            }">
                <LucideShare class="size-4" />提交
            </FormButton>
        </div>
    </div>
</template>