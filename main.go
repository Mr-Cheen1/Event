package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
	"github.com/Mr-Cheen1/Event/events"
)

func main() {
	log.Println("🚀 Запуск Event Bot...")

	// Загрузка конфигурации
	log.Println("📋 Загрузка конфигурации...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Ошибка при загрузке конфигурации: %v", err)
	}
	log.Printf("✅ Конфигурация загружена успешно: время уведомлений %s, часовой пояс %s",
		cfg.NotificationTime, cfg.Timezone)

	// Инициализация бота
	log.Println("🤖 Инициализация Telegram бота...")
	bot, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании бота: %v", err)
	}
	log.Printf("✅ Бот успешно инициализирован, ID чата для уведомлений: %d", cfg.ChatID)

	// Загрузка событий
	log.Println("📅 Загрузка списка событий...")
	eventsList, err := events.Load("events.yml")
	if err != nil {
		log.Fatalf("❌ Ошибка при загрузке событий: %v", err)
	}
	log.Printf("✅ Загружено событий: %d", len(eventsList))

	// Запуск планировщика
	log.Println("⏰ Запуск планировщика событий...")
	scheduler := events.NewScheduler(cfg, bot, eventsList)
	scheduler.Start()

	// Информация о запуске
	loc, _ := time.LoadLocation(cfg.Timezone)
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	log.Printf("🎉 Event Bot успешно запущен! Текущее время: %s", currentTime)

	// Блокировка основного потока, чтобы приложение не завершалось
	fmt.Println("=================================================")
	fmt.Println("🟢 Бот запущен и работает! Нажмите Ctrl+C для выхода.")
	fmt.Println("=================================================")

	// Держим приложение запущенным
	select {}
}
