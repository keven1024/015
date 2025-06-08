<script setup lang="ts">
import dayjs from "dayjs";
import AsyncButton from "@/components/ui/button/AsyncButton.vue";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";
import { isBoolean } from "lodash-es";
import { LucideCheck, LucideX } from "lucide-vue-next";
import { cx } from "class-variance-authority";
import { toast } from "vue-sonner";
import MarkdownRender from "@/components/MarkdownRender.vue";
dayjs.extend(duration);
dayjs.extend(relativeTime);

const props = defineProps<{
  data: any;
}>();

const { getShareToken } = useMyAppShare();

const expireSeconds = computed(() => {
  return dayjs(props?.data?.expire_at * 10e2).unix() - dayjs().unix();
});

const { remaining, start } = useCountdown(expireSeconds.value);

onMounted(() => {
  start();
});

const fileShareInfo = computed(() => {
  return [
    { label: "需要密码", value: props?.data?.has_password ?? false },
    {
      label: "过期时间",
      value: dayjs.duration(remaining.value, "seconds").format(`D天 HH:mm:ss`),
    },
    { label: "剩余浏览次数", value: props?.data?.download_nums ?? 0 },
  ];
});
const previewText = ref<string | null>(null);

const handlePreview = async () => {
  try {
    const token = await getShareToken(props?.data?.id);
    const r = await $fetch<{
      code: number;
      data: {
        data: string;
      };
    }>(`/api/download?token=${token}`);
    previewText.value = r?.data?.data;
  } catch (error) {
    toast.error((error as any)?.data?.message || error);
  }
};
</script>
<template>
  <div
    :class="
      cx(
        'flex flex-col max-h-full',
        !!previewText ? 'gap-3' : 'gap-16 items-center',
      )
    "
  >
    <h1 class="text-xl">查看文本</h1>
    <template v-if="!previewText">
      <div class="flex flex-col gap-2 md:flex-row w-full">
        <div
          class="flex flex-row md:flex-col md:gap-1 justify-between items-center md:flex-1"
          v-for="item in fileShareInfo"
        >
          <div class="text-xs opacity-75">{{ item?.label }}</div>
          <component
            v-if="isBoolean(item?.value)"
            :is="item?.value ? LucideCheck : LucideX"
            class="size-6"
          />
          <div v-else class="md:text-xl">{{ item?.value }}</div>
        </div>
      </div>
      <div class="w-full">
        <AsyncButton @click="handlePreview" class="w-full">浏览</AsyncButton>
      </div>
    </template>
    <template v-else>
      <MarkdownRender
        :markdown="previewText"
        class="rounded-md bg-white/70 p-3 w-full max-w-full min-h-80 overflow-y-auto"
      />
    </template>
  </div>
</template>
