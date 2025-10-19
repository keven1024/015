const useMyAppConfig = () => {
    const { data } = useFetch<{
        data: {
            site_title: Record<string, string>
            site_desc: Record<string, string>
            site_url: string
            site_icon: string
            site_bg_url: string
        }
    }>('/api/config')
    return computed(() => data?.value?.data)
}

export default useMyAppConfig
