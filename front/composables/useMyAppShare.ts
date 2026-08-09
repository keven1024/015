let shareIdTokenMap: WeakMap<{ share_id: string }, string>

const getShareToken = async (
    share_id: string,
    options?: {
        password?: string
    }
): Promise<string | undefined> => {
    if (!shareIdTokenMap) {
        shareIdTokenMap = new WeakMap()
    }
    let token = shareIdTokenMap.get({ share_id })
    const { password } = options || {}
    if (!token) {
        const data = await $fetch<{
            code: number
            message: string
            data: {
                token?: string
            }
        }>(`/api/download`, {
            method: 'POST',
            body: {
                share_id,
                password,
            },
        })
        if (!data?.data?.token) {
            throw new Error(data?.message || '获取token失败')
        }
        token = data.data.token
        shareIdTokenMap.set({ share_id }, token)
    }
    return token
}

const downloadFile = (token: string) => {
    const a = document.createElement('a')
    a.href = `/api/download?token=${token}`
    a.download = ''
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
}

const downloadFileByShareId = async (share_id: string) => {
    const token = await getShareToken(share_id)
    if (!token) {
        throw new Error('获取token失败')
    }
    return downloadFile(token)
}

const createShare = async (data: any) => {
    return await $fetch<{
        code: number
        data: {
            id?: string
            download_nums?: number
            expire_at?: number
            files?: { id?: string; file_name: string }[]
            pickup_code?: string
        }
    }>(`/api/share`, {
        method: 'POST',
        body: data,
    })
}

const createFileShare = async (data: {
    files: { id: string; file_name: string }[]
    config: {
        download_nums: number
        expire_time: number
        has_pickup_code?: boolean
        has_password?: boolean
        pickup_code?: string
        password?: string
        notify_email?: string
    }
}) => {
    const { files, config } = data || {}
    return await createShare({ type: 'file', files, config })
}

const createTextShare = async (data: { text: string; config: any }) => {
    const { text, config } = data || {}
    return await createShare({
        type: 'text',
        text,
        config,
    })
}

const useMyAppShare = () => {
    return {
        downloadFile,
        downloadFileByShareId,
        createShare,
        createFileShare,
        createTextShare,
        getShareToken,
    }
}

export default useMyAppShare
