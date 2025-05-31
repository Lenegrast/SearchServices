package handlers

import (
	"SearchServices/internal/prediction"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	userStates = make(map[int64]string) // Храним состояния пользователей
)

// Обработчик нажатия кнопки
func AskSearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	// Устанавливаем состояние пользователя
	userStates[chatID] = "awaiting_service_name"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введите наименование услуги",
	})
}

// Обработчик текстовых сообщений
func InputHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Проверяем состояние пользователя
	if state, ok := userStates[chatID]; ok && state == "awaiting_service_name" {
		// Обрабатываем ввод только если пользователь в нужном состоянии
		answer := prediction.FinalResponce(text)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   answer,
		})

		// Сбрасываем состояние
		delete(userStates, chatID)
		return
	}

	// Если состояние не установлено, показываем дефолтную клавиатуру
	DefaultHandler(ctx, b, update)
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Поиск услуги", CallbackData: "button"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Выберите действие",
		ReplyMarkup: kb,
	})
}
