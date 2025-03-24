package main

//import (
//	"log"
//
//	"gorm.io/driver/sqlite"
//	"gorm.io/gorm"
//)

//// Глобальная переменная БД
//var db *gorm.DB
//
//func ConnectDB() {
//	var err error
//	db, err = gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
//	if err != nil {
//		log.Fatal("❌ Ошибка подключения к БД:", err)
//	}
//
//	// Создаём таблицы, если их нет
//	db.AutoMigrate(&Message{})
//}
