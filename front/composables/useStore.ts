import { defineStore } from 'pinia'
import type { DrawerItem } from '@/components/GlobalDrawer.vue'

type StoreProps = {
    tiptapCommandsView?: any
    drawer?: DrawerItem<any>[]
}

const useStore = defineStore<any, StoreProps>('store', () => {
    return {
        tiptapCommandsView: null,
        drawer: [],
    }
})

export default useStore
