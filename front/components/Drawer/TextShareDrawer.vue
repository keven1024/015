<script setup lang="ts">
import {
  LucideShare,
  LucideImage,
  LucideBot,
  LucideLanguages,
} from "lucide-vue-next";
import { cx } from "class-variance-authority";
import showDrawer from "@/lib/showDrawer";
import TextShareHandle from "@/components/Preprocessing/TextShareHandle.vue";

const props = defineProps<{
  hide: () => void;
  text: string;
  onTextHandle: ({ type, config }: { type: string; config: any }) => void;
}>();
const { t } = useI18n();
const actions = [
  {
    label: t("text.handleType.text-share"),
    icon: LucideShare,
    className: "bg-green-300",
    onClick: () => {
      showDrawer({
        render: ({ hide }) => h(TextShareHandle, { ...props, hide }),
      });
    },
  },
  // {
  //     label: '生成配图', icon: LucideImage, className: 'bg-red-300', onClick: () => {
  //         console.log('复制链接')
  //     }
  // },
  // {
  //     label: '问大模型', icon: LucideBot, className: 'bg-blue-300', onClick: () => {
  //         console.log('复制链接')
  //     }
  // },
  // {
  //     label: '文本翻译', icon: LucideLanguages, className: 'bg-orange-300', onClick: () => {
  //         console.log('复制链接')
  //     }
  // },
];
</script>
<template>
  <div class="flex flex-col gap-5 p-5">
    <div class="flex flex-row gap-2">
      <div
        v-for="item in actions"
        :key="item.label"
        class="flex flex-col items-center gap-2 max-w-20"
        @click="
          () => {
            props?.hide();
            item?.onClick();
          }
        "
      >
        <div
          :class="
            cx(
              'size-14 flex justify-center items-center rounded-full mx-3',
              item?.className,
            )
          "
        >
          <component :is="item?.icon" />
        </div>
        <div class="text-xs truncate w-full text-center">{{ item?.label }}</div>
      </div>
    </div>
  </div>
</template>
