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

# Установка пакета часовых поясов
RUN apk add --no-cache tzdata

# Настройка часового пояса
ENV TZ=Europe/Moscow

# Создание рабочей директории
WORKDIR /app

# Копирование собранного приложения из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .
COPY --from=builder /app/.env .

# Указываем порт, который будет использовать приложение
EXPOSE 8080

# Запуск приложения
CMD ["./main"] 