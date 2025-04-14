const showDrawer = (props: any) => {
    const { render } = props || {}
    const store = useStore()
    store?._set('drawer', [...(store?._get('drawer')?.value || []), { render }])
}

export default showDrawer
