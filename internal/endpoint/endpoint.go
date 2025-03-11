package endpoint

import (
	handlers "SearchServices/internal/telegram"
	"github.com/go-telegram/bot"
)

func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypeExact, handlers.AskSearchHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, handlers.InputHandler)
}
