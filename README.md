# 📅 Event Bot
[![Linting](https://github.com/Mr-Cheen1/Event/actions/workflows/lint.yml/badge.svg)](https://github.com/Mr-Cheen1/Event/actions/workflows/lint.yml) [![Testing](https://github.com/Mr-Cheen1/Event/actions/workflows/test.yml/badge.svg)](https://github.com/Mr-Cheen1/Event/actions/workflows/test.yml) [![Build binary](https://github.com/Mr-Cheen1/Event/actions/workflows/build.yml/badge.svg)](https://github.com/Mr-Cheen1/Event/actions/workflows/build.yml)

Telegram бот для отслеживания важных дат и событий с поддержкой различных типов праздников, написанный на Go.

## 🚀 Возможности

* 📆 Поддержка различных типов событий:
  * ✨ Фиксированные даты (например, Новый год)
  * 🔄 N-ый день недели месяца (например, третье воскресенье июля)
  * 📍 Последний определенный день месяца (например, последнее воскресенье августа)
  * 🔢 N-ый день года (например, 256-й день года - День программиста)
  * 📅 Еженедельные события
  * 📌 Ежемесячные события
* ⏰ Уведомления за 2 дня до события и в день события
* 🌍 Поддержка часовых поясов
* ⚙️ Настраиваемое время отправки уведомлений
* 📝 Конфигурация через YAML файл
* 🔍 Подробное логирование событий

## 📥 Установка

### Требования

* Go 1.21 или выше
* Telegram Bot Token
* Chat ID для отправки уведомлений

### Сборка из исходного кода

1. Клонируйте репозиторий:

```bash
git clone https://github.com/Mr-Cheen1/Event.git
cd Event
```

2. Установите зависимости:

```bash
go mod download
```

3. Создайте необходимые конфигурационные файлы:

`.env`:

```env
BOT_TOKEN=your_telegram_bot_token
CHAT_ID=your_chat_id
```

`config.json`:

```json
{
    "notification_time": "10:00",
    "timezone": "Europe/Moscow"
}
```

`events.yml`:

```yaml
# День металлурга (третье воскресенье июля)
- event: День металлурга
  rule:
    day_of_week: 7    # воскресенье
    week_of_month: 3  # третья неделя
    month: 7         # июль

# День программиста (256-й день года)
- event: День программиста
  rule:
    day_of_year: 256  # 256-й день года
```

## 🏗️ Структура проекта

```
Event/
├── 📁 bot/                # Реализация Telegram бота
│   ├── bot.go            # Основной код бота
│   └── bot_test.go       # Тесты бота
├── 📁 config/            # Конфигурация приложения
│   ├── config.go         # Загрузка конфигурации
│   └── config_test.go    # Тесты конфигурации
├── 📁 events/            # Логика обработки событий
│   ├── events.go         # Основная логика событий
│   └── events_test.go    # Тесты событий
├── 📝 main.go            # Точка входа в приложение
├── 📋 events.yml         # Список событий
├── ⚙️ config.json        # Настройки приложения
├── 🔑 .env               # Переменные окружения
└── 📖 README.md          # Документация
```

## 🔧 Как это работает

### Основные компоненты

1. **События (events/)**

   * `Event`: структура для хранения информации о событии
   * `EventRule`: правила определения даты события
   * Включает методы проверки и определения дат
2. **Конфигурация (config/)**

   * Загрузка настроек из JSON и переменных окружения
   * Управление часовым поясом и временем уведомлений
3. **Telegram бот (bot/)**

   * Отправка уведомлений в чат
   * Обработка ошибок коммуникации

### Конфигурация событий (events.yml)

```yaml
# День металлурга (третье воскресенье июля)
- event: День металлурга
  rule:
    day_of_week: 7    # воскресенье
    week_of_month: 3  # третья неделя
    month: 7         # июль

# День программиста (256-й день года)
- event: День программиста
  rule:
    day_of_year: 256  # 256-й день года
```

### Процесс работы

1. При запуске:

   * 📥 Загружает конфигурацию и список событий
   * 🤖 Инициализирует Telegram бота
   * 🕒 Запускает планировщик
2. Планировщик:

   * ⏰ Проверяет события каждый день в заданное время
   * 📊 Определяет ближайшие события
   * 📨 Отправляет уведомления при необходимости
3. Логирование:

   * 📝 Записывает все проверки событий
   * ✅ Подтверждает отправку уведомлений
   * ❌ Фиксирует ошибки

## 🧪 Тестирование

Проект включает модульные тесты для всех компонентов:

```bash
go test ./... -v
```

## 🐳 Развертывание с Docker на Timeweb Cloud

### Требования

- Аккаунт в Docker Hub
- VPS на Timeweb Cloud (Ubuntu 22.04)

### Шаг 1: Создание Docker-образа

1. Создайте `Dockerfile` в корне проекта:

```
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.json .
COPY --from=builder /app/events.yml .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

2. Соберите Docker-образ:

```bash
docker build -t event-bot .
```

3. Протестируйте локально:

```bash
docker run -d --name event-bot event-bot
```

### Шаг 2: Публикация образа в Docker Hub

1. Авторизуйтесь в Docker Hub:

```bash
docker login
```

2. Создайте тег для образа:

```bash
docker tag event-bot ваш_пользователь/event-bot:latest
```

3. Загрузите образ:

```bash
docker push ваш_пользователь/event-bot:latest
```

### Шаг 3: Настройка VPS на Timeweb Cloud

1. Создайте VPS на Timeweb Cloud (Ubuntu 22.04)
2. Подключитесь к серверу по SSH:

```bash
ssh root@ваш_ip_адрес
```

3. Обновите систему:

```bash
apt update
```

4. Установите Docker:

```bash
apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt update
apt install -y docker-ce
```

5. Настройте автозапуск Docker:

```bash
systemctl enable docker
systemctl start docker
```

### Шаг 4: Запуск бота на сервере

1. Скачайте образ с Docker Hub:

```bash
docker pull ваш_пользователь/event-bot:latest
```

2. Запустите контейнер с автоперезапуском:

```bash
docker run -d --name event-bot --restart always ваш_пользователь/event-bot:latest
```

3. Проверьте статус контейнера:

```bash
docker ps
```

4. Просмотрите логи:

```bash
docker logs event-bot
```

### Если нужно сохранение данных между перезапусками

Если ваш бот хранит состояние на диске, используйте тома Docker:

```bash
docker run -d --name event-bot \
  -v event-data:/app/data \
  --restart always \
  ваш_пользователь/event-bot:latest
```

### Обновление бота

Для обновления бота:

1. Остановите и удалите существующий контейнер:

```bash
docker stop event-bot
docker rm event-bot
```

2. Загрузите последний образ:

```bash
docker pull ваш_пользователь/event-bot:latest
```

3. Запустите новую версию бота:

```bash
docker run -d --name event-bot --restart always ваш_пользователь/event-bot:latest
```

### Мониторинг и поддержка

- Проверка статуса: `docker ps`
- Просмотр логов: `docker logs event-bot`
- Просмотр логов в реальном времени: `docker logs -f event-bot`
- Перезапуск бота: `docker restart event-bot`

При возникновении проблем Docker автоматически перезапустит бот благодаря флагу `--restart always`.

## 📝 Лицензия

Copyright © 2025 Mr-Cheen1
