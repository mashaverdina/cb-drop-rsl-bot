package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type NotFoundCommand struct {
	BotCommand
}

func (c *NotFoundCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	return chatutils.TextTo(&user, "Команда не найдена", keyboards.MainMenuKeyboard), nil
}

func (c *NotFoundCommand) CanHandle(cmd string) bool {
	return true
}
