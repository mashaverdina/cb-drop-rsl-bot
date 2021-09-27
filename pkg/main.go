package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/bot"
)

func token() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("BOT_TOKEN env variable is required")
	}
	return token
}

func main() {
	tgbot, err := tgbotapi.NewBotAPI(token())
	if err != nil {
		log.Panic(err)
	}

	tgbot.Debug = true

	log.Printf("Authorized on account %s", tgbot.Self.UserName)

	bot := rslbot.NewBot(tgbot)
	if err := bot.Start(context.Background()); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	<-c

	log.Println("stopping bot")
	if err := bot.Stop(time.Second); err != nil {
		log.Fatalf("bot was not stopped: %v", err)
	}
}
