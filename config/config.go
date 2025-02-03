package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ConfigFile - имя файла конфигурации.
var ConfigFile = "config.json"

type Config struct {
	NotificationTime string `json:"notification_time"`
	Timezone         string `json:"timezone"`
	BotToken         string
	ChatID           int64
	EventsFile       string
}

func Load() (*Config, error) {
	// Пробуем загрузить .env файл, но не возвращаем ошибку если его нет
	_ = godotenv.Load()

	// Получаем значения из переменных окружения
	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("CHAT_ID не задан в переменных окружения")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("неверный формат CHAT_ID: %w", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN не задан в переменных окружения")
	}

	cfg := &Config{
		BotToken:   botToken,
		ChatID:     chatID,
		EventsFile: "events.yml",
	}

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("ошибка разбора файла конфигурации: %w", err)
	}

	if cfg.NotificationTime == "" {
		return nil, fmt.Errorf("время уведомления не задано в config.json")
	}
	if cfg.Timezone == "" {
		return nil, fmt.Errorf("часовой пояс не задан в config.json")
	}

	return cfg, nil
}
