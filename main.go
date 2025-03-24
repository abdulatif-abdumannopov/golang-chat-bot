package main

import (
	"log"
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
	// Запуск бота
	bot.Start()
}
