package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/entities"
)

type BotCommand interface {
	Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error)
	CanHandle(cmd string) bool
}
