package handlers

import (
	"github.com/go-telegram/bot"
)

func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, InputHandler)
}
