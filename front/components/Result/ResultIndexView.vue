<script lang="ts" setup>
import FileShareResult from '@/components/Result/FileShareResult.vue'
type basehandleData = { config: any, handle_type: string }

type filehandleData = { file: File, file_id: string } & basehandleData
type texthandleData = { text: string } & basehandleData

const props = defineProps<{
    data: filehandleData 
}>()

const emit = defineEmits<{
    (e: 'change', key: string): void
}>()

// console.log(props.data)

const handleList = [
    { component: FileShareResult, key: 'file-share' },
    // { component: FileShareResult, key: 'file-share' },
]
const handleComponent = computed(() => {
    return handleList.find((item) => item.key === props?.data?.handle_type)?.component
})

</script>
<template>
    <div class="">
        <component :is="handleComponent" :data="data" @change="(key: string) => {
            emit('change', key)
        }" />
    </div>
</template>