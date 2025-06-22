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
            link: [
                { rel: 'icon', href: '/logo.png', sizes: 'any' },
                // { rel: 'icon', href: '/favicon.svg', sizes: 'any', type: 'image/svg+xml' },
                { rel: 'apple-touch-icon', sizes: '180x180', href: '/logo.png' },
            ],
            meta: [
                // used on some mobile browsers
                { name: 'theme-color', content: '#395276' },
            ],
            ...head,
            title: title ? `${title} - ${seoMeta?.value?.site_title}` : seoMeta?.value?.site_title,
        })
        useSeoMeta({
            ...seo,
            title: seoMeta?.value?.site_title,
            description: seoMeta?.value?.site_desc,
            ogTitle: seoMeta?.value?.site_title,
            ogDescription: seoMeta?.value?.site_desc,
            ogImage: {
                url: `${seoMeta?.value?.site_url}/logo.png`,
                width: 1024,
                height: 1024,
                alt: 'logo',
                type: 'image/png',
            },
            twitterCard: 'summary',
        })
    }
    return
}

export default useSeo
