# Сборка приложения
FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Запускаемый контейнер
FROM alpine:3.20
RUN apk add --no-cache curl tzdata ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .
RUN chmod +x /app/main
RUN chown -R nobody:nobody /app
USER nobody

EXPOSE 8080
CMD ["./main"] 