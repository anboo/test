FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY migration ./migration
ENV LISTEN_PORT=:8080
ENV MIGRATION_PATH=/app/migration
CMD ["./app"]
