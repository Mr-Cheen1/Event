package main

import (
	"log"
	"net/http"

	"github.com/Mr-Cheen1/Event/bot"
	"github.com/Mr-Cheen1/Event/config"
	"github.com/Mr-Cheen1/Event/events"
)

func main() {
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

	// Запускаем HTTP сервер для проверки работоспособности
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bot is running"))
	})

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Ошибка HTTP сервера: %v", err)
		}
	}()

	// Запускаем планировщик
	scheduler.Start()
}
