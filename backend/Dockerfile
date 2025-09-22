# syntax=docker/dockerfile:1

FROM golang:1.24.3 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o backend

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
FROM alpine:3
COPY --from=builder /app/backend /bin/backend
ENTRYPOINT ["/bin/backend"]
EXPOSE 8080