interface DrawerProps {
    render: (props: { hide: () => void }) => Component
}

const showDrawer = (props: DrawerProps) => {
    const key = Math.random().toString(36).slice(2, 8)
    return new Promise<void>((res) => {
        const { render } = props || {}
        const onClose = () => {
            store._set(
                'drawer',
                (store._get('drawer')?.value ?? [])?.filter((item: any) => item.key !== key)
            )
            res()
        }
        const store = useStore()
        store._set('drawer', [...(store._get('drawer')?.value || []), { render, onClose, key }])
    })
}

export default showDrawer
