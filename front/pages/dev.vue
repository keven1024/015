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
import MultiSelectField from '@/components/Field/MultiSelectField.vue'
import ComboboxField from '@/components/Field/ComboboxField.vue'
import MultiComboboxField from '@/components/Field/MultiComboboxField.vue'
import SwitchField from '@/components/Field/SwitchField.vue'
import DateRangePickerField from '@/components/Field/DateRangePickerField.vue'
import FileUploadBlockProgressView from '@/components/FileUploadBlockProgressView.vue'
import FileUploadHeatMapView from '@/components/FileUploadHeatMapView.vue'

type DevChunkStatus = 'success' | 'error' | 'processing'

type DevUploadInfo = {
    chunks: Record<number, { status: DevChunkStatus; createdAt: number }>
    chunkLength: number
}

const createDevChunks = (chunkLength: number, uploadedChunkLength = chunkLength) => {
    const statuses: DevChunkStatus[] = ['success', 'error', 'processing']
    return Object.fromEntries(
        Array.from({ length: Math.min(chunkLength, uploadedChunkLength) }, (_, index) => [
            index,
            {
                status: statuses[Math.floor(Math.random() * statuses.length)]!,
                createdAt: Date.now() - Math.floor(Math.random() * chunkLength) * 1000,
            },
        ])
    ) as Record<number, { status: DevChunkStatus; createdAt: number }>
}

const devUploadInfo = reactive<DevUploadInfo>({
    chunkLength: 300,
    chunks: createDevChunks(300, 150),
})

const randomizeDevUploadInfo = () => {
    devUploadInfo.chunks = createDevChunks(devUploadInfo.chunkLength, Math.floor(devUploadInfo.chunkLength / 2))
}

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
                        // @ts-ignore vue-sonner's template binding is incorrectly unwrapped by vue-tsc
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
                        // @ts-ignore vue-sonner's template binding is incorrectly unwrapped by vue-tsc
                        toast.success('成功')
                    }
                "
                >异步按钮</AsyncButton
            >
        </div>
        <VeeForm
            :initialValues="{
                select: '1',
                combobox: '1',
                multiCombobox: ['1'],
            }"
        >
            <div class="flex flex-col gap-5 bg-white p-5 rounded-xl">
                <div>表单测试</div>
                <InputField name="input" label="input" rules="required" />
                <div class="flex gap-3">
                    <SelectField
                        name="select"
                        label="select"
                        :options="[
                            { label: 'one', value: '1' },
                            { label: 'two', value: '2' },
                        ]"
                        rules="required"
                    />
                    <MultiSelectField
                        name="multiselect"
                        label="multiselect"
                        :options="[
                            { label: 'one', value: '1' },
                            { label: 'two', value: '2' },
                        ]"
                        rules="required"
                    />
                </div>
                <div class="flex gap-3">
                    <ComboboxField
                        name="combobox"
                        label="combobox"
                        class="w-20"
                        placeholder="请选择"
                        :options="[
                            { label: 'one', value: '1' },
                            { label: 'two', value: '2' },
                        ]"
                        rules="required"
                    />
                    <MultiComboboxField
                        name="multiCombobox"
                        label="multiCombobox"
                        placeholder="请选择"
                        :options="[
                            { label: 'one', value: '1' },
                            { label: 'two', value: '2' },
                        ]"
                        rules="required"
                    />
                </div>
                <div class="flex flex-wrap items-center gap-3">
                    <DateRangePickerField name="dateRange" />
                </div>
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

        <div class="flex flex-col gap-3 rounded-xl bg-white p-5">
            <div class="flex items-center justify-between gap-3">
                <div class="font-semibold">FileUploadBlockProgressView</div>
                <Button size="sm" variant="outline" @click="randomizeDevUploadInfo">随机更新</Button>
            </div>
            <div class="h-7 w-full">
                <FileUploadBlockProgressView :data="devUploadInfo" />
            </div>
            <div class="font-semibold">FileUploadHeatMapView</div>
            <div class="h-36 w-full rounded bg-black/5 p-2">
                <FileUploadHeatMapView :size="12" :data="devUploadInfo" />
            </div>
        </div>
    </BaseCard>
</template>
<style scoped></style>
