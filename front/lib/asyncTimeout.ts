import asyncWait from './asyncWait'

const asyncTimeout = <T>(fn: () => Promise<T>, ms: number) => {
    return Promise.race([
        fn(),
        async () => {
            await asyncWait(ms)
            throw new Error('timeout')
        },
    ])
}

export default asyncTimeout
