package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type HelpCommand struct {
	BotCommand
}

func (c *HelpCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	return chatutils.TextTo(&user, "#TODO Тест для хелпа", keyboards.MainMenuKeyboard), nil
}

func (c *HelpCommand) CanHandle(cmd string) bool {
	return cmd == "help"
}
