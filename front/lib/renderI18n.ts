const renderI18n = (json: Record<string, string>, defaultKey: string, locale?: string) => {
    if (!json) return ''
    if (!locale) {
        return json?.[defaultKey]
    }
    if (!json?.[locale]) {
        const [baseLocaleKey, subLocaleKey] = locale?.split('-') || []
        if (!baseLocaleKey || !subLocaleKey) {
            return json?.[defaultKey]
        }
        return renderI18n(json, defaultKey, baseLocaleKey)
    }
    return json?.[locale]
}

export default renderI18n
