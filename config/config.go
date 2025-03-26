package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Загружаем .env ТОЛЬКО если работаем не в Render
	if os.Getenv("RENDER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ .env не найден, используем переменные окружения")
		}
	}
}

func GetBotToken() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ Переменная окружения BOT_TOKEN не задана!")
	}
	return token
}
