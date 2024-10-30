package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// API - интерфейс для API бота.
type API interface {
	SendMessage(chatID int64, text string) error
}

// TgBotAPI - интерфейс для tgbotapi.BotAPI.
type TgBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Bot struct {
	api TgBotAPI
}

func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	return err
}

// MockBot - мок для Bot.
type MockBot struct {
	SentMessages []string
}

func (m *MockBot) SendMessage(_ int64, text string) error {
	m.SentMessages = append(m.SentMessages, text)
	return nil
}
