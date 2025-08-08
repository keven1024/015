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

const { NODE_ENV } = process.env || {}
const isDev = NODE_ENV === 'development'
if (!isDev) {
    navigateTo('/')
}
</script>

<template>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 my-5 flex flex-col gap-5">
        <h1>Dev</h1>
        <div class="flex flex-row gap-5 items-center">
            <Button
                @click="
                    async () => {
                        await showDrawer({
                            render: ({ ...rest }) => h('div', { style: { height: '30vh' } }, '内容'),
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
    </div>
</template>
<style scoped></style>
