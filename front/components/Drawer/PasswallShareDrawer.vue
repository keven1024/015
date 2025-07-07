<script setup lang="ts">
import VeeForm from '@/components/VeeForm.vue'
import FormButton from '@/components/Field/FormButton.vue'
import InputField from '@/components/Field/InputField.vue'
import type { FormContext, GenericObject } from 'vee-validate'
import { toast } from 'vue-sonner'
const props = defineProps<{
    share_id: string
    hide: any
}>()
const { getShareToken } = useMyAppShare()

const handleSubmit = async (form: FormContext<GenericObject, GenericObject>) => {
    try {
        const password = form.values.password
        const token = await getShareToken(props.share_id, { password })
        if (!token) {
            toast.error('密码错误')
            form.resetForm()
            return
        }
        props.hide?.value(token)
        return
    } catch (error) {
        toast.error('密码错误')
        form.resetForm()
    }
}
</script>

<template>
    <VeeForm>
        <div class="flex flex-col gap-5">
            <div class="text-xl font-bold">输入密码</div>
            <InputField name="password" type="password" rules="required" placeholder="请输入密码" />
            <FormButton @click="handleSubmit">提交</FormButton>
        </div>
    </VeeForm>
</template>
