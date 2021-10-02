package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type StartCommand struct {
	BotCommand
}

func (c *StartCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	return chatutils.TextTo(&user, "Добро пожаловать в RSL.CB бот. Используй клавиатуру внизу", keyboards.MainMenuKeyboard), nil
}

func (c *StartCommand) CanHandle(cmd string) bool {
	return cmd == "start"
}
