package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/bot"
	pg2 "vkokarev.com/rslbot/pkg/pg"
)

func token() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("BOT_TOKEN env variable is required")
	}
	return token
}

func loadPgVariables() (pgHosts []string, pgPort int64, pgUser string, pgPassword string, dbName string) {
	if hosts := os.Getenv("DB_HOST"); hosts == "" {
		panic("missing var DB_HOST")
	} else {
		pgHosts = strings.Split(hosts, ",")
	}

	if portStr := os.Getenv("DB_PORT"); portStr == "" {
		panic("missing var DB_HOST")
	} else {
		var err error
		pgPort, err = strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	if pgUser = os.Getenv("DB_USER"); pgUser == "" {
		panic("missing var DB_USER")
	}
	if pgPassword = os.Getenv("DB_PASSWORD"); pgPassword == "" {
		panic("missing var DB_PASSWORD")
	}
	if dbName = os.Getenv("DB_NAME"); dbName == "" {
		panic("missing var DB_NAME")
	}
	return
}

func main() {
	tgbot, err := tgbotapi.NewBotAPI(token())
	if err != nil {
		log.Panic(err)
	}

	tgbot.Debug = true

	log.Printf("Authorized on account %s", tgbot.Self.UserName)

	pgHosts, pgPort, pgUser, pgPassword, dbName := loadPgVariables()
	pg, err := pg2.NewPGClient(pgHosts, int(pgPort), pgUser, pgPassword, dbName)
	if err != nil {
		panic(err)
	}

	bot := rslbot.NewBot(tgbot, pg)
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
