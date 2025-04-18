<template>
  <div class="flex flex-col gap-5 py-5 items-center w-full h-full">
    <VeeForm>
      <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-200" v-if="type === 'file'">
        1
      </div>
    </VeeForm>
    <VeeForm v-slot="{ setValues, getValues }">
      <div class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 gap-5 flex flex-col"
        v-if="type === 'text'">
        <div class="text-xl font-normal">输入文本</div>
        <div class="relative">
          <MarkdownInputField name="text" class="max-h-[50vh] min-h-40 overflow-y-auto max-w-full [&>*]:pr-10"
            rules="required" />
          <Button variant="ghost" size="icon" class="absolute right-2 top-2 hover:bg-black/10" @click="() => {
            setValues({ text: '' })
            console.log('text', getValues())
          }">
            <LucideX />
          </Button>
        </div>
        <div class="flex flex-row gap-3">
          <FormButton @click="(form) => {
            console.log('text form', form)
            showDrawer({ render: ({ hide }) => h(TextShareDrawer, { hide }) })
          }">发送</FormButton>
        </div>
      </div>
    </VeeForm>
  </div>
</template>

<script setup lang="ts">
import { isString } from 'lodash-es'
import VeeForm from '@/components/VeeForm.vue'
import MarkdownInputField from '@/components/Field/MarkdownInputField.vue'
import FormButton from '@/components/Field/FormButton.vue'
import Button from '@/components/ui/button/Button.vue'
import showDrawer from '@/lib/showDrawer'
import { h } from 'vue'
import TextShareDrawer from '@/components/Drawer/TextShareDrawer.vue'
const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
onMounted(() => {
  if (!isString(type.value) || type.value?.length === 0) {
    router.push({ query: { type: 'file' }, replace: true })
  }
})
</script>