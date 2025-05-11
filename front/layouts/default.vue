<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
</script>
<template>
  <div class="h-screen w-screen ">
    <GlobalDrawer />
    <img class="w-full h-full object-cover absolute inset-0 -z-[1]" src="https://fuwari.vercel.app/_astro/demo-banner.DFyx781H_Z1gN7UP.webp" />
    <Toaster position="top-center" richColors closeButton  />
    <div class=" h-full w-full flex flex-col items-center lg:p-10 p-5">
      <Navbar />
      <slot />
    </div>

  </div>
</template>