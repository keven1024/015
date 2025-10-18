<script setup lang="ts">
import { Drawer, DrawerContent } from '@/components/ui/drawer'
import { createVNode } from 'vue'
import useStore from '@/composables/useStore'

const store = useStore()
const currentDrawer = computed(() => store.drawer?.[store.drawer?.length - 1])

const render = computed(() => currentDrawer?.value?.render)
const hide = computed(() => currentDrawer?.value?.onClose)
const Children = () =>
    render.value
        ? createVNode(render.value, {
              hide: hide?.value,
          })
        : null
</script>

<template>
    <Drawer
        :open="!!store.drawer?.[store.drawer?.length - 1]"
        @update:open="
            (open) => {
                if (!open && store?.drawer?.length && hide) {
                    hide()
                }
            }
        "
    >
        <DrawerContent>
            <div class="mx-auto w-full max-w-lg pb-10 px-3">
                <Children />
            </div>
        </DrawerContent>
    </Drawer>
</template>
