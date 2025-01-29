package events

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
	"gopkg.in/yaml.v3"
)

// convertToGoWeekday преобразует день недели из российского формата (1-7) в формат Go (0-6)
func convertToGoWeekday(russianWeekday int) int {
	if russianWeekday == 7 {
		return 0 // воскресенье
	}
	return russianWeekday
}

// convertFromGoWeekday преобразует день недели из формата Go (0-6) в российский формат (1-7)
func convertFromGoWeekday(goWeekday int) int {
	if goWeekday == 0 {
		return 7 // воскресенье
	}
	return goWeekday
}

type EventRule struct {
	DayOfWeek   *int `yaml:"day_of_week,omitempty"`   // 1-7 (понедельник = 1, воскресенье = 7)
	DayOfMonth  *int `yaml:"day_of_month,omitempty"`  // 1-31
	Month       *int `yaml:"month,omitempty"`         // 1-12
	DayOfYear   *int `yaml:"day_of_year,omitempty"`   // 1-366
	WeekOfMonth *int `yaml:"week_of_month,omitempty"` // 1-5 (5 = последняя неделя)
}

type Event struct {
	Event string    `yaml:"event"`
	Rule  EventRule `yaml:"rule"`
}

func Load(filename string) ([]Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := yaml.Unmarshal(data, &events); err != nil {
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

func (s *Scheduler) isEventToday(event Event, now time.Time) bool {
	rule := event.Rule

	// Проверяем день года (например, для дня программиста)
	if rule.DayOfYear != nil {
		currentDayOfYear := now.YearDay()
		return currentDayOfYear == *rule.DayOfYear
	}

	// Проверяем месяц
	if rule.Month != nil && int(now.Month()) != *rule.Month {
		return false
	}

	// Проверяем день месяца
	if rule.DayOfMonth != nil && now.Day() != *rule.DayOfMonth {
		return false
	}

	// Проверяем день недели (конвертируем из Go в российский формат)
	if rule.DayOfWeek != nil {
		currentWeekday := convertFromGoWeekday(int(now.Weekday()))
		if currentWeekday != *rule.DayOfWeek {
			return false
		}
	}

	// Проверяем неделю месяца
	if rule.WeekOfMonth != nil {
		weekNum := (now.Day()-1)/7 + 1
		isLastWeek := now.Day() >= 22 // Примерное определение последней недели

		if *rule.WeekOfMonth == 5 && !isLastWeek {
			return false
		} else if *rule.WeekOfMonth != 5 && weekNum != *rule.WeekOfMonth {
			return false
		}
	}

	return true
}

func (s *Scheduler) getNextEventDate(event Event, now time.Time) time.Time {
	current := now

	// Максимум 366 дней вперед (для високосного года)
	for i := 0; i < 366; i++ {
		if s.isEventToday(event, current) {
			return current
		}
		current = current.AddDate(0, 0, 1)
	}

	return now // Если что-то пошло не так, вернем текущую дату
}

func (s *Scheduler) checkEvents(now time.Time) {
	todayDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	log.Printf("=== Проверка событий [%s] ===", now.Format("2006-01-02 15:04:05"))

	hasEvents := false
	for _, event := range s.events {
		eventDate := s.getNextEventDate(event, todayDate)
		daysUntilEvent := int(eventDate.Sub(todayDate).Hours() / 24)

		// Логируем только ближайшие события (в пределах недели)
		if daysUntilEvent <= 7 {
			hasEvents = true
			log.Printf("🗓 %s: %s (через %d дн.)",
				eventDate.Format("02.01.2006"),
				event.Event,
				daysUntilEvent)
		}

		var message string
		if daysUntilEvent == 2 {
			message = fmt.Sprintf("🤔 Скоро %s 👀", event.Event)
		} else if daysUntilEvent == 0 {
			message = fmt.Sprintf("😏 Сегодня %s 😉", event.Event)
		}

		if message != "" {
			if err := s.bot.SendMessage(s.config.ChatID, message); err != nil {
				log.Printf("❌ Ошибка при отправке сообщения: %v", err)
			} else {
				log.Printf("✅ Отправлено уведомление: %s", message)
			}
		}
	}

	if !hasEvents {
		log.Printf("В ближайшие 7 дней событий нет")
	}
	log.Printf("=== Проверка завершена ===\n")
}
