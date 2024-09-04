package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Event struct {
	Date  string `json:"date"`
	Event string `json:"event"`
}

type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func LoadEvents(filename string) []Event {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла событий: %v", err)
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		log.Fatalf("Ошибка при парсинге файла событий: %v", err)
	}

	return events
}

func CheckEvents(bot BotAPI, events []Event, now time.Time, chatID int64) {
	for _, event := range events {
		eventDate, err := time.Parse("02.01", event.Date)
		if err != nil {
			log.Printf("Ошибка при парсинге даты события: %v", err)
			continue
		}
		eventDate = time.Date(now.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, now.Location())

		// Рассчитываем количество дней до события
		todayDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		daysUntilEvent := int(eventDate.Sub(todayDate).Hours() / 24)

		log.Printf("Событие: %s, Дата: %s, Дней до события: %d",
			event.Event, eventDate.Format("02.01"), daysUntilEvent)

		var message string
		if daysUntilEvent == 2 {
			message = fmt.Sprintf("🤔 Скоро %s 👀", event.Event)
		} else if daysUntilEvent == 0 {
			message = fmt.Sprintf("😏 Сегодня %s 😉", event.Event)
		}

		if message != "" {
			msg := tgbotapi.NewMessage(chatID, message)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
			}
		}
	}
}
