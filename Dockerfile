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

# Установка часового пояса
RUN apk add --no-cache tzdata

# Создание рабочей директории
WORKDIR /app

# Копирование собранного приложения из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .

# Установка переменных окружения по умолчанию
ENV TZ=Europe/Moscow

# Запуск приложения
CMD ["./main"] 