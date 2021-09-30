package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/entities"
)

type Processor interface {
	Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error)
	CancelFor(userID int64)
}

type ProcessingMessage struct {
	User    entities.User
	ChatID  int64
	Text    string
	Message *tgbotapi.Message
}

func (p ProcessingMessage) Chat() int64 {
	return p.ChatID
}

func (p ProcessingMessage) MessageID() int {
	return p.Message.MessageID
}

func (p ProcessingMessage) Original() string {
	return p.Message.Text
}
