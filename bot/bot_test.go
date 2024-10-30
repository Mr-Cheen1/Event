package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MockTgBotAPI - мок для tgbotapi.BotAPI.
type MockTgBotAPI struct {
	SentMessages []tgbotapi.MessageConfig
}

func (m *MockTgBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	msg, ok := c.(tgbotapi.MessageConfig)
	if !ok {
		return tgbotapi.Message{}, nil
	}
	m.SentMessages = append(m.SentMessages, msg)
	return tgbotapi.Message{}, nil
}

func TestSendMessage(t *testing.T) {
	mockAPI := &MockTgBotAPI{}
	bot := &Bot{api: mockAPI}

	chatID := int64(123456)
	text := "Test message"

	err := bot.SendMessage(chatID, text)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(mockAPI.SentMessages) != 1 {
		t.Errorf("Expected 1 sent message, got %d", len(mockAPI.SentMessages))
	}

	sentMsg := mockAPI.SentMessages[0]
	if sentMsg.ChatID != chatID {
		t.Errorf("Expected ChatID %d, got %d", chatID, sentMsg.ChatID)
	}
	if sentMsg.Text != text {
		t.Errorf("Expected text '%s', got '%s'", text, sentMsg.Text)
	}
}
