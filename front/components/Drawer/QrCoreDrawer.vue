<script setup lang="ts">
import QRCode from "qrcode";
const props = defineProps<{
  hide: () => void;
  data: any;
}>();
const { state } = useAsyncState(async () => {
  return await QRCode.toDataURL(props.data);
}, null);
</script>

<template>
  <div class="flex flex-col gap-5">
    <div class="text-xl font-bold">分享二维码</div>
    <div class="flex flex-row justify-center">
      <img :src="state" v-if="!!state" />
      <Skeleton class="size-20" v-else />
    </div>
  </div>
</template>
