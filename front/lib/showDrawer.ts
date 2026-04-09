import type { DrawerItem } from '~/components/GlobalDrawer.vue'

type DrawerProps<T = unknown> = Pick<DrawerItem<T>, 'render'>

function showDrawer<T = unknown>(props: DrawerProps<T>): Promise<T | undefined> {
    const key = Math.random().toString(36).slice(2, 8)
    return new Promise((res) => {
        const { render } = props || {}
        const onClose = (data?: T) => {
            store.drawer = (store.drawer ?? [])?.filter((item) => item.key !== key)
            res(data)
        }
        const store = useStore()
        store.drawer = [...(store.drawer || []), { render, onClose, key }]
    })
}

export default showDrawer
