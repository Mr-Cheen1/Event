package events

import (
	"os"
	"testing"
	"time"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
)

func TestLoad(t *testing.T) {
	// Создаем временный файл событий
	tmpfile, err := os.CreateTemp("", "events.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Записываем тестовые данные во временный файл
	testEvents := []byte(`[
		{"date": "01.01", "event": "Новый год"},
		{"date": "07.01", "event": "Рождество"}
	]`)
	if _, err := tmpfile.Write(testEvents); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Загружаем события
	events, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Ошибка при загрузке событий: %v", err)
	}

	// Проверяем загруженные события
	if len(events) != 2 {
		t.Errorf("Ожидалось 2 события, получено %d", len(events))
	}
	if events[0].Date != "01.01" || events[0].Event != "Новый год" {
		t.Errorf("Неверное первое событие: %v", events[0])
	}
	if events[1].Date != "07.01" || events[1].Event != "Рождество" {
		t.Errorf("Неверное второе событие: %v", events[1])
	}
}

func TestCheckEvents(t *testing.T) {
	// Создаем мок-объекты
	mockBot := &bot.MockBot{}
	mockConfig := &config.Config{
		ChatID: 123456,
	}
	testEvents := []Event{
		{Date: "01.01", Event: "Новый год"},
		{Date: "07.01", Event: "Рождество"},
	}

	scheduler := NewScheduler(mockConfig, mockBot, testEvents)

	// Тестируем проверку событий
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	scheduler.checkEvents(now)

	// Проверяем, что было отправлено правильное сообщение
	if len(mockBot.SentMessages) != 1 {
		t.Errorf("Ожидалось 1 отправленное сообщение, получено %d", len(mockBot.SentMessages))
	}
	expectedMessage := "😏 Сегодня Новый год 😉"
	if mockBot.SentMessages[0] != expectedMessage {
		t.Errorf("Ожидалось сообщение '%s', получено '%s'", expectedMessage, mockBot.SentMessages[0])
	}
}
