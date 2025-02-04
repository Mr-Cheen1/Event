# Сборка приложения
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp main.go

# Запускаемый контейнер
FROM alpine:latest
WORKDIR /app
COPY --from=builder /myapp /app/myapp
COPY config.json /app/config.json

EXPOSE 8080
ENTRYPOINT ["/app/myapp"] 