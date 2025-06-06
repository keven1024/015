<script lang="ts" setup>
import { Toaster } from "vue-sonner";
await useSeo();
</script>
<template>
  <div class="h-screen w-screen">
    <GlobalDrawer />
    <Toaster position="top-center" richColors closeButton />
    <img
      class="w-full h-full object-cover absolute inset-0 -z-[1]"
      src="https://img.fudaoyuan.icu/api/1/random/?scale_min=1.5&webp=true&md=false&format=302"
    />
    <div
      class="h-full w-full flex flex-col items-center lg:p-10 p-5 overflow-y-auto"
    >
      <Navbar />
      <slot />
    </div>
  </div>
</template>
