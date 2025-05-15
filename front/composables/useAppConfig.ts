const useAppConfig = () => {
    const { data } = useFetch('/config')
    return data
}

export default useAppConfig