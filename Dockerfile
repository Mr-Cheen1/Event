# Сборка приложения
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp main.go

# Запускаемый контейнер
FROM alpine:3.20
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .
RUN chmod +x /app/myapp

EXPOSE 8080
ENTRYPOINT ["/app/myapp"] 