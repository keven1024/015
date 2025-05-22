<script setup lang="ts">
import VeeForm from '@/components/VeeForm.vue';
import FormButton from '@/components/Field/FormButton.vue';
import PinInputField from '@/components/Field/PinInputField.vue';
import type { FormContext, GenericObject } from 'vee-validate';
import { toast } from 'vue-sonner';
const handleSubmit = async (form: FormContext<GenericObject, GenericObject>) => {
    try {
        const code = form.values.code
        const data = await $fetch<{
            data: {
                share_id: string
            }
        }>(`/api/share/pickup/${code}`)
        if (!data.data.share_id) {
            toast.error('取件码错误')
            form.resetForm()
            return
        }
        return
    } catch (error) {
        toast.error('取件码错误')
        form.resetForm()
    }
}
</script>

<template>
    <VeeForm>
        <div class="flex flex-col gap-5">
            <div class="text-xl font-bold">输入取件码</div>
            <PinInputField name="code" :rules="(value) => {
                if (value.length !== 4) {
                    return false
                }
                return true
            }" />
            <FormButton @click="handleSubmit">提交</FormButton>
        </div>
    </VeeForm>
</template>