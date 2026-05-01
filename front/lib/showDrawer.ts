import type { DrawerItem } from '~/components/GlobalDrawer.vue'
import asyncWait from './asyncWait'

type DrawerProps<T = unknown> = Pick<DrawerItem<T>, 'render'>

function showDrawer<T = unknown>(props: DrawerProps<T>): Promise<T | undefined> {
    const key = Math.random().toString(36).slice(2, 8)
    return new Promise((res) => {
        const { render } = props || {}
        const store = useStore()
        const onClose = async (data?: T) => {
            store.drawer = store.drawer?.map((d) => (d.key === key ? { ...d, visible: false } : d))
            await asyncWait(500)
            store.drawer = (store.drawer ?? []).filter((d) => d.key !== key)
            res(data)
        }
        store.drawer = [...(store.drawer || []), { render, onClose, key, visible: true }]
    })
}

export default showDrawer
