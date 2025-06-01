<script setup lang="ts">
import { useQuery } from "@tanstack/vue-query";
import { AsyncButton, Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { filesize } from "filesize";
import useAppShare from "~/composables/useShare";
const props = defineProps<{
  data: { file: File; config: any; handle_type: string; file_id: string };
}>();

const { data } = useQuery({
  queryKey: ["create-image-compress", props?.data?.file_id],
  queryFn: async () => {
    const { file_id } = props?.data || {};
    const data = await $fetch<{
      code: number;
      data: {
        id?: string;
      };
    }>(`/api/image/compress`, {
      method: "POST",
      body: {
        file_id,
      },
    });
    return data?.data;
  },
  staleTime: Infinity,
});

const taskId = computed(() => data?.value?.id);

const { data: taskData, refetch } = useQuery({
  queryKey: ["image-compress-task", taskId],
  queryFn: async () => {
    const data = await $fetch<{
      code: number;
      data: {
        result: {
          old_file: {
            id: string;
            size: number;
          };
          new_file: {
            id: string;
            size: number;
          };
        }[];
        status: "success" | "processing" | "failed";
      };
    }>(`/api/image/compress/${taskId.value}`);
    return data?.data;
  },
  enabled: !!taskId.value,
});

const { downloadFile, createFileShare } = useAppShare();

const { counter, pause } = useInterval(2000, { controls: true });

watch(
  () => counter.value,
  () => {
    if (taskData.value?.status === "success") {
      pause();
      return;
    }
    refetch();
  },
);
</script>
<template>
  <div class="flex flex-col gap-3">
    <h2 class="text-lg">上传成功</h2>
    <div class="flex flex-col gap-1 items-center">
      <div class="flex flex-col h-30 items-center justify-center">
        <FilePreviewView :value="props?.data?.file" />
      </div>
    </div>
    <div
      v-if="taskData?.status === 'success'"
      class="flex flex-col gap-2"
      v-for="item in taskData?.result"
    >
      <div
        class="bg-white/80 p-2 rounded-md w-full flex flex-row items-center justify-between"
      >
        <div class="flex flex-row gap-2 items-center">
          <div
            class="flex flex-row items-center justify-center rounded-md bg-black/5 p-2"
          >
            <LucideImage />
          </div>
          {{ props?.data?.file?.name }}
          <div class="flex flex-row gap-2 items-center text-sm">
            <span class="opacity-75">{{
              filesize(item.new_file.size ?? 0)
            }}</span>
            <span
              class="bg-green-200 text-green-600 rounded-md px-1 py-0.5 flex flex-row gap-1 items-center text-xs"
            >
              <LucideChevronDown class="size-4" />
              {{
                ((1 - item.new_file.size / item.old_file.size) * 100).toFixed(
                  2,
                )
              }}%
            </span>
          </div>
        </div>
        <AsyncButton
          variant="outline"
          class="bg-black/5"
          size="icon"
          @click="
            async () => {
              const data = await createFileShare({
                file_id: item.new_file.id,
                config: {
                  download_nums: 1,
                  expire_time: 60,
                  has_pickup_code: false,
                  has_password: false,
                },
                file_name: props?.data?.file?.name,
              });
              const { id } = data?.data || {};
              if (!id) {
                return;
              }
              await downloadFile(id);
            }
          "
        >
          <LucideDownload />
        </AsyncButton>
      </div>
    </div>
    <div v-else class="flex flex-col gap-2">
      <Skeleton
        class="w-full h-16 flex flex-row items-center justify-between"
        v-for="i in 3"
      />
    </div>
  </div>
</template>
