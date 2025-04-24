<script lang="ts" setup>
import VeeForm from '@/components/VeeForm.vue'
import FileUploadField from '@/components/Field/FileUploadField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import showDrawer from '@/lib/showDrawer'
import { h } from 'vue'
import FileShareDrawer from '@/components/Drawer/FileShareDrawer.vue'

</script>
<template>
  <VeeForm>
    <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 gap-5 flex flex-col">
      <div class="text-xl font-normal">上传文件</div>
      <FileUploadField name="file" rules="required" />
      <div class="flex flex-row gap-3">
        <FormButton @click="(form) => {
          const { file } = form?.values || {}
          showDrawer({ render: ({ hide }) => h(FileShareDrawer, { hide, file }) })
        }">
          <LucideShare class="size-4" />提交
        </FormButton>
      </div>
    </div>
  </VeeForm>
</template>