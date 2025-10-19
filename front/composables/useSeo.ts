import getApiBaseUrl from '~/lib/getApiBaseUrl'
import renderI18n from '~/lib/renderI18n'

type UseSeoProps = {
    head?: Record<string, any>
    seo?: Record<string, any>
    locale?: string
}
const useSeo = async (props: UseSeoProps = {}) => {
    const { head, seo, locale } = props || {}
    const seoMeta = ref<{
        site_title: Record<string, string>
        site_desc: Record<string, string>
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
        const siteTitle = computed(() => renderI18n(seoMeta?.value?.site_title || {}, 'en', locale))
        const siteDesc = computed(() => renderI18n(seoMeta?.value?.site_desc || {}, 'en', locale))
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
            title: title ? `${title} - ${siteTitle.value}` : siteTitle.value,
        })
        useSeoMeta({
            ...seo,
            title: siteTitle.value,
            description: siteDesc.value,
            ogTitle: siteTitle.value,
            ogDescription: siteDesc.value,
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
