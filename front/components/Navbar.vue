<template>
  <div
    class="flex flex-row bg-white/50 backdrop-blur-xl p-2 rounded-full gap-1"
  >
    <div
      v-for="item in routes"
      :key="item.key"
      :class="
        cx(
          'flex flex-row items-center text-sm px-4 py-2 font-bold rounded-full relative select-none cursor-pointer',
          !isActive(item) && 'hover:bg-black/5',
          item?.name && 'gap-2',
          item?.className,
        )
      "
      @click="handleClick(item)"
    >
      <motion.div
        v-if="isActive(item)"
        layoutId="navbar-active"
        class="absolute inset-0 rounded-full w-full h-full bg-black/10"
      />
      <component :is="item.icon" />
      <div class="hidden sm:block">{{ item.name }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { cx } from "class-variance-authority";
import { LucideClipboardType, LucidePaperclip } from "#components";
import { motion } from "motion-v";
import { LucideGlobe } from "lucide-vue-next";
import showDrawer from "@/lib/showDrawer";
import I18nSwitchDrawer from "./Drawer/I18nSwitchDrawer.vue";
const { t } = useI18n();
const routes = [
  {
    key: "about",
    icon: () =>
      h("img", {
        class: "size-10 rounded-full border-2 border-white/50",
        src: "/logo.png",
      }),
    onClick: () => {
      router.push("/about");
    },
    isActive: (item: { key: string }) => route.path?.endsWith(item.key),
    className: "!p-1.5",
  },
  { name: t("navbar.file"), key: "file", icon: LucidePaperclip },
  { name: t("navbar.text"), key: "text", icon: LucideClipboardType },
  {
    key: "i18n",
    icon: LucideGlobe,
    onClick: () => {
      showDrawer({
        render: () => h(I18nSwitchDrawer),
      });
    },
    className: "size-12 !p-1.5 justify-center items-center",
  },
];
const route = useRoute();
const router = useRouter();
const type = computed(() => route?.query?.type);

const isActive = (item: {
  key: string;
  isActive?: (item: { key: string }) => boolean;
}) => {
  const { key, isActive } = item || {};
  return isActive ? isActive(item) : type.value === key;
};

const handleClick = (item: { key: string; onClick?: () => void }) => {
  const { key, onClick } = item || {};
  if (onClick) {
    onClick();
    return;
  }
  router.push({
    path: "/",
    query: {
      ...route.query,
      type: key,
    },
  });
};
</script>
