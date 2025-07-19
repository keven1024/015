import { defineStore } from 'pinia'

type renderComponent = Component | ((props: { hide: () => null }) => VNode)

type StoreProps = {
    tiptapCommandsView?: any
    drawer?: { render?: renderComponent; onClose: (data?: any) => void; key: string }[]
}

const useStore = defineStore<any, StoreProps>('store', () => {
    return {
        tiptapCommandsView: null,
        drawer: [],
    }
})

export default useStore
