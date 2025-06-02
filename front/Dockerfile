FROM node:22-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat 
WORKDIR /app
COPY package.json ./
RUN corepack enable pnpm && pnpm i


FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN corepack enable pnpm && pnpm i && pnpm build

FROM base AS runner
ARG BUILD_TAG
WORKDIR /app
RUN apk add --no-cache curl openssl
ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Only `.output` folder is needed from the build stage
COPY --from=builder --chown=nuxtjs:nodejs /app/.output/ ./

# Change the port and host
ENV PORT 80
ENV HOST 0.0.0.0

EXPOSE 80

CMD ["node", "/app/server/index.mjs"]