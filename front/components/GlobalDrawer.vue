<script setup lang="ts">
import {
  Drawer,
  DrawerContent,
} from '@/components/ui/drawer'
import { createVNode } from 'vue'
import useStore from '@/composables/useStore'

const store = useStore('drawer')
const drawer = computed(() => store?._get('drawer'))
const currentDrawer = computed(() => drawer?.value?.[drawer?.value?.length - 1])

const render = computed(() => currentDrawer?.value?.render)
const Children = () =>
  createVNode(render.value, {
    hide: () => store?._set('drawer', drawer?.value?.slice(0, -1)),
  })

</script>

<template>
  <Drawer :open="!!store?._get('drawer')?.[store?._get('drawer')?.length - 1]" @update:open="(open) => {
    if (!open && drawer?.length > 0) {
      store?._set('drawer', drawer?.slice(0, -1))
    }
  }">
    <DrawerContent>
      <div class="mx-auto w-full max-w-sm pb-10">
        <Children />
      </div>
    </DrawerContent>
  </Drawer>
</template>