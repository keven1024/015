<script setup lang="ts">
import VeeForm from '@/components/VeeForm.vue'
import FormButton from '@/components/Field/FormButton.vue'
import PinInputField from '@/components/Field/PinInputField.vue'
import type { FormContext, GenericObject } from 'vee-validate'
import { toast } from 'vue-sonner'
const router = useRouter()
const { t } = useI18n()
const props = defineProps<{
    hide: () => void
}>()
const handleSubmit = async (form: FormContext<GenericObject, GenericObject>) => {
    try {
        const code = form.values.code
        const data = await $fetch<{
            data: {
                share_id: string
            }
        }>(`/api/share/pickup/${code}`)
        if (!data.data.share_id) {
            toast.error(t('pickup.codeError'))
            form.resetForm()
            return
        }
        const { share_id } = data.data
        props.hide()
        router.push({
            path: `/s/${share_id}`,
        })
        return
    } catch (error) {
        toast.error(t('pickup.codeError'))
        form.resetForm()
    }
}
</script>

<template>
    <VeeForm>
        <div class="flex flex-col gap-5">
            <div class="text-xl font-bold">{{ t('pickup.title') }}</div>
            <PinInputField
                name="code"
                :rules="
                    (value) => {
                        if (value?.length !== 4) {
                            return false
                        }
                        return true
                    }
                "
            />
            <FormButton @click="handleSubmit">{{ t('btn.submit') }}</FormButton>
        </div>
    </VeeForm>
</template>
