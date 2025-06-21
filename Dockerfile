FROM node:22-alpine AS front-base

# Install dependencies only when needed
FROM front-base AS front-deps
RUN apk add --no-cache libc6-compat 
WORKDIR /app
COPY . .
RUN corepack enable pnpm && pnpm i && pnpm --filter=015-front deploy dist


FROM front-base AS front-builder
WORKDIR /app
COPY --from=front-deps /app/dist/ .
RUN corepack enable pnpm && pnpm build

FROM golang:1.23.1 AS backend-builder
WORKDIR /app
# Download Go modules
COPY backend/go.mod backend/go.sum ./
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY backend/ .
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o backend


FROM front-base AS runner
ARG VERSION
ARG BUILD_TIME
WORKDIR /app
RUN apk add --no-cache curl openssl
ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxtjs

# Only `.output` folder is needed from the build stage
COPY --from=front-builder --chown=nuxtjs:nodejs /app/.output/ ./
COPY --from=backend-builder /app/backend /bin/backend
COPY 015.sh /app/015.sh

# Change the port and host
ENV PORT=80 HOST=0.0.0.0
ENV SITE_URL="http://localhost"
ENV UPLOAD_PATH="/uploads"
ENV VERSION=${VERSION}
ENV BUILD_TIME=${BUILD_TIME}

EXPOSE 80

CMD ["/bin/sh", "/app/015.sh"]