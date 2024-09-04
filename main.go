package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Mr-Cheen1/Event/config"
	"github.com/Mr-Cheen1/Event/events"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const (
	configFile = "config.json"
	eventsFile = "events.json"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании chatID: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	config := config.LoadConfig(configFile)
	eventsList := events.LoadEvents(eventsFile)

	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Fatalf("Ошибка при загрузке часового пояса: %v", err)
	}

	for {
		now := time.Now().In(loc)
		notificationTime, err := time.ParseInLocation("15:04", config.NotificationTime, loc)
		if err != nil {
			log.Printf("Ошибка при парсинге времени уведомления: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		notificationTime = time.Date(
			now.Year(), now.Month(), now.Day(),
			notificationTime.Hour(), notificationTime.Minute(),
			0, 0, loc,
		)

		if now.Equal(notificationTime) || now.After(notificationTime) && now.Before(notificationTime.Add(1*time.Minute)) {
			events.CheckEvents(bot, eventsList, now, chatID)
			time.Sleep(24 * time.Hour)
		} else {
			time.Sleep(time.Until(notificationTime))
		}
	}
}
