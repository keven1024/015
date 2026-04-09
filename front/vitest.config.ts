import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vitest/config'

const rootDir = fileURLToPath(new URL('./', import.meta.url))

export default defineConfig({
    resolve: {
        alias: [
            { find: /^@\/(.*)$/, replacement: `${rootDir}$1` },
            { find: /^~\/(.*)$/, replacement: `${rootDir}$1` },
            { find: /^@@\/(.*)$/, replacement: `${rootDir}$1` },
            { find: /^~~\/(.*)$/, replacement: `${rootDir}$1` },
        ],
    },
})
