FROM node:22-alpine AS front-base
WORKDIR /app

# Install dependencies only when needed
FROM front-base AS front-builder
RUN apk add --no-cache gcompat
ENV CI=true
COPY . .
RUN corepack enable pnpm && pnpm i && pnpm --filter=015-front build && pnpm --dir pkg/mail export

FROM golang:1.26.3 AS backend-builder
WORKDIR /app
# Workspace and module manifests for cache
COPY go.work go.work.sum ./
COPY backend/ ./backend/
COPY worker/ ./worker/
COPY pkg/ ./pkg/
# Inject built email templates so Go can embed them
COPY --from=front-builder /app/pkg/mail/out/ ./pkg/mail/out/
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download
# Build from workspace root so pkg/utils, pkg/models, pkg/services resolve
RUN CGO_ENABLED=0 GOOS=linux go build -o backend-bin ./backend


FROM front-base AS runner
ARG VERSION
ARG BUILD_TIME
RUN apk add --no-cache curl openssl
ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Only `.output` folder is needed from the build stage
COPY --from=front-builder --chown=nuxtjs:nodejs /app/front/.output/ ./
COPY --from=backend-builder /app/backend-bin /bin/backend
COPY 015.sh /app/015.sh

# Change the port and host
ENV PORT=80 HOST=0.0.0.0
ENV VERSION=${VERSION}
ENV BUILD_TIME=${BUILD_TIME}

EXPOSE 80

CMD ["/bin/sh", "/app/015.sh"]