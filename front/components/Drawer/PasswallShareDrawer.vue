<script setup lang="ts">
import VeeForm from '@/components/VeeForm.vue'
import FormButton from '@/components/Field/FormButton.vue'
import InputField from '@/components/Field/InputField.vue'
import type { FormContext, GenericObject } from 'vee-validate'
import { toast } from 'vue-sonner'
const { t } = useI18n()
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
            toast.error(t('page.shareView.passwall.passwordError'))
            form.resetForm()
            return
        }
        props?.hide(token)
        return
    } catch (error) {
        toast.error(t('page.shareView.passwall.passwordError'))
        form.resetForm()
    }
}
</script>

<template>
    <VeeForm>
        <div class="flex flex-col gap-5">
            <div class="text-xl font-bold">{{ t('page.shareView.passwall.title') }}</div>
            <InputField name="password" type="password" rules="required" :placeholder="t('page.shareView.passwall.passwordPlaceholder')" />
            <FormButton @click="handleSubmit">{{ t('btn.submit') }}</FormButton>
        </div>
    </VeeForm>
</template>
