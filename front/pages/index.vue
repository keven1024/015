<template>
  <div class="flex flex-col gap-5 py-5 items-center w-full h-full">
    <component :is="renderComponent" />
  </div>
</template>

<script setup lang="ts">
import FileUploadView from '~/components/Home/File/FileUploadIndexView.vue'
import TextUploadView from '~/components/Home/Text/TextUploadIndexView.vue'
import { isString } from 'lodash-es'
const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
onMounted(() => {
  if (!isString(type.value) || type.value?.length === 0) {
    router.push({ query: { type: 'file' }, replace: true })
  }
})

const renderList = [
  { key: 'file', component: FileUploadView },
  { key: 'text', component: TextUploadView },
]
const renderComponent = computed(() => renderList.find(item => item.key === type.value)?.component)
</script>