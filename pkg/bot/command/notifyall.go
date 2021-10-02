package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/storage"
)

type NotifyAllCommand struct {
	userStorage *storage.UserStorage
	msgQueue    chan<- []tgbotapi.Chattable
}

func NewNotifyAllCommand(userStorage *storage.UserStorage, msgQueue chan<- []tgbotapi.Chattable) *NotifyAllCommand {
	return &NotifyAllCommand{
		userStorage: userStorage,
		msgQueue:    msgQueue,
	}
}

func (c *NotifyAllCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	if !user.HasSudo {
		msg := tgbotapi.NewMessage(user.UserID, "Для данной команды требуются супер права")
		return []tgbotapi.Chattable{msg}, nil
	}
	err := c.NotifyAll(ctx, arguments)
	return nil, err
}

func (c *NotifyAllCommand) CanHandle(cmd string) bool {
	return cmd == "notifyall"
}

func (c *NotifyAllCommand) NotifyAll(ctx context.Context, arguments string) error {
	allUsers, err := c.userStorage.All(ctx)
	if err != nil {
		return err
	}
	for _, user := range allUsers {
		select {
		case c.msgQueue <- []tgbotapi.Chattable{tgbotapi.NewMessage(user.UserID, arguments)}:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
