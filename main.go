package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"telegram/config"
	"telegram/handlers"
	"telegram/models"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	// Загружаем переменные окружения и БД
	config.LoadEnv()
	db := models.InitDB()

	// Создаём бота
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.GetBotToken(),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	handlers.SetupHandlers(bot, db)

	// Получаем порт от Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Значение по умолчанию
	}

	// 🔥 Запускаем HTTP-сервер в отдельной горутине
	go func() {
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(w, "pong")
		})

		log.Printf("🔹 Запускаем HTTP сервер на порту %s", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("❌ Ошибка запуска сервера:", err)
		}
	}()

	// Запуск бота
	bot.Start()
}
