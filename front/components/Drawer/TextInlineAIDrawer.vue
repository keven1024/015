<script setup lang="ts">
import VeeForm from '../VeeForm.vue'
import FormButton from '@/components/Field/FormButton.vue'
import InputField from '@/components/Field/InputField.vue'
import { LucideCheckCheck, LucideLanguages, LucideSparkle, LucideWandSparkles } from 'lucide-vue-next'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import asyncWait from '~/lib/asyncWait'

const props = defineProps<{
    hide: () => void
    data: { html: string }
}>()
const formRef = ref<InstanceType<typeof VeeForm>>()

const actions = [
    {
        icon: LucideWandSparkles,
        label: '总结',
        onClick: () => {
            formRef.value?.form?.setFieldValue('input', '总结这段话')
        },
    },
    {
        icon: LucideCheckCheck,
        label: '修正拼写和语法错误',
        onClick: () => {
            formRef.value?.form?.setFieldValue('input', '修正这段话的拼写和语法错误')
        },
    },
    {
        icon: LucideLanguages,
        label: '翻译成',
        children: [
            {
                label: '英语',
                onClick: () => {
                    formRef.value?.form?.setFieldValue('input', '把这段话翻译成英语')
                },
            },
            {
                label: '韩文',
                onClick: () => {
                    formRef.value?.form?.setFieldValue('input', '把这段话翻译成韩文')
                },
            },
            {
                label: '中文',
                onClick: () => {
                    formRef.value?.form?.setFieldValue('input', '把这段话翻译成中文')
                },
            },
        ],
    },
]
</script>

<template>
    <div class="flex flex-col gap-5">
        <div class="text-xl font-bold">询问AI</div>
        <div class="overflow-y-auto max-h-[160px] p-3 rounded-md bg-primary/5 border border-primary/10 text-black/80 text-sm">
            <div v-html="data.html" />
        </div>
        <VeeForm ref="formRef">
            <InputField name="input" placeholder="您想如何处理该文本" rules="required" />
            <div class="flex flex-row gap-3">
                <template v-for="action in actions">
                    <Button variant="outline" v-if="!action?.children" @click="action.onClick">
                        <component :is="action.icon" /> {{ action.label }}
                    </Button>
                    <DropdownMenu v-else>
                        <DropdownMenuTrigger>
                            <Button variant="outline"> <component :is="action.icon" /> {{ action.label }} </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent>
                            <DropdownMenuItem v-for="item in action?.children" @click="item.onClick">{{ item?.label }}</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </template>
            </div>

            <FormButton
                @click="
                    async () => {
                        await asyncWait(3000)
                    }
                "
            >
                <LucideSparkle />
                生成
            </FormButton>
        </VeeForm>
    </div>
</template>
