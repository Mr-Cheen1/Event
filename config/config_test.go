package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Создаем временный .env файл
	envFile := ".env"
	err := os.WriteFile(envFile, []byte("BOT_TOKEN=test_token\nCHAT_ID=123456"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(envFile)

	// Создаем временный файл конфигурации
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	testConfig := []byte(`{
		"notification_time": "12:00",
		"timezone": "Europe/Moscow"
	}`)
	if _, err := tmpfile.Write(testConfig); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	oldConfigFile := ConfigFile
	ConfigFile = tmpfile.Name()
	defer func() { ConfigFile = oldConfigFile }()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	if cfg.NotificationTime != "12:00" {
		t.Errorf("Ожидалось время уведомления 12:00, получено %s", cfg.NotificationTime)
	}
	if cfg.Timezone != "Europe/Moscow" {
		t.Errorf("Ожидался часовой пояс Europe/Moscow, получен %s", cfg.Timezone)
	}
	if cfg.BotToken != "test_token" {
		t.Errorf("Ожидался токен бота test_token, получен %s", cfg.BotToken)
	}
	if cfg.ChatID != 123456 {
		t.Errorf("Ожидался ID чата 123456, получен %d", cfg.ChatID)
	}
}

func TestLoadMissingEnvFile(t *testing.T) {
	if err := os.Remove(".env"); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии .env файла, но ошибки не было")
	}
}

func TestLoadMissingValues(t *testing.T) {
	// Создаем временный .env файл
	envFile := ".env"
	err := os.WriteFile(envFile, []byte("BOT_TOKEN=test_token\nCHAT_ID=123456"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(envFile)

	// Создаем временный файл конфигурации
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	oldConfigFile := ConfigFile
	ConfigFile = tmpfile.Name()
	defer func() { ConfigFile = oldConfigFile }()

	// Тест на отсутствие времени уведомления
	testConfig := []byte(`{"timezone": "Europe/Moscow"}`)
	if err := os.WriteFile(tmpfile.Name(), testConfig, 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии времени уведомления, но ошибки не было")
	}

	// Тест на отсутствие часового пояса
	testConfig = []byte(`{"notification_time": "12:00"}`)
	if err := os.WriteFile(tmpfile.Name(), testConfig, 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии часового пояса, но ошибки не было")
	}

	// Тест на пустые значения
	testConfig = []byte(`{"notification_time": "", "timezone": ""}`)
	if err := os.WriteFile(tmpfile.Name(), testConfig, 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	if err == nil {
		t.Error("Ожидалась ошибка при пустых значениях, но ошибки не было")
	}
}
