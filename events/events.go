package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
)

type Event struct {
	Date  string `json:"date"`
	Event string `json:"event"`
}

func Load(filename string) ([]Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, err
	}

	return events, nil
}

type Scheduler struct {
	config *config.Config
	bot    bot.API
	events []Event
}

func NewScheduler(cfg *config.Config, bot bot.API, events []Event) *Scheduler {
	return &Scheduler{
		config: cfg,
		bot:    bot,
		events: events,
	}
}

func (s *Scheduler) Start() {
	loc, err := time.LoadLocation(s.config.Timezone)
	if err != nil {
		log.Fatalf("Ошибка при загрузке часового пояса: %v", err)
	}

	for {
		now := time.Now().In(loc)
		notificationTime, err := time.ParseInLocation("15:04", s.config.NotificationTime, loc)
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
			s.checkEvents(now)
			time.Sleep(24 * time.Hour)
		} else {
			time.Sleep(time.Until(notificationTime))
		}
	}
}

func (s *Scheduler) checkEvents(now time.Time) {
	for _, event := range s.events {
		eventDate, err := time.Parse("02.01", event.Date)
		if err != nil {
			log.Printf("Ошибка при парсинге даты события: %v", err)
			continue
		}
		eventDate = time.Date(now.Year(), eventDate.Month(), eventDate.Day(), 0, 0, 0, 0, now.Location())

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
			if err := s.bot.SendMessage(s.config.ChatID, message); err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
			}
		}
	}
}
