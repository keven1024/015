import asyncWait from './asyncWait'

interface AsyncRetryOptions {
    retryCount?: number
    delay?: number
}

const asyncRetry = async <T>(fn: () => Promise<T>, options: AsyncRetryOptions = {}) => {
    const { retryCount = 3, delay = 1000 } = options || {}
    await asyncWait(delay)
    try {
        return await fn()
    } catch (error) {
        if (retryCount > 0) {
            return await asyncRetry(fn, { retryCount: retryCount - 1, delay: delay * 2 })
        }
        throw error
    }
}

export default asyncRetry
