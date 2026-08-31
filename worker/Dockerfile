# syntax=docker/dockerfile:1

FROM node:22-alpine AS mail-builder
RUN apk add --no-cache gcompat
ENV CI=true
WORKDIR /app
COPY . .
RUN corepack enable pnpm && pnpm i && pnpm --dir pkg/mail export:ci

FROM golang:1.25.5 AS builder

# Set destination for COPY
WORKDIR /app
COPY go.work go.work.sum ./
COPY backend/ ./backend/
COPY worker/ ./worker/
COPY pkg/ ./pkg/
# Inject built email templates so Go can embed them
COPY --from=mail-builder /app/pkg/mail/out/ ./pkg/mail/out/
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-bin ./worker

FROM alpine:3
RUN apk add --no-cache pngquant jpegoptim imagemagick
COPY --from=builder /app/worker-bin /bin/worker
ENTRYPOINT ["/bin/worker"]
