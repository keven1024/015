export default defineEventHandler(() => {
    const { SITE_TITLE, SITE_DESC, SITE_URL } = process.env || {}
    return {
        site_title: SITE_TITLE,
        site_desc: SITE_DESC,
        site_url: SITE_URL,
    }
})