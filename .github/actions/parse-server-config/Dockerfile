# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o parse-server-config

FROM alpine:latest
COPY --from=builder /app/parse-server-config /parse-server-config
ENTRYPOINT ["/parse-server-config"]
