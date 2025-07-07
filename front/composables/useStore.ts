import { defineStore } from 'pinia'
import { cloneDeep, get, isEmpty, isUndefined, set, isString } from 'lodash-es'

type StoreType = Record<string, any>

const initState: StoreType = {}
// 做了一点小小的改进，可以传入key，会自动初始化，如果不初始化的话容易导致不存在值而丢失响应式
const useStore = (key?: string) => {
    const store = defineStore('store', {
        state: () => ({
            ...initState,
        }),
        actions: {
            _get(path?: string) {
                if (isEmpty(path) || isUndefined(path)) {
                    return this.$state
                }
                return get(this.$state, path)
            },
            _set(path: string, value: any) {
                const newState = cloneDeep(this.$state)
                set(newState, path, value)
                this.$patch(newState)
            },
        },
    })()
    if (!isEmpty(key) && isString(key) && isUndefined(store?._get(key))) {
        // console.log('reset', key, store?._get(key))
        store?._set(key, null)
    }
    return store
}

export default useStore
