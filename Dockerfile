# Этап сборки
FROM golang:1.21-alpine AS builder

# Установка необходимых инструментов для сборки
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Финальный этап
FROM alpine:latest

# Установка и настройка часовых поясов
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/Samara /etc/localtime && \
    echo "Europe/Samara" > /etc/timezone

# Создание рабочей директории
WORKDIR /app

# Копирование собранного приложения из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .

# Даем права на выполнение
RUN chmod +x /app/main

# Указываем порт, который будет использовать приложение
EXPOSE 8080

# Устанавливаем переменные окружения
ENV BOT_TOKEN=""
ENV CHAT_ID=""

# Запускаем приложение с выводом всей информации
CMD echo "=== Environment ===" && \
    env | grep -E "BOT_TOKEN|CHAT_ID" && \
    echo "=== Files ===" && \
    ls -la && \
    echo "=== Starting ===" && \
    /app/main 