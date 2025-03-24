package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Ошибка загрузки .env")
	}
}

func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}
