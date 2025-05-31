<script setup lang="ts">
import { cx } from "class-variance-authority";

const { availableLocales, setLocale, locale: currentLocale, t } = useI18n();
const route = useRoute();

const localeMap = {
  "zh-CN": "简体中文",
  en: "English",
  // 'ja': '日本語',
  // 'ko': '한국어',
  // 'fr': 'Français',
  // 'de': 'Deutsch',
};

const switchLocale = async (locale: string) => {
  await setLocale(locale as keyof typeof localeMap);
  navigateTo(route.path, {
    external: true,
  });
};
</script>

<template>
  <div class="flex flex-col gap-1 py-2">
    <div class="text-xl font-bold mb-3">{{ t("i18n.switchLocale") }}</div>
    <div
      v-for="locale in availableLocales"
      :key="locale"
      :class="
        cx(
          'rounded-md hover:bg-black/10 p-2 cursor-pointer',
          currentLocale === locale && 'bg-black/10 font-bold',
        )
      "
      @click="() => switchLocale(locale)"
    >
      {{ localeMap?.[locale as keyof typeof localeMap] }}
    </div>
  </div>
</template>
