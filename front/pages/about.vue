<script setup lang="ts">
import { CurveType } from "@unovis/ts";
import { AreaChart } from "@/components/ui/chart-area";
import { cx } from "class-variance-authority";
import { useQuery } from "@tanstack/vue-query";
import { Skeleton } from "@/components/ui/skeleton";
import AboutChartTooltip from "@/components/AboutChartTooltip.vue";
import { filesize } from "filesize";
import SparkMD5 from "spark-md5";
import useMyAppConfig from "@/composables/useMyAppConfig";

const appConfig = useMyAppConfig();
const { site_title, site_desc } = appConfig.value || {};

const { data, isLoading } = useQuery({
  queryKey: ["stat"],
  queryFn: async () => {
    const data = await $fetch<{ data: any }>("/api/stat");
    return data.data;
  },
});

const { t } = useI18n();

const chartTabs = computed(() => {
  return [
    {
      label: t("about.file"),
      value: "storage",
      total:
        data.value?.chart?.storage?.reduce(
          (acc: number, curr: { file_size: number; file_num: number }) =>
            acc + curr.file_num,
          0,
        ) ?? 0,
    },
    {
      label: t("about.task"),
      value: "queue",
      total:
        data.value?.chart?.queue?.reduce(
          (acc: number, curr: { processed: number; failed: number }) =>
            acc + curr.processed + curr.failed,
          0,
        ) ?? 0,
    },
  ];
});

const currentFileSize = computed(() => {
  return (
    data.value?.chart?.storage?.reduce(
      (acc: number, curr: { file_size: number; file_num: number }) =>
        acc + curr.file_size,
      0,
    ) ?? 0
  );
});

const currentChartTab = ref<"storage" | "queue">("storage");
const currentChartData = computed(() => {
  const { storage, queue } = data.value?.chart || {};
  if (currentChartTab.value === "storage") {
    return {
      data: storage,
      index: "date",
      categories: ["file_size", "file_num"],
      colors: ["#22d3ee", "#c084fc"],
    };
  }
  return {
    data: queue,
    index: "date",
    categories: ["processed", "failed"],
    colors: ["#4ade80", "#f87171"],
  };
});

const genUserAvatar = ({ email }: { email: string }) => {
  if (!email) {
    return "/logo.png";
  }
  return `https://www.gravatar.com/avatar/${SparkMD5.hash(email)}?d=retro`;
};

const handleUserClick = ({ url, email }: { url: string; email: string }) => {
  if (url) {
    return navigateTo(url, { external: true });
  }
  if (email) {
    return navigateTo(`mailto:${email}`, { external: true });
  }
  return null;
};

const users = computed(() => {
  const { email, name, url } = data.value?.admin || {};
  return [
    ...(!!name
      ? [
          {
            title: t("about.admin"),
            email,
            name,
            url: url ?? (email ? `mailto:${email}` : null),
          },
        ]
      : []),
    {
      title: t("about.author"),
      name: "keven1024",
      email: "keven@fudaoyuan.icu",
      url: "https://github.com/keven1024",
    },
  ];
});
</script>

<template>
  <div
    class="rounded-xl p-5 bg-white/50 backdrop-blur-xl w-full lg:w-200 my-5 flex flex-col gap-5"
  >
    <div class="text-xl font-normal">{{ t("about.title") }}</div>
    <div class="flex flex-col gap-2 items-center">
      <NuxtImg src="/logo.png" class="size-20 rounded-xl" />
      <div class="text-xl">{{ site_title ?? "015" }}</div>
      <div class="text-sm opacity-75 text-center px-5">
        {{ site_desc ?? t("seo.desc") }}
      </div>
    </div>
    <div class="font-semibold">{{ t("about.systemInfo") }}</div>
    <template v-if="isLoading">
      <div class="flex flex-row gap-2">
        <Skeleton class="w-full h-20 rounded-xl" v-for="i in 2" :key="i" />
      </div>
    </template>
    <template v-else>
      <div class="grid grid-cols-2 gap-2">
        <div class="rounded-xl bg-white/50 flex-1 flex flex-col p-3">
          <div class="opacity-75 text-xs">{{ t("about.systemVersion") }}</div>
          <div class="text-xl font-semibold">
            {{ data?.version ?? "dev" }}
          </div>
        </div>
        <div class="rounded-xl bg-white/50 flex-1 flex flex-col p-3">
          <div class="opacity-75 text-xs">{{ t("about.storage") }}</div>
          <div class="text-right flex flex-row items-baseline">
            <span class="text-lg font-semibold">{{
              filesize(currentFileSize ?? 0)
            }}</span>
            <span class="text-md opacity-75"
              >/ {{ filesize(data?.max_limit?.file_size ?? 0) }}</span
            >
          </div>
          <div class="rounded-full w-full h-1 bg-black/10">
            <div
              class="rounded-full h-full bg-blue-500"
              :style="{
                width: `${(currentFileSize / (data?.max_limit?.file_size ?? 0)) * 100}%`,
              }"
            ></div>
          </div>
        </div>
      </div>
    </template>
    <div class="font-semibold">{{ t("about.analysis") }}</div>
    <template v-if="isLoading">
      <div class="flex flex-row gap-2">
        <Skeleton class="w-full h-96 rounded-xl" />
      </div>
    </template>
    <template v-else>
      <div class="flex flex-col gap-2 bg-white/50 w-full rounded-xl py-5">
        <div class="flex flex-row gap-2 px-5">
          <div
            :class="
              cx(
                'rounded-md min-w-30 flex flex-col px-3 py-1.5 cursor-pointer',
                currentChartTab === tab.value && 'bg-black/10',
              )
            "
            v-for="tab in chartTabs"
            :key="tab.value"
            @click="
              () => {
                currentChartTab = tab.value as 'storage' | 'queue';
              }
            "
          >
            <div class="opacity-75 text-xs">{{ tab.label }}</div>
            <div class="text-lg font-semibold">{{ tab.total }}</div>
          </div>
        </div>
        <AreaChart
          v-if="currentChartData"
          class="h-64 w-full"
          :key="currentChartTab"
          :index="currentChartData.index"
          :data="currentChartData.data"
          :categories="currentChartData.categories"
          :show-grid-line="false"
          :show-legend="false"
          :show-y-axis="true"
          :show-x-axis="true"
          :colors="currentChartData.colors"
          :custom-tooltip="AboutChartTooltip"
          :curve-type="CurveType.CatmullRom"
        />
      </div>
    </template>
    <template v-if="isLoading">
      <div class="flex flex-row gap-2">
        <Skeleton class="w-full h-20 rounded-xl" v-for="i in 2" :key="i" />
      </div>
    </template>
    <template v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
        <div class="flex flex-col gap-3" v-for="user in users" :key="user.name">
          <div class="font-semibold">{{ user.title }}</div>
          <div
            class="rounded-xl bg-white/50 hover:bg-white/40 flex-1 flex flex-row items-center gap-2 p-3 cursor-pointer"
            @click="
              () => {
                handleUserClick(user);
              }
            "
          >
            <div class="size-10 rounded-full bg-white/50">
              <NuxtImg
                :src="genUserAvatar(user)"
                class="size-full rounded-full"
              />
            </div>
            <div class="text-md font-semibold">{{ user.name }}</div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
