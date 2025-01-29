package events

import (
	"fmt"
	"testing"
	"time"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
)

func intPtr(i int) *int {
	return &i
}

func TestWeekdayConversion(t *testing.T) {
	tests := []struct {
		russianWeekday int
		goWeekday      int
	}{
		{1, 1}, // Понедельник
		{2, 2}, // Вторник
		{3, 3}, // Среда
		{4, 4}, // Четверг
		{5, 5}, // Пятница
		{6, 6}, // Суббота
		{7, 0}, // Воскресенье
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Russian %d to Go", tt.russianWeekday), func(t *testing.T) {
			if got := convertToGoWeekday(tt.russianWeekday); got != tt.goWeekday {
				t.Errorf("convertToGoWeekday(%d) = %d, want %d", tt.russianWeekday, got, tt.goWeekday)
			}
		})

		t.Run(fmt.Sprintf("Go %d to Russian", tt.goWeekday), func(t *testing.T) {
			if got := convertFromGoWeekday(tt.goWeekday); got != tt.russianWeekday {
				t.Errorf("convertFromGoWeekday(%d) = %d, want %d", tt.goWeekday, got, tt.russianWeekday)
			}
		})
	}
}

func TestIsEventToday(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		date     time.Time
		expected bool
	}{
		{
			name: "День металлурга (третье воскресенье июля)",
			event: Event{
				Event: "День металлурга",
				Rule: EventRule{
					DayOfWeek:   intPtr(7), // воскресенье
					Month:       intPtr(7), // июль
					WeekOfMonth: intPtr(3), // третья неделя
				},
			},
			date:     time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC), // третье воскресенье июля 2024
			expected: true,
		},
		{
			name: "День программиста (256-й день года)",
			event: Event{
				Event: "День программиста",
				Rule: EventRule{
					DayOfYear: intPtr(256),
				},
			},
			date:     time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC), // 256-й день 2024 года (високосный)
			expected: true,
		},
		{
			name: "День шахтера (последнее воскресенье августа)",
			event: Event{
				Event: "День шахтера",
				Rule: EventRule{
					DayOfWeek:   intPtr(7), // воскресенье
					Month:       intPtr(8), // август
					WeekOfMonth: intPtr(5), // последняя неделя
				},
			},
			date:     time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), // последнее воскресенье августа 2024
			expected: true,
		},
		{
			name: "Новый год (фиксированная дата)",
			event: Event{
				Event: "Новый год",
				Rule: EventRule{
					DayOfMonth: intPtr(1),
					Month:      intPtr(1),
				},
			},
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "Неверная дата для Дня металлурга",
			event: Event{
				Event: "День металлурга",
				Rule: EventRule{
					DayOfWeek:   intPtr(7),
					Month:       intPtr(7),
					WeekOfMonth: intPtr(3),
				},
			},
			date:     time.Date(2024, 7, 14, 0, 0, 0, 0, time.UTC), // второе воскресенье июля
			expected: false,
		},
	}

	scheduler := NewScheduler(&config.Config{}, &bot.MockBot{}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scheduler.isEventToday(tt.event, tt.date)
			if result != tt.expected {
				t.Errorf("isEventToday() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetNextEventDate(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		date     time.Time
		expected time.Time
	}{
		{
			name: "День металлурга - следующая дата",
			event: Event{
				Event: "День металлурга",
				Rule: EventRule{
					DayOfWeek:   intPtr(0),
					Month:       intPtr(7),
					WeekOfMonth: intPtr(3),
				},
			},
			date:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "День программиста - следующая дата",
			event: Event{
				Event: "День программиста",
				Rule: EventRule{
					DayOfYear: intPtr(256),
				},
			},
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	scheduler := NewScheduler(&config.Config{}, &bot.MockBot{}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scheduler.getNextEventDate(tt.event, tt.date)
			if !result.Equal(tt.expected) {
				t.Errorf("getNextEventDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	events, err := Load("../events.json")
	if err != nil {
		t.Fatalf("Failed to load events: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected events to be loaded, got empty slice")
	}

	// Проверяем наличие обязательных праздников
	expectedEvents := map[string]bool{
		"День металлурга":   false,
		"День программиста": false,
		"День шахтера":      false,
		"Новый год":         false,
	}

	for _, event := range events {
		if _, ok := expectedEvents[event.Event]; ok {
			expectedEvents[event.Event] = true
		}
	}

	for event, found := range expectedEvents {
		if !found {
			t.Errorf("Expected event %q not found in loaded events", event)
		}
	}
}
