import getApiBaseUrl from '~/lib/getApiBaseUrl'

type UseSeoProps = {
    head?: Record<string, any>
    seo?: Record<string, any>
}
const useSeo = async (props: UseSeoProps = {}) => {
    const { head, seo } = props || {}
    const seoMeta = ref<{
        site_title: string
        site_desc: string
        site_url: string
        site_icon: string
        site_bg_url: string
    }>()
    if (import.meta.server) {
        await fetch(`${getApiBaseUrl()}/config`)
            .then((res) => res.json())
            .then(({ data }) => {
                seoMeta.value = data
            })
        const { title } = head || {}
        useHead({
            link: [
                { rel: 'icon', href: seoMeta.value?.site_icon || '/logo.png', sizes: 'any' },
                // { rel: 'icon', href: '/favicon.svg', sizes: 'any', type: 'image/svg+xml' },
                { rel: 'apple-touch-icon', sizes: '180x180', href: seoMeta.value?.site_icon || '/logo.png' },
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
                url: `${seoMeta?.value?.site_url}${seoMeta?.value?.site_icon || '/logo.png'}`,
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
