package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	NotificationTime string `json:"notification_time"`
	Timezone         string `json:"timezone"`
}

func LoadConfig(filename string) Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла конфигурации: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Ошибка при парсинге файла конфигурации: %v", err)
	}

	return config
}
