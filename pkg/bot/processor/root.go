package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/entities"
)

type Processor interface {
	Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error)
}

type ProcessingMessage struct {
	User      entities.User
	ChatID    int64
	MessageID int
	Text      string
}

func (p ProcessingMessage) Chat() int64 {
	return p.ChatID
}

func (p ProcessingMessage) Message() int {
	return p.MessageID
}
