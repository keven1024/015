const useMyAppConfig = () => {
    const { data } = useFetch<{
        data: {
            site_title: Record<string, string>
            site_desc: Record<string, string>
            site_url: string
            site_icon: string
            site_bg_url: string
            version: string
            build_time: number
        }
    }>('/api/config')
    return computed(() => data?.value?.data)
}

export default useMyAppConfig
