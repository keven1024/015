<script setup lang="ts">
import { Button, AsyncButton } from '@/components/ui/button'
import asyncWait from '~/lib/asyncWait'
import { toast } from 'vue-sonner'
import showDrawer from '~/lib/showDrawer'
import { h } from 'vue'
import VeeForm from '@/components/VeeForm.vue'
import InputField from '@/components/Field/InputField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import SelectField from '@/components/Field/SelectField.vue'
import SwitchField from '@/components/Field/SwitchField.vue'
import dayjs from 'dayjs'

const { NODE_ENV } = process.env || {}
const isDev = NODE_ENV === 'development'
if (!isDev) {
    navigateTo('/')
}
</script>

<template>
    <BaseCard class="my-5 flex flex-col gap-5" title="dev">
        <div class="flex flex-row gap-5 items-center">
            <Button
                @click="
                    async () => {
                        await showDrawer({
                            render: ({ hide }) =>
                                h('div', { class: 'flex h-[30vh] flex-col gap-4' }, [
                                    h('div', '第一层内容'),
                                    h(
                                        Button,
                                        {
                                            variant: 'outline',
                                            onClick: () =>
                                                showDrawer({
                                                    render: ({ hide: hideInner }) =>
                                                        h('div', { class: 'flex h-[30vh] flex-col gap-4' }, [
                                                            h('div', '第二层内容'),
                                                            h(
                                                                Button,
                                                                {
                                                                    onClick: () => hideInner(),
                                                                },
                                                                () => '关闭第二层'
                                                            ),
                                                        ]),
                                                }),
                                        },
                                        () => '打开第二层(showDrawer)'
                                    ),
                                    h(
                                        Button,
                                        {
                                            onClick: () => hide(),
                                        },
                                        () => '关闭第一层'
                                    ),
                                ]),
                        })
                        toast.success('被关闭')
                    }
                "
                >普通按钮(showDrawer)</Button
            >
            <AsyncButton
                variant="outline"
                @click="
                    async () => {
                        await asyncWait(3000)
                        toast.success('成功')
                    }
                "
                >异步按钮</AsyncButton
            >
        </div>
        <VeeForm
            :initialValues="{
                select: '1',
            }"
        >
            <div class="flex flex-col gap-5 bg-white p-5 rounded-xl">
                <div>表单测试</div>
                <InputField name="input" label="input" rules="required" />
                <SelectField
                    name="select"
                    label="select"
                    :options="[
                        { label: 'one', value: '1' },
                        { label: 'two', value: '2' },
                    ]"
                    rules="required"
                />
                <SwitchField name="switch" label="我同意xxx" :rules="(v: any) => !!v" />
                <FormButton
                    @click="
                        async (form) => {
                            console.log('form值', form.values)
                        }
                    "
                >
                    获取表单值(console.log)
                </FormButton>
            </div>
        </VeeForm>
    </BaseCard>
</template>
<style scoped></style>
