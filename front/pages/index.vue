<template>
  <div class="flex flex-col gap-5 py-5 items-center w-full h-full">
    <FileUploadView v-if="type === 'file'" />
    <TextUploadView v-if="type === 'text'" />
  </div>
</template>

<script setup lang="ts">
import FileUploadView from '~/components/Upload/File/FileUploadIndexView.vue'
import TextUploadView from '~/components/Upload/TextUploadView.vue'
import { isString } from 'lodash-es'

const route = useRoute()
const router = useRouter()
const type = computed(() => route?.query?.type)
onMounted(() => {
  if (!isString(type.value) || type.value?.length === 0) {
    router.push({ query: { type: 'file' }, replace: true })
  }
})
</script>