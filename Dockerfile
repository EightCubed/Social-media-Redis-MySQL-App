# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o social-media-app ./cmd/server/main.go

# Stage 2: Run
FROM alpine:latest AS runner

WORKDIR /root/

RUN apk add --no-cache mysql-client ca-certificates

COPY --from=builder /app/social-media-app .

RUN chmod +x social-media-app

EXPOSE 8080

ENV DB_HOST="mysql" \
    DB_PORT="3306" \
    DB_NAME="social_media_app" \
    REDIS_HOST="redis" \
    REDIS_PORT="6379"

CMD ["./social-media-app"]