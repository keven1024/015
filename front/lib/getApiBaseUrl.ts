const getApiBaseUrl = () => {
    return import.meta.env.API_BASE_URL?.replace(/\/$/, '') || 'http://127.0.0.1:5001'
}

export default getApiBaseUrl
