<script setup lang="ts">
import { Drawer, DrawerContent } from '@/components/ui/drawer'
import { createVNode } from 'vue'
import useStore from '@/composables/useStore'

const store = useStore('drawer')
const drawer = computed(() => store?._get('drawer'))
const currentDrawer = computed(() => drawer?.value?.[drawer?.value?.length - 1])

const render = computed<() => Component>(() => currentDrawer?.value?.render)
const hide = computed<() => void>(() => currentDrawer?.value?.onClose)
const Children = () =>
    createVNode(render.value, {
        hide: hide?.value,
    })
</script>

<template>
    <Drawer
        :open="!!drawer?.[drawer?.length - 1]"
        @update:open="
            (open) => {
                if (!open && drawer?.length > 0) {
                    hide()
                }
            }
        "
    >
        <DrawerContent>
            <div class="mx-auto w-full max-w-sm pb-10">
                <Children />
            </div>
        </DrawerContent>
    </Drawer>
</template>
