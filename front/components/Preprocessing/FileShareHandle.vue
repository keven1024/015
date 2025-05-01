<script setup lang="ts">
import SwitchField from '../Field/SwitchField.vue';
import InputField from '../Field/InputField.vue';
import SelectField from '../Field/SelectField.vue';
import FormButton from '../Field/FormButton.vue';
const props = defineProps<{
    hide: () => void
    file: File
    onFileHandle: ({ type, config }: { type: string, config: any }) => void
}>()

</script>

<template>
    <VeeForm v-slot="{ values }" :initialValues="{ download_nums: 1, expire_time: 1440 }">
        <div class="flex flex-col gap-3">
            <h2 class="text-lg font-bold">分享选项</h2>
            <div class="flex flex-row items-center gap-2 text-sm">

                <SelectField name="download_nums" label="下载次数" :options="[
                    { label: '1次下载', value: 1 },
                    { label: '2次下载', value: 2 },
                    { label: '3次下载', value: 3 },
                    { label: '5次下载', value: 4 },
                    { label: '10次下载', value: 5 },
                ]" />
                或
                <SelectField name="expire_time" label="过期时间" :options="[
                    { label: '5分钟', value: 5 },
                    { label: '1小时', value: 60 },
                    { label: '1天', value: 1440 },
                    { label: '3天', value: 4320 },
                ]" />
                后过期
            </div>
            <div class="flex flex-row gap-3 min-h-9">
                <SwitchField name="has_password" label="密码保护" />
                <InputField v-if="!!values.has_password" name="password" placeholder="请输入密码" />
            </div>
            <div class="flex flex-row gap-3 min-h-9">
                <SwitchField name="has_download_notify" label="下载通知" />
                <InputField v-if="!!values.has_download_notify" name="download_notify_email" placeholder="请输入邮箱" />
            </div>
            <FormButton @click="(form) => {
                onFileHandle({ type: 'file-share', config: values })
                hide()
            }">提交</FormButton>
        </div>
    </VeeForm>
</template>