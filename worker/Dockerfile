# syntax=docker/dockerfile:1

FROM golang:1.23.1 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker

FROM alpine:3
RUN apk add --no-cache pngquant jpegoptim
COPY --from=builder /app/worker /bin/worker
ENTRYPOINT ["/bin/worker"]