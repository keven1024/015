<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useClipboard } from "@vueuse/core";
import { toast } from "vue-sonner";
import { useQuery } from "@tanstack/vue-query";
import useMyAppShare from "@/composables/useMyAppShare";
import useMyAppConfig from "@/composables/useMyAppConfig";
import showDrawer from "@/lib/showDrawer";
import QrCoreDrawer from "@/components/Drawer/QrCoreDrawer.vue";

const props = defineProps<{
  data: { text: string; config: any; handle_type: string };
}>();

const emit = defineEmits<{
  (e: "change", key: string): void;
}>();

const { createTextShare } = useMyAppShare();
const { data } = useQuery({
  queryKey: ["create-share", props?.data?.text],
  queryFn: async () => {
    const { config, text } = props?.data || {};
    const data = await createTextShare({
      text,
      config,
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
    <div class="flex flex-col md:flex-row gap-2">
      <div class="flex flex-row justify-between md:basis-1/2">
        <h2 class="text-lg">分享成功</h2>
        <Button
          variant="outline"
          class="bg-white/70"
          size="icon"
          @click="
            () => {
              emit('change', 'input');
            }
          "
        >
          <LucideHome />
        </Button>
      </div>
      <div class="flex flex-row gap-2 flex-1">
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
    <div
      class="prose rounded-md bg-white/70 p-3 w-full max-w-full min-h-[30vh]"
      v-html="props?.data?.text"
    />
  </div>
</template>
