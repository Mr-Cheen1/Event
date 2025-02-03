# Этап сборки
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Финальный этап
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/Samara /etc/localtime
EXPOSE 8080
CMD /app/main 