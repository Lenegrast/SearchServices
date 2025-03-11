package handlers

import (
	"SearchServices/internal/prediction"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

var (
	userStates = make(map[int64]string)
)

func InputHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	switch userStates[chatId] {
	case "search":
		userSend := update.Message.Text
		answer := prediction.FinalResponce(userSend)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   answer,
		})
		userStates[chatId] = "default"
	}
}
func AskSearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		logrus.Fatal("ad")
	}

	chatId := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Введите наименование услуги",
	})

	if err != nil {
		return
	}

	userStates[chatId] = "search"
}
