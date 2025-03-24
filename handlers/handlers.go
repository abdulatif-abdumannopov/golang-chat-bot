package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
	"sync"
)

const botUsername = "anonym_reply_bot"

var activeChats = make(map[int64]int64)
var replyChats = make(map[int64]int64)
var mu sync.Mutex

const ADMIN_ID int64 = 7751874405

var replyMarkup = &telebot.ReplyMarkup{}
var replyBtn = replyMarkup.Data("Ответить", "reply")

func SetupHandlers(bot *telebot.Bot, db *sql.DB) {
	// Передаём db через замыкание
	bot.Handle("/start", func(c telebot.Context) error {
		return HandleStart(db, c)
	})
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return HandleAnonymousMessage(bot, c)
	})
	bot.Handle(&telebot.Btn{Unique: "reply"}, func(c telebot.Context) error {
		return HandleReplyMode(c)
	})
}

// HandleStart /start

func HandleStart(db *sql.DB, c telebot.Context) error {
	args := c.Args()

	// Если команда вызвана без параметров
	if len(args) == 0 {
		userTelegramID := c.Sender().ID
		var userID int

		// Проверяем, есть ли пользователь в БД
		err := db.QueryRow("SELECT id FROM users WHERE telegram = ?", userTelegramID).Scan(&userID)
		if errors.Is(err, sql.ErrNoRows) {
			// Если пользователя нет, добавляем его
			result, err := db.Exec("INSERT INTO users (telegram) VALUES (?)", userTelegramID)
			if err != nil {
				log.Printf("❌ Ошибка при добавлении пользователя %d: %v", userTelegramID, err)
				return c.Send("⚠ Ошибка при сохранении в БД.")
			}

			// Получаем сгенерированный ID
			lastInsertID, _ := result.LastInsertId()
			userID = int(lastInsertID)
		} else if err != nil {
			log.Printf("❌ Ошибка при запросе пользователя %d: %v", userTelegramID, err)
			return c.Send("⚠ Ошибка при проверке БД.")
		}

		// Генерируем ссылку
		link := fmt.Sprintf("https://t.me/%s?start=%d", botUsername, userID)

		// Отправляем пользователю его ссылку
		return c.Send(fmt.Sprintf("👋 Привет! Твоя ссылка для получения анонимных сообщений:\n\n[%s](%s)", link, link), telebot.ModeMarkdown)
	}

	// 🔹 Если пользователь перешёл по ссылке (args[0] = user_id)
	targetID, err := strconv.Atoi(args[0])
	if err != nil {
		return c.Send("⚠ Неверный формат ссылки.")
	}

	// Проверяем, существует ли этот ID в БД
	var telegram int
	err = db.QueryRow("SELECT telegram FROM users WHERE id = ?", targetID).Scan(&telegram)
	if err != nil {
		return c.Send("Ошибка поиска пользователя")
	}

	// Запоминаем связь "отправитель -> получатель"
	mu.Lock()
	activeChats[c.Sender().ID] = int64(telegram)
	mu.Unlock()

	return c.Send(fmt.Sprintf("✉ Теперь ты можешь анонимно написать пользователю! Просто отправь мне сообщение."), telebot.ModeMarkdown)
}

// HandleAnonymousMessage - пересылка сообщений
func HandleAnonymousMessage(bot *telebot.Bot, c telebot.Context) error {
	sender := c.Sender()
	senderTelegramID := sender.ID

	mu.Lock()
	receiverTelegramID, exists := activeChats[senderTelegramID]
	mu.Unlock()

	if !exists {
		return c.Send("⚠ Ты ещё не выбрал, кому писать. Перейди по ссылке из /start!")
	}

	// Экранируем HTML-теги
	senderName := sender.FirstName
	if sender.Username != "" {
		senderName = fmt.Sprintf(`<a href="https://t.me/%s">%s</a>`, sender.Username, sender.FirstName)
	} else {
		senderName = fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, senderTelegramID, sender.FirstName)
	}

	// Данные отправителя
	senderInfo := fmt.Sprintf("👤 Отправитель: %s\n📎 ID: <code>%d</code>", senderName, senderTelegramID)

	// Создаём новую разметку для кнопки "Ответить"
	replyMarkup := &telebot.ReplyMarkup{}
	replyBtn := replyMarkup.Data("Ответить", "reply", strconv.FormatInt(senderTelegramID, 10))
	replyMarkup.Inline(replyMarkup.Row(replyBtn))

	// Подготавливаем текст сообщения
	messageText := fmt.Sprintf("📩 Новое анонимное сообщение:\n\n%s", c.Text())

	// Если сообщение отправляется админу — добавить данные отправителя
	if receiverTelegramID == ADMIN_ID {
		messageText += fmt.Sprintf("\n\n%s", senderInfo)
	}

	// Отправляем сообщение получателю с кнопкой "Ответить"
	_, err := bot.Send(&telebot.User{ID: receiverTelegramID}, messageText, replyMarkup, telebot.ModeHTML) // Заменил ModeMarkdown → ModeHTML
	if err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
		return c.Send("⚠ Не удалось отправить сообщение.")
	}

	// Запоминаем связь "получатель -> отправитель"
	mu.Lock()
	replyChats[receiverTelegramID] = senderTelegramID
	mu.Unlock()

	return c.Send("✅ Сообщение отправлено!")
}

// HandleReplyMode - обработка нажатия "Ответить"
func HandleReplyMode(c telebot.Context) error {
	receiverTelegramID := c.Sender().ID

	// Проверяем, есть ли активный чат для ответа
	mu.Lock()
	senderTelegramID, exists := replyChats[receiverTelegramID]
	mu.Unlock()

	if !exists {
		return c.Send("⚠ Нет сообщений для ответа.")
	}

	// Запоминаем связь "ответчик -> отправитель"
	mu.Lock()
	activeChats[receiverTelegramID] = senderTelegramID
	mu.Unlock()

	return c.Send("✉ Теперь ты можешь анонимно ответить! Просто отправь сообщение.")
}
