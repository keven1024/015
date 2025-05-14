type UseSeoProps = {
    head?: Record<string, any>
    seo?: Record<string, any>
}
const useSeo = async (props: UseSeoProps = {}) => {
    const { head, seo } = props || {}
    const { data } = await useFetch<any>('/config')
    const seoMeta = computed(() => data.value?.data)
    const { title } = head || {}
    useHead({
        ...head,
        title: title ? `${title} - ${seoMeta?.value?.site_title}` : seoMeta?.value?.site_title,
    })
    useSeoMeta({
        ...seo,
        title: seoMeta?.value?.site_title,
        description: seoMeta?.value?.site_desc,
        ogTitle: seoMeta?.value?.site_title,
        ogDescription: seoMeta?.value?.site_desc,
        // ogImage: seoMeta?.value?.site_url,
        // twitterCard: 'summary_large_image',
    })
    return
}

export default useSeo