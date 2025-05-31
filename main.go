package main

import (
	webendpoint "SearchServices/internal/endpoint"
	"SearchServices/internal/env"
	handlers "SearchServices/internal/telegram"
	"context"
	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

var b *bot.Bot

func main() {
	//	setupTelegramBot()
	http.HandleFunc("/get", webendpoint.GetHandler)
	http.HandleFunc("/main", webendpoint.PostHandler)

	// Запуск сервера
	logrus.Infof("Server started. Port %s", "")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		logrus.Fatalf("Failed with start server. Error: %s", err)
	}

}

func setupTelegramBot() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token := env.BotToken()
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, handlers.AskSearchHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.DefaultHandler),
		bot.WithMessageTextHandler("", bot.MatchTypePrefix, handlers.InputHandler),
	}

	var err error
	b, err = bot.New(token, opts...)
	if err != nil {
		panic(err)
	}
	//b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, handlers.InputHandler)
	//endpoint.RegisterHandlers(b)
	b.Start(ctx)

}
