<script setup lang="ts">
import { filesize } from "filesize";

const props = defineProps<{
  data: { name: string; value: string; color: string }[];
  title: string;
}>();
const dataKeyMap = {
  file_size: "文件大小",
  file_num: "文件数量",
  processed: "处理数量",
  failed: "失败数量",
};
</script>

<template>
  <div class="rounded-md bg-white p-2 flex flex-col gap-2">
    <div class="text-sm font-medium">{{ title }}</div>
    <div v-for="(item, index) in data" :key="index">
      <div class="flex flex-row items-center gap-2">
        <div
          class="h-5 w-2 rounded-full"
          :style="{ backgroundColor: item.color ?? '#222' }"
        ></div>
        <div class="text-xs font-medium">
          {{ dataKeyMap?.[item.name as keyof typeof dataKeyMap] ?? item.name }}
        </div>
        <div class="text-sm">
          {{
            ["file_size"]?.includes(item?.name)
              ? filesize(item.value)
              : item.value
          }}
        </div>
      </div>
    </div>
  </div>
</template>
