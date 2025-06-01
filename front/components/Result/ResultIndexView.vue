<script lang="ts" setup>
import FileShareResult from "@/components/Result/FileShareResult.vue";
import TextShareResult from "@/components/Result/TextShareResult.vue";
import ImageCompressResult from "@/components/Result/ImageCompressResult.vue";

type basehandleData = { config: any; handle_type: string };

type filehandleData = { file: File; file_id: string } & basehandleData;
type texthandleData = { text: string } & basehandleData;

const props = defineProps<{
  data: filehandleData | texthandleData;
}>();

const emit = defineEmits<{
  (e: "change", key: string): void;
}>();

const handleList = [
  { component: FileShareResult, key: "file-share" },
  { component: TextShareResult, key: "text-share" },
  { component: ImageCompressResult, key: "file-image-compress" },
];

const activeHandle = computed(() => {
  return handleList.find((item) => item.key === props?.data?.handle_type);
});
// vue这个ts蠢的没边了，本来想写component: FileShareResult | TextShareResult，结果不行
</script>
<template>
  <div>
    <component
      v-if="'file' in data"
      :is="activeHandle?.component"
      :data="data"
      @change="(key: string) => emit('change', key)"
    />
    <component
      v-if="'text' in data"
      :is="activeHandle?.component"
      :data="data"
      @change="(key: string) => emit('change', key)"
    />
  </div>
</template>
