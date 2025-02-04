package main

import (
	"fmt"
	"log"
	"net/http"

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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

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

	// Блокировка основного потока
	select {}
}
