package events

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MockBotAPI struct {
	Messages []tgbotapi.MessageConfig
}

func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	msg, ok := c.(tgbotapi.MessageConfig)
	if !ok {
		return tgbotapi.Message{}, nil
	}
	m.Messages = append(m.Messages, msg)
	return tgbotapi.Message{}, nil
}

func TestLoadEvents(t *testing.T) {
	events := LoadEvents("../events.json")
	if len(events) == 0 {
		t.Error("Expected events to be loaded, but got none")
	}
}

func TestCheckEvents(t *testing.T) {
	mockBot := &MockBotAPI{}
	events := []Event{
		{Date: "15.10", Event: "День рождения"},
		{Date: "17.10", Event: "Годовщина"},
	}

	now := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
	chatID := int64(12345)

	CheckEvents(mockBot, events, now, chatID)

	if len(mockBot.Messages) != 2 {
		t.Errorf("Ожидалось 2 сообщения, но получено %d", len(mockBot.Messages))
	}

	expectedMessages := []string{
		"😏 Сегодня День рождения 😉",
		"🤔 Скоро Годовщина 👀",
	}

	for i, msg := range mockBot.Messages {
		if msg.Text != expectedMessages[i] {
			t.Errorf("Ожидалось сообщение '%s', но получено '%s'", expectedMessages[i], msg.Text)
		}
	}
}
