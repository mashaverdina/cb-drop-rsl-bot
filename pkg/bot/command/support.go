package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
)

type SupportCommand struct {
	BotCommand
}

func (c *SupportCommand) Handle(ctx context.Context, user entities.User, command string, arguments string) ([]tgbotapi.Chattable, error) {
	return chatutils.TextTo(&user, "Весь фидбек и сообщения о проблемах смело отправляй [моему ботюне](https://t.me/rsl_cb_drop_support_bot)", keyboards.MainMenuKeyboard), nil
}

func (c *SupportCommand) CanHandle(cmd string) bool {
	return cmd == "support" || cmd == "feedback"
}
