package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// "net/http" // Удалено, так как HTTP сервер не используется

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
	"github.com/Mr-Cheen1/Event/events"
)

func main() {
	log.Println("Starting application...")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime panic: %v", r)
		}
	}()

	// Проверка наличия файлов
	if _, err := os.Stat("config.json"); err != nil {
		log.Fatalf("config.json not found: %v", err)
	}
	if _, err := os.Stat("events.yml"); err != nil {
		log.Fatalf("events.yml not found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	log.Printf("Config loaded successfully: %+v", cfg)

	bot, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	eventsList, err := events.Load("events.yml")
	if err != nil {
		log.Fatalf("Ошибка при загрузке событий: %v", err)
	}

	scheduler := events.NewScheduler(cfg, bot, eventsList)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	scheduler.Start()
	log.Println("Scheduler started")

	// Блокировка основного потока
	select {}
}
