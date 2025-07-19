interface DrawerProps {
    render: (props: { hide: () => void }) => Component
}

const showDrawer = (props: DrawerProps) => {
    const key = Math.random().toString(36).slice(2, 8)
    return new Promise<void>((res) => {
        const { render } = props || {}
        const onClose = (data?: any) => {
            store.drawer = (store.drawer ?? [])?.filter((item: any) => item.key !== key)
            res(data)
        }
        const store = useStore()
        store.drawer = [...(store.drawer || []), { render, onClose, key }]
    })
}

export default showDrawer
