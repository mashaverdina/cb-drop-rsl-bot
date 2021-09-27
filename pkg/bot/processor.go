package rslbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor interface {
	Handle(msg *tgbotapi.Message) error
}
