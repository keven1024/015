<script setup lang="ts">
import { Button } from "@/components/ui/button";
import FilePreviewView from "@/components/FilePreviewView.vue";
import { Input } from "@/components/ui/input";
import { useClipboard } from "@vueuse/core";
import { toast } from "vue-sonner";
import { useQuery } from "@tanstack/vue-query";
import useMyAppShare from "@/composables/useMyAppShare";
import useMyAppConfig from "@/composables/useMyAppConfig";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/zh-cn"; // 导入中文语言包
import showDrawer from "@/lib/showDrawer";
import QrCoreDrawer from "@/components/Drawer/QrCoreDrawer.vue";
dayjs.extend(relativeTime); // 扩展 relativeTime 插件
dayjs.locale("zh-cn"); // 设置语言为中文

const props = defineProps<{
  data: { file: File; config: any; handle_type: string; file_id: string };
}>();
const emit = defineEmits<{
  (e: "change", key: string): void;
}>();
const { createFileShare } = useMyAppShare();
const { data } = useQuery({
  queryKey: ["create-share", props?.data?.file_id],
  queryFn: async () => {
    const { file_id, config, file } = props?.data || {};
    const { name } = file || {};
    const data = await createFileShare({
      file_id,
      config,
      file_name: name,
    });
    return data?.data;
  },
});

const appConfig = useMyAppConfig();
const url = computed(() => {
  const { id } = data?.value || {};
  return `${appConfig?.value?.site_url}/s/${id}`;
});

const { copy } = useClipboard();
</script>

<template>
  <div class="flex flex-col gap-3">
    <h2 class="text-lg">上传成功</h2>
    <div class="flex flex-col gap-3 items-center">
      <div class="flex flex-col h-30 items-center">
        <FilePreviewView :value="props?.data?.file" />
      </div>
      <div
        class="flex flex-col md:flex-row gap-5 rounded-md p-5 bg-white/20 backdrop-blur-xl w-full"
      >
        <div class="flex flex-col gap-2 flex-1">
          <div class="text-sm font-semibold">信息</div>
          <div class="grid grid-cols-2 gap-2">
            <div class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1">
              <div class="text-xs font-semibold">下载次数</div>
              <div class="text-3xl font-light">{{ data?.download_nums }}</div>
            </div>
            <div class="rounded-xl flex flex-col bg-black/5 px-3 py-2 gap-1">
              <div class="text-xs font-semibold">过期时间</div>
              <div class="text-md font-light">
                {{
                  dayjs(data?.expire_at * 1000).format("YYYY-MM-DD HH:mm:ss")
                }}
              </div>
            </div>
            <div
              class="rounded-xl flex flex-col bg-black/10 px-3 py-2 gap-1"
              v-if="data?.pickup_code"
            >
              <div class="flex flex-row justify-between w-full items-center">
                <div class="text-xs font-semibold">提取码</div>
                <Button
                  variant="outline"
                  class="bg-white/70 p-0 size-6"
                  size="icon"
                  @click="
                    () => {
                      copy(data?.pickup_code);
                      toast.success('复制成功');
                    }
                  "
                >
                  <LucideCopy class="size-3" />
                </Button>
              </div>
              <div class="flex flex-row gap-2">
                <div v-for="s in data?.pickup_code" class="text-2xl font-light">
                  {{ s }}
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-5 flex-1">
          <div class="text-sm font-semibold">链接</div>
          <div class="flex flex-row gap-2">
            <Input v-model="url" class="bg-white/70" readonly />
            <Button
              variant="outline"
              class="bg-white/70"
              size="icon"
              @click="
                () => {
                  copy(url);
                  toast.success('复制成功');
                }
              "
            >
              <LucideCopy />
            </Button>

            <Button
              variant="outline"
              class="bg-white/70"
              size="icon"
              @click="
                () => {
                  showDrawer({
                    render: ({ ...rest }) =>
                      h(QrCoreDrawer, {
                        ...rest,
                        data: url,
                      }),
                  });
                }
              "
            >
              <LucideQrCode />
            </Button>
          </div>
        </div>
      </div>
      <Button
        class="hover:bg-white/50 w-40"
        @click="
          () => {
            emit('change', 'input');
          }
        "
      >
        确定
      </Button>
    </div>
  </div>
</template>
