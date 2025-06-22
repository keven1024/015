type UseSeoProps = {
    head?: Record<string, any>
    seo?: Record<string, any>
}
const useSeo = async (props: UseSeoProps = {}) => {
    const { head, seo } = props || {}
    const seoMeta = ref<any>()
    if (import.meta.server) {
        const { SITE_TITLE, SITE_DESC, SITE_URL } = process.env || {}
        seoMeta.value = {
            site_title: SITE_TITLE,
            site_desc: SITE_DESC,
            site_url: SITE_URL,
        }
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
    }
    return
}

export default useSeo
