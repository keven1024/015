interface DrawerProps {
  render: (props: { hide: () => void }) => Component
}

const showDrawer = (props: DrawerProps) => {
    const { render } = props || {}
    const store = useStore()
    store?._set('drawer', [...(store?._get('drawer')?.value || []), { render }])
}

export default showDrawer
