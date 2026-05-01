<script setup lang="ts">
import { Drawer, DrawerContent } from '@/components/ui/drawer'
import type { VNode } from 'vue'
import useStore from '@/composables/useStore'
import { isFunction } from 'lodash-es'

type DrawerOnclose<T = unknown> = (data?: T) => Promise<void>
type DrawerRender<T = unknown> = VNode | ((props: { hide: DrawerOnclose<T> }) => VNode)
export type DrawerItem<T = unknown> = {
    render?: DrawerRender<T>
    onClose: DrawerOnclose<T>
    key: string
    visible: boolean
}

const store = useStore()

const Children = ({ drawer }: { drawer: DrawerItem }) => {
    if (!drawer.render) {
        return null
    }
    return isFunction(drawer.render) ? drawer.render({ hide: drawer.onClose }) : drawer.render
}
</script>

<template>
    <Drawer
        v-for="item in store.drawer"
        :key="item.key"
        :open="item.visible"
        @update:open="
            (open) => {
                if (!open) {
                    item.onClose()
                }
            }
        "
    >
        <DrawerContent>
            <div class="mx-auto min-w-lg max-w-[80vw] pb-10 px-3">
                <Children :drawer="item" />
            </div>
        </DrawerContent>
    </Drawer>
</template>
