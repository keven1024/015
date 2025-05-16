<script lang="ts" setup>
import VeeForm from '@/components/VeeForm.vue'
import FileUploadInputFileView from './FileUploadInputFileView.vue'
import FileUploadProgressView from './FileUploadProgressView.vue'
import FileUploadResultView from './FileUploadResultView.vue'

const fileStepList = [
    { component: FileUploadInputFileView, key: 'input' },
    { component: FileUploadProgressView, key: 'progress' },
    { component: FileUploadResultView, key: 'result' },
]

const step = ref('input')

const renderComponent = computed(() => {
    return fileStepList.find((item) => item.key === step.value)?.component
})
const formRef = ref<InstanceType<typeof VeeForm>>()
watch(() => step.value, (newVal) => {
    if (newVal === 'input') {
        formRef.value?.form?.resetForm()
        formRef.value?.form?.setValues({ file: null })
    }
})
</script>
<template>
    <VeeForm ref="formRef" v-slot="{ values }" :keepValues="true">
        <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200">
            <component :is="renderComponent" :data="values" @change="(key: string) => {
                step = key
            }" />
        </div>
    </VeeForm>
</template>