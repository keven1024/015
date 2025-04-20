<template>
  <div class="flex flex-col gap-5 py-5 items-center w-full h-full">
    <VeeForm>
      <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 gap-5 flex flex-col"
        v-if="type === 'file'">
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
    <VeeForm v-slot="{ setValues, getValues }">
      <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 gap-5 flex flex-col"
        v-if="type === 'text'">
        <div class="text-xl font-normal">输入文本</div>
        <div class="relative">
          <MarkdownInputField name="text" placeholder="使用我们的文本处理器轻松分享，翻译，总结，生成图片，询问大模型"
            class="max-h-[50vh] min-h-40 overflow-y-auto max-w-full [&>*]:pr-10" rules="required" />
          <Button variant="ghost" size="icon" :class="cx('absolute right-2 top-2 hover:bg-black/10 transition-all duration-300',
            get(getValues(), 'text')?.length > 0 ? 'opacity-100' : 'opacity-0 pointer-events-none'
          )" @click="() => {
            setValues({ text: '' })
          }">
            <LucideX />
          </Button>
        </div>
        <div class="flex flex-row gap-3">
          <FormButton @click="(form) => {
            const { text } = form?.values || {}
            showDrawer({ render: ({ hide }) => h(TextShareDrawer, { hide, text }) })
          }">
            <LucideShare class="size-4" />提交
          </FormButton>
        </div>
      </div>
    </VeeForm>
  </div>
</template>

<script setup lang="ts">
import { isString, get } from 'lodash-es'
import VeeForm from '@/components/VeeForm.vue'
import MarkdownInputField from '@/components/Field/MarkdownInputField.vue'
import FileUploadField from '@/components/Field/FileUploadField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import Button from '@/components/ui/button/Button.vue'
import showDrawer from '@/lib/showDrawer'
import { h } from 'vue'
import TextShareDrawer from '@/components/Drawer/TextShareDrawer.vue'
import FileShareDrawer from '@/components/Drawer/FileShareDrawer.vue'
import { cx } from 'class-variance-authority'

const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
onMounted(() => {
  if (!isString(type.value) || type.value?.length === 0) {
    router.push({ query: { type: 'file' }, replace: true })
  }
})
</script>