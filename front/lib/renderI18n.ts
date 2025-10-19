import { useI18n } from 'vue-i18n'

const renderI18n = (json: Record<string, string>, defaultKey: string, locale?: string) => {
    const { locale: _locale } = useI18n()
    if (!json) return ''
    if (!locale) {
        locale = _locale.value
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
