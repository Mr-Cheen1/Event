package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Сохраняем оригинальные значения переменных окружения
	origToken := os.Getenv("BOT_TOKEN")
	origChatID := os.Getenv("CHAT_ID")

	// Устанавливаем тестовые значения
	os.Setenv("BOT_TOKEN", "test_token")
	os.Setenv("CHAT_ID", "123456")

	// Восстанавливаем оригинальные значения после теста
	defer func() {
		os.Setenv("BOT_TOKEN", origToken)
		os.Setenv("CHAT_ID", origChatID)
	}()

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

func TestLoadMissingEnvVars(t *testing.T) {
	// Сохраняем оригинальные значения
	origToken := os.Getenv("BOT_TOKEN")
	origChatID := os.Getenv("CHAT_ID")

	// Очищаем переменные окружения
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("CHAT_ID")

	// Восстанавливаем оригинальные значения после теста
	defer func() {
		os.Setenv("BOT_TOKEN", origToken)
		os.Setenv("CHAT_ID", origChatID)
	}()

	_, err := Load()
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии переменных окружения, но ошибки не было")
	}
}

func TestLoadMissingConfigValues(t *testing.T) {
	// Устанавливаем тестовые значения окружения
	os.Setenv("BOT_TOKEN", "test_token")
	os.Setenv("CHAT_ID", "123456")

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
}
