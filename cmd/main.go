package main

import (
	"SearchServices/internal/endpoint"
	"SearchServices/internal/env"
	handlers "SearchServices/internal/telegram"
	"context"
	"github.com/go-telegram/bot"
	"os"
	"os/signal"
)

var b *bot.Bot

func main() {
	setupTelegramBot()
}

func setupTelegramBot() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token := env.BotToken()
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.AskSearchHandler),
	}
	var err error
	b, err = bot.New(token, opts...)
	if err != nil {
		panic(err)
	}
	endpoint.RegisterHandlers(b)
	b.Start(ctx)
}
