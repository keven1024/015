<script lang="ts" setup>
import showDrawer from '@/lib/showDrawer'
import FileShareDrawer from '@/components/Drawer/FileShareDrawer.vue'
import FileUploadField from '@/components/Field/FileUploadField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import PickupShareBtn from '@/components/PickupShareBtn.vue'

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

const handleFormSubmit = async (form: any) => {
  const { file } = form?.values || {}
  showDrawer({
    render: ({ hide }) => h(FileShareDrawer, {
      hide, file, onFileHandle: ({ type, config }) => {
        form.setFieldValue('handle_type', type)
        form.setFieldValue('config', config)
        emit('change', 'progress')
      }
    })
  })
}
</script>


<template>
    <div class="gap-5 flex flex-col">
        <div class="text-xl font-normal">上传文件</div>
        <FileUploadField name="file" rules="required" />
        <div class="flex flex-row gap-3">
            <FormButton @click="handleFormSubmit">
                <LucideShare class="size-4" />提交
            </FormButton>
            <PickupShareBtn />
        </div>
    </div>
</template>